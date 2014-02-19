/* Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package repo

import (
	"bytes"
	"fmt"
	c "github.com/taichi/gotive/config"
	"github.com/taichi/gotive/log"
	"github.com/taichi/osutil"
	"github.com/taichi/rand"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type gotiveRepos struct {
	config c.Config
	rs     *rand.RandomStringer
}

type RepoMaker interface {
	MakeRepo() (Repo, error)
	LoadRepo(repoid string) (Repo, error)
}

func New(c c.Config) RepoMaker {
	if err := os.MkdirAll(c.Repo, os.ModeDir); err != nil {
		panic(err)
	}
	return &gotiveRepos{
		config: c,
		rs:     rand.Alnum(),
	}
}

var FailToMakeRepo = fmt.Errorf("Fail to make repository")

type gotiveRepo struct {
	config   c.Config
	id, root string
}

type Repo interface {
	Id() string
	Desc() (string, error)
	ApplyDesc(desc string) error
	Add(name, content string) error
	Commit(name, email string) error
	Walk(fn filepath.WalkFunc) error
	ReadFile(path string) ([]byte, error)
}

// TODO use promise or future?
// https://sites.google.com/site/gopatterns/concurrency/futures
func (r *gotiveRepos) MakeRepo() (Repo, error) {
	for i := 0; i < 3; i++ {
		newid := r.rs.Next(16)
		newone := filepath.Join(r.config.Repo, newid)
		if osutil.IsExist(newone) {
			continue
		}
		if err := os.MkdirAll(newone, os.ModeDir); err != nil {
			log.Debug(err)
			continue
		}
		if err := run(r.config, newone, []string{"init"}); err == nil {
			return &gotiveRepo{
				id:     newid,
				config: r.config,
				root:   newone}, nil
		} else {
			log.Debug(err)
		}
	}
	return nil, FailToMakeRepo
}

func (r *gotiveRepos) LoadRepo(repoid string) (Repo, error) {
	repo := filepath.Join(r.config.Repo, repoid)
	if _, err := os.Lstat(repo); err != nil {
		return nil, err
	}
	return &gotiveRepo{id: repoid, config: r.config, root: repo}, nil
}

func (r *gotiveRepo) Id() string {
	return r.id
}

func (r *gotiveRepo) DescPath() string {
	return filepath.Join(r.root, ".git/description")
}

func (r *gotiveRepo) Desc() (string, error) {
	b, err := ioutil.ReadFile(r.DescPath())
	return string(b), err
}

func (r *gotiveRepo) ApplyDesc(desc string) error {
	return ioutil.WriteFile(r.DescPath(), []byte(desc), 0)
}

func run(c c.Config, root string, options []string, env ...map[string]string) error {
	cmd := exec.Command(c.Git, options...)
	cmd.Dir = root
	cmd.Env = mergeEnv(env...)

	if log.IsDebugEnabled() {
		var out, err bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &err
		defer func() {
			for _, b := range []bytes.Buffer{out, err} {
				s := b.String()
				if 0 < len(s) {
					log.Debug(s)
				}
			}
		}()
	}
	return cmd.Run() // TODO timeout
}

func mergeEnv(newmaps ...map[string]string) []string {
	out := os.Environ()
	for _, m := range newmaps {
		for k, v := range m {
			prefix := fmt.Sprintf("%s=", k)
			kv := prefix + v
			if x, i := find(out, matcher(prefix)); x {
				out[i] = kv
			} else {
				out = append(out, kv)
			}
		}
	}
	return out
}

func matcher(prefix string) Matcher {
	return func(v string) bool {
		return strings.HasPrefix(v, prefix)
	}
}

type Matcher func(value string) bool

func find(target []string, fn Matcher) (found bool, index int) {
	for i, v := range target {
		if fn(v) {
			return true, i
		}
	}
	return false, -1
}

func (r *gotiveRepo) Add(name, content string) error {
	p := filepath.Join(r.root, name)

	if osutil.Contains(r.root, p) == false {
		return fmt.Errorf("Unsupported path %s", name)
	}

	if osutil.IsExist(p) {
		return fmt.Errorf("Already Exists %s", name) // override? merge?
	}

	if err := os.MkdirAll(filepath.Dir(p), 0644); err != nil {
		return err
	}

	if err := ioutil.WriteFile(p, []byte(content), 0644 /*-rw-r--r--*/); err != nil {
		return err
	}
	if rel, err := filepath.Rel(r.root, p); err != nil {
		return err
	} else {
		return run(r.config, r.root, []string{"add", rel})
	}
}

func (r *gotiveRepo) Commit(name, email string) error {
	return run(r.config, r.root, []string{"commit", "--allow-empty-message", "-m", ""}, r.makeEnv(name, email))
}

func (r *gotiveRepo) makeEnv(name, email string) map[string]string {
	env := map[string]string{}

	resolve := func(val, def string) string {
		if 0 < len(val) {
			return val
		}
		return def
	}
	n := resolve(name, r.config.Commit.Name)
	env["GIT_AUTHOR_NAME"] = n
	env["GIT_COMMITTER_NAME"] = n

	e := resolve(email, r.config.Commit.Email)
	env["GIT_AUTHOR_EMAIL"] = e
	env["GIT_COMMITTER_EMAIL"] = e
	return env
}

func (r *gotiveRepo) Walk(fn filepath.WalkFunc) error {
	return filepath.Walk(r.root, func(path string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(r.root, path)
		if err != nil {
			return err
		}
		if isGit(rel) {
			return nil
		}
		return fn(rel, info, err)
	})
}

func isGit(path string) bool {
	return -1 < strings.Index(path, ".git")
}

func (r *gotiveRepo) ReadFile(path string) ([]byte, error) {
	p := filepath.Join(r.root, path)

	if osutil.Contains(r.root, p) == false || isGit(p) {
		return nil, fmt.Errorf("Unsupported path %s", p)
	}

	return ioutil.ReadFile(p)
}
