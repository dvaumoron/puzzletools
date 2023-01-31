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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultServiceAddr = "localhost:50051"

var rootCmd = &cobra.Command{
	Use:   "puzzletools [command]",
	Short: "puzzletools includes diverse features helping in puzzle project.",
	Long: `puzzletools includes the following features:
- prepare templates
- init login db
- init right db`,
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", ".env", "config file")
}

func initConfig() {
	viper.SetConfigFile(cfgFile)

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	viper.SetDefault(rightServiceAddrName, defaultServiceAddr)
	viper.SetDefault(saltServiceAddrName, defaultServiceAddr)
	viper.SetDefault(loginServiceAddrName, defaultServiceAddr)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file :", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
