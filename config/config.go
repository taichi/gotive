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
package config

import (
	"github.com/BurntSushi/toml"
	"github.com/taichi/gotive/log"
	"os/exec"
)

type commitDefaults struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type gotiveConfig struct {
	Port   uint           `toml:"port"`
	Repo   string         `toml:"repo"`
	Git    string         `toml:"git"`
	Commit commitDefaults `toml:"commit_defaults"`
}

type Config *gotiveConfig

func New() Config {
	return &gotiveConfig{
		Port: 8080,
		Repo: "./repo",
		Git:  "git",
		Commit: commitDefaults{
			Name:  "anonymous",
			Email: "anonymous@example.com",
		},
	}
}

func Load(path string) Config {
	config := New()

	if _, err := toml.DecodeFile(path, config); err != nil {
		log.Warn(err)
		log.Warnf("Default values are %v", config)
	}

	if _, err := exec.LookPath(config.Git); err != nil {
		log.Fatal(err)
	}

	return config
}
