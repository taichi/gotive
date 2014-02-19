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
)

var configpath string
var verbose bool

func Execute() {
	rootCmd := newRootCmd()
	addCommands(rootCmd)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "gotive",
		Run: helpFn,
	}
	rootCmd.PersistentFlags().StringVarP(&configpath, "config", "c", "config.toml", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "print debug informations")
	return rootCmd
}

func addCommands(cmd *cobra.Command) {
	addServerCommands(cmd)
}

func helpFn(cmd *cobra.Command, args []string) { cmd.Help() }

type cmdFn func(cmd *cobra.Command, args []string)

func wrapRunFn(f func(cmd *cobra.Command, c config.Config, args []string)) cmdFn {
	return func(cmd *cobra.Command, args []string) {
		if verbose {
			log.VerboseLog()
		}
		f(cmd, config.Load(configpath), args)
	}
}
