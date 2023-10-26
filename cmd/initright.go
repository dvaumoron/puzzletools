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
	"strconv"

	"github.com/dvaumoron/puzzletools/initrightdb"
	"github.com/urfave/cli/v2"
)

var errNoUserId = errors.New("userId required")

var rightServiceAddr string

var initRightCmd = &cli.Command{
	Name:      "initright",
	Usage:     "init right database",
	ArgsUsage: "userId",
	Description: `init right database :
- init default role
- give the user with the id in argument, this role`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "right-service-addr",
			Usage:       "Address of the right service",
			EnvVars:     []string{"RIGHT_SERVICE_ADDR"},
			Destination: &rightServiceAddr,
		},
	},
	Action: func(cCtx *cli.Context) error {
		err := errNoUserId
		if args := cCtx.Args(); args.Present() {
			var userId uint64
			if userId, err = strconv.ParseUint(args.Get(0), 10, 64); err == nil {
				if rightServiceAddr == "" {
					err = errUnknownServiceAddr
				} else {
					err = initrightdb.MakeUserAdmin(rightServiceAddr, userId)
				}
			}
		}
		return err
	},
}
