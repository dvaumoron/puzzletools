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
	"errors"

	"github.com/dvaumoron/puzzletools/initlogindb"
	"github.com/urfave/cli/v2"
)

var errNoLoginOrPassword = errors.New("userLogin and userPassword required")

var saltServiceAddr string
var loginServiceAddr string

var initLoginCmd = &cli.Command{
	Name:        "initlogin",
	Usage:       "init login database",
	ArgsUsage:   "userLogin userPassword",
	Description: "init login database : create an user with the arguments",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "salt-service-addr",
			Usage:       "Address of the salt service",
			EnvVars:     []string{"SALT_SERVICE_ADDR"},
			Destination: &saltServiceAddr,
		},
		&cli.StringFlag{
			Name:        "login-service-addr",
			Usage:       "Address of the right service",
			EnvVars:     []string{"LOGIN_SERVICE_ADDR"},
			Destination: &loginServiceAddr,
		},
	},
	Action: func(cCtx *cli.Context) error {
		err := errUnknownServiceAddr
		if saltServiceAddr != "" && loginServiceAddr != "" {
			err = errNoLoginOrPassword
			if args := cCtx.Args(); args.Len() > 1 {
				err = initlogindb.InitUser(saltServiceAddr, loginServiceAddr, args.Get(0), args.Get(1))
			}
		}
		return err
	},
}
