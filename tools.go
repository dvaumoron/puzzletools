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
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	projectPath := os.Args[1]
	if last := len(projectPath) - 1; projectPath[last] == '/' {
		projectPath = projectPath[:last]
	}
	inPath := projectPath + "/fragments/"
	outPath := projectPath + "/templates/"

	data, err := os.ReadFile(outPath + "main.html")
	if err == nil {
		tmpl := string(data)

		jsIndex := strings.Index(tmpl, "{{.WidgetHeader}}")
		part1 := tmpl[0:jsIndex]
		jsIndexEnd := jsIndex + 7
		bodyIndex := strings.Index(tmpl, "{{.WidgetBody}}")
		part2 := tmpl[jsIndexEnd:bodyIndex]
		bodyIndexEnd := bodyIndex + 9
		part3 := tmpl[bodyIndexEnd:]

		inSize := len(inPath)
		initJs := "<script type=\"text/javascript\" src=\"/static/"
		endJs := "\"/>\n"

		err = filepath.WalkDir(inPath, func(path string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() {
				destPath := outPath + path[inSize:]
				if destPath[len(destPath)-5:] == ".html" {
					var jsRefs []string
					var bodyLines []string

					jsRefs, bodyLines, err = parseHtmlFragment(path)
					if err == nil {
						err = makeDirectory(destPath, len(d.Name()))
						if err == nil {
							var bodyBuilder strings.Builder
							bodyBuilder.WriteString(part1)
							for _, jsRef := range jsRefs {
								bodyBuilder.WriteString(initJs)
								bodyBuilder.WriteString(jsRef)
								bodyBuilder.WriteString(endJs)
							}
							bodyBuilder.WriteString(part2)
							for _, line := range bodyLines {
								bodyBuilder.WriteString(line)
								bodyBuilder.WriteByte('\n')
							}
							bodyBuilder.WriteString(part3)
							body := []byte(bodyBuilder.String())

							err = os.WriteFile(destPath, body, 0644)
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

func parseHtmlFragment(path string) ([]string, []string, error) {
	var jsRefs []string
	var bodyLines []string

	file, err := os.Open(path)
	if err == nil {
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if trimmed := strings.TrimSpace(scanner.Text()); trimmed != "" {
				if trimmed == "Body:" {
					break
				} else {
					jsRefs = append(jsRefs, trimmed)
				}
			}
		}
		for scanner.Scan() {
			if trimmed := strings.TrimSpace(scanner.Text()); trimmed != "" {
				bodyLines = append(bodyLines, trimmed)
			}
		}

		if err = scanner.Err(); err == nil {
			if len(bodyLines) == 0 {
				bodyLines = jsRefs
				jsRefs = nil
			}
		}
	}
	return jsRefs, bodyLines, err
}

func makeDirectory(path string, nameSize int) error {
	i := len(path) - nameSize
	return os.MkdirAll(path[:i], 0755)
}
