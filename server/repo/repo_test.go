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
package repo_test

import (
	. "."
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/taichi/gotive/config"
	"github.com/taichi/gotive/ginkgo"
	"github.com/taichi/gotive/log"
	"github.com/taichi/osutil"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ = Describe("RepoMaker", func() {
	var (
		c    config.Config
		rm   RepoMaker
		root string
	)
	BeforeEach(func() {
		log.VerboseLog()
		c = config.New()
		root = filepath.Join(os.TempDir(), "repo")
		maerr := os.MkdirAll(root, os.ModeDir)
		Expect(maerr).To(BeNil())
		p, err := ioutil.TempDir(root, "")
		Expect(err).To(BeNil())
		c.Repo = p
		rm = New(c)
	})
	AfterEach(func() {
		Expect(osutil.ForceRemoveAll(root)).To(BeNil())
	})
	Context("Repo", func() {

		repoOk := func(r Repo, err error) Repo {
			Expect(err).To(BeNil())
			Expect(r).NotTo(BeNil())
			return r
		}

		It("should work normally", func() {
			r := repoOk(rm.MakeRepo())
			file, fe := os.Open(filepath.Join(c.Repo, r.Id()))
			Expect(fe).To(BeNil())
			Expect(file).NotTo(BeNil())

			r2 := repoOk(rm.LoadRepo(r.Id()))
			Expect(r2.Id()).To(Equal(r.Id()))

			repodir, fe := os.Open(filepath.Join(c.Repo, r.Id(), ".git"))
			Expect(fe).To(BeNil())
			info, se := repodir.Stat()
			Expect(se).To(BeNil())
			Expect(info.IsDir()).To(BeTrue())

		})
		str := "abcdefghijklmnop"
		seed := str + str + str + str

		It("should return error", ginkgo.FixSeed(seed, func() {
			rm.MakeRepo()
			repo, err := rm.MakeRepo()
			Expect(repo).To(BeNil())
			Expect(err).To(Equal(FailToMakeRepo))
		}))

		It("add contens normally", func() {
			r := repoOk(rm.MakeRepo())
			name, content := "hoge.txt", "hogehoge"
			err := r.Add(name, content)
			Expect(err).To(BeNil())
			read, re := ioutil.ReadFile(filepath.Join(c.Repo, r.Id(), name))
			Expect(re).To(BeNil())
			Expect(string(read)).To(Equal(content))
		})

		It("commit normally", func() {
			r := repoOk(rm.MakeRepo())
			name, content := "hoge/moge.txt", "hogehoge"
			err := r.Add(name, content)
			Expect(err).To(BeNil())
			Expect(r.Commit("way", "wayway@example.com")).To(BeNil())
		})
	})
})
