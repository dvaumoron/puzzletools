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
package puzzletools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	args := os.Args
	inPath := args[1]
	outPath := args[2]

	if inPath[len(inPath)-1] != '/' {
		inPath += "/"
	}
	if outPath[len(outPath)-1] != '/' {
		outPath += "/"
	}

	tmplName := "main.html"
	data, err := os.ReadFile(inPath + tmplName)
	if err == nil {
		tmpl := string(data)

		err = filepath.WalkDir(inPath, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				name := path[len(inPath):]
				if name[len(name)-5:] == ".html" {
					if name != tmplName {
						var data []byte
						data, err = os.ReadFile(path)
						if err == nil {
							body := strings.Replace(tmpl, "{{.Body}}", string(data), 1)
							dest := outPath + name
							err = os.WriteFile(dest, []byte(body), d.Type().Perm())
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
