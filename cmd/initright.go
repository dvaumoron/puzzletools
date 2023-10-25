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
	"strconv"

	"github.com/dvaumoron/puzzletools/initrightdb"
	"github.com/urfave/cli/v2"
)

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
		userId, err := strconv.ParseUint(cCtx.Args().Get(0), 10, 64)
		if err == nil {
			err = initrightdb.MakeUserAdmin(rightServiceAddr, userId)
		}
		return cli.Exit(err, 1)
	},
}
