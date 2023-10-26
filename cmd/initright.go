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
	"github.com/spf13/cobra"
)

var rightServiceAddr string

func newInitRightCmd(defaultRightServiceAddr string) *cobra.Command {
	initRightCmd := &cobra.Command{
		Use:   "initright userId",
		Short: "Init right database",
		Long: `Init right database :
 - init default role
 - give the user with the id in argument, this role`,
		Args: cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			userId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			if rightServiceAddr == "" {
				return errUnknownServiceAddr
			}
			return initrightdb.MakeUserAdmin(rightServiceAddr, userId)
		},
	}

	initRightCmd.Flags().StringVar(
		&rightServiceAddr, "right-service-addr", defaultRightServiceAddr,
		"Address of the right service",
	)

	return initRightCmd
}
