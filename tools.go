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
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	inPath := addSlash(args[1])
	outPath := addSlash(args[2])

	tmplName := "main.html"
	data, err := os.ReadFile(outPath + tmplName)
	if err == nil {
		tmpl := string(data)

		err = filepath.WalkDir(inPath, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				name := path[len(inPath):]
				if name[len(name)-5:] == ".html" {
					var data []byte
					data, err = os.ReadFile(path)
					if err == nil {
						destPath := outPath + name
						err = makeDirectory(destPath)
						if err == nil {
							body := strings.Replace(tmpl, "{{.Body}}", string(data), 1)
							err = os.WriteFile(destPath, []byte(body), 0666)
						}
					}
				}
			}
			return err
		})
	}

	if err != nil {
		fmt.Println("An error occurred :", err)
	}
}

func addSlash(path string) string {
	if path[len(path)-1] != '/' {
		path += "/"
	}
	return path
}

func makeDirectory(path string) error {
	i := len(path) - 1
	for path[i] != '/' {
		i--
	}
	return os.MkdirAll(path[:i], 0755)
}
