/*
 *
 * Copyright 2022 puzzletools authors.
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
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// TODO check arguments number on subcommands
// TODO also check service addr presence
var app = &cli.App{
	Usage:     "puzzletools includes diverse features helping in puzzle project.",
	UsageText: "puzzletools command",
	Commands: []*cli.Command{
		initRightCmd,
		initLoginCmd,
		prepareCmd,
		checkPasswordCmd,
	},
	Suggest: true,
}

func Execute() {
	if godotenv.Overload() == nil {
		fmt.Println("Loaded .env file")
	}
	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
