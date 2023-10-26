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

	"github.com/spf13/cobra"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var errNoDefaultPassword = errors.New("no default password to compare")

func newCheckPassword(defaultPass string) *cobra.Command {
	return &cobra.Command{
		Use:   "check password",
		Short: "Check the strength of a password",
		Long:  "Check the strength of the password : return indications to improve your proposal",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			err := errNoDefaultPassword
			if defaultPass != "" {
				minEntropy := passwordvalidator.GetEntropy(defaultPass)
				if err = passwordvalidator.Validate(args[0], minEntropy); err == nil {
					fmt.Println("password seems strong")
				}
			}
			return err
		},
	}
}
