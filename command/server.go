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
package command

import (
	"github.com/spf13/cobra"
	"github.com/taichi/gotive/config"
	"github.com/taichi/gotive/log"
	"github.com/taichi/gotive/server"
)

func addServerCommands(cmd *cobra.Command) {
	serverCmd := &cobra.Command{
		Use: "server",
		Run: helpFn,
	}
	serverCmd.AddCommand(&cobra.Command{
		Use: "start",
		Run: wrapRunFn(start),
	})
	cmd.AddCommand(serverCmd)
}

func start(cmd *cobra.Command, c config.Config, args []string) {
	if err := server.Start(c); err != nil {
		log.Fatal(err)
	}
}
