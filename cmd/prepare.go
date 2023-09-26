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
	"os"

	"github.com/dvaumoron/puzzletools/preparetemplates"
	"github.com/spf13/cobra"
)

func newPrepareCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "prepare [projectPath]",
		Short: "prepare templates.",
		Long: `prepare templates :
- without arg work in the current folder
- walk subfolder "/fragments" and write in subfolder "/templates".
- inject the walked file in "/templates/main.html" to generate complete file`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			var path string
			if len(args) == 0 {
				var err error
				path, err = os.Getwd()
				if err != nil {
					return err
				}
			} else {
				path = args[0]
			}
			return preparetemplates.PrepareTemplates(path)
		},
	}
}
