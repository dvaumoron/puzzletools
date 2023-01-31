/*
 *
 * Copyright 2023 puzzletools authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package cmd

import (
	"github.com/dvaumoron/puzzletools/initlogindb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const saltServiceAddrName = "salt-service-addr"
const loginServiceAddrName = "login-service-addr"

var saltServiceAddr string
var loginServiceAddr string

func init() {
	initLoginCmd := &cobra.Command{
		Use:   "initlogin userLogin userPassword",
		Short: "init login database.",
		Long:  "init login database : create an user with the arguments",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return initlogindb.InitUser(saltServiceAddr, loginServiceAddr, args[0], args[1])
		},
	}

	initLoginCmdFlags := initLoginCmd.Flags()
	initLoginCmdFlags.StringVar(
		&saltServiceAddr, saltServiceAddrName, "", "Address of the salt service",
	)
	viper.BindPFlag(saltServiceAddrName, initLoginCmdFlags.Lookup(saltServiceAddrName))
	initLoginCmdFlags.StringVar(
		&loginServiceAddr, loginServiceAddrName, "", "Address of the login service",
	)
	viper.BindPFlag(loginServiceAddrName, initLoginCmdFlags.Lookup(loginServiceAddrName))

	rootCmd.AddCommand(initLoginCmd)
}
