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
package handler

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/taichi/gotive/config"
	"github.com/taichi/gotive/server/repo"
	"os"
)

type content struct {
	Name, Content string
}

func ViewEntry(res render.Render, p martini.Params, c config.Config) {
	maker := repo.New(c)
	r, err := maker.LoadRepo(p["id"])
	if err != nil {
		handleError(res, err)
		return
	}

	model := map[string]interface{}{}

	if desc, err := r.Desc(); err != nil {
		handleError(res, err)
		return
	} else {
		model["desc"] = desc
	}

	contents := []content{}
	walkRepo := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			if c, err := r.ReadFile(path); err == nil {
				contents = append(contents, content{Name: path, Content: string(c)})
			}
		}
		return nil
	}
	if err := r.Walk(walkRepo); err != nil {
		handleError(res, err)
		return
	}
	model["contents"] = contents

	res.HTML(200, "render", model)
}
