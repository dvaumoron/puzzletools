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
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var errNoDefaultPassword = errors.New("no default password to compare")

var checkPasswordCmd = &cli.Command{
	Name:        "check",
	Usage:       "check the strength of a password",
	ArgsUsage:   "password",
	Description: "check the strength of the password : return indications to improve your proposal",
	Action: func(cCtx *cli.Context) error {
		err := errNoDefaultPassword
		if defaultPass := os.Getenv("DEFAULT_PASSWORD"); defaultPass != "" {
			minEntropy := passwordvalidator.GetEntropy(defaultPass)
			if err = passwordvalidator.Validate(cCtx.Args().Get(0), minEntropy); err == nil {
				fmt.Println("password seems strong")
			}
		}
		return err
	},
}
