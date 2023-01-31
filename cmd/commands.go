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
	"github.com/spf13/cobra"
)

const UnknownServiceAddr = "service address unknown"

var rootCmd = &cobra.Command{
	Use:   "puzzletools [command]",
	Short: "puzzletools includes diverse features helping in puzzle project.",
	Long: `puzzletools includes the following features:
- prepare templates
- init login db
- init right db`,
}

func init() {
	if godotenv.Overload() == nil {
		fmt.Println("Loaded .env file")
	}

	rootCmd.AddCommand(newInitRightCmd(os.Getenv("RIGHT_SERVICE_ADDR")))
	rootCmd.AddCommand(newInitLoginCmd(os.Getenv("SALT_SERVICE_ADDR"), os.Getenv("LOGIN_SERVICE_ADDR")))
	rootCmd.AddCommand(newPrepareCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
