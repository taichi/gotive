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
	"fmt"
	"github.com/martini-contrib/render"
	"github.com/taichi/gotive/config"
	"github.com/taichi/gotive/server/repo"
	"net/http"
)

func NewEntry(req *http.Request, c config.Config, res render.Render) {
	maker := repo.New(c)
	r, err := maker.MakeRepo()

	if err != nil {
		handleError(res, err)
		return
	}

	if err := r.ApplyDesc(req.FormValue("d")); err != nil {
		handleError(res, err)
		return
	}

	contents := req.Form["c"]
	clen := len(contents)
	for index, filename := range req.Form["n"] {
		if 0 < len(filename) && index < clen {
			content := contents[index]
			if err := r.Add(filename, content); err != nil {
				handleError(res, err)
				return
			}
		}
	}

	// TODO login
	if err := r.Commit("", ""); err != nil {
		handleError(res, err)
		return
	}
	res.Redirect(fmt.Sprintf("/%s", r.Id()))
}
