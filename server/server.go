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
package server

import (
	"fmt"
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/taichi/gotive/config"
	"github.com/taichi/gotive/server/handler"
	"net/http"
)

func classic() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("server/public"))
	m.Use(render.Renderer(render.Options{
		Directory:  "server/template",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
	}))
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}

func Start(c config.Config) error {
	m := classic()
	m.Map(c)
	handler.AddHandlers(m)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.Port), m)
}
