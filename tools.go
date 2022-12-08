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
	outPath := addSlash(os.Args[1])
	inPath := outPath + "fragments/"

	placeHolderCss := "{{.Css}}"
	placeHolderJs := "{{.Js}}"
	placeHolderBody := "{{.Body}}"

	data, err := os.ReadFile(outPath + "templates/main.html")
	if err == nil {
		tmpl := string(data)

		data, err = os.ReadFile(outPath + "static/main.css")
		if err == nil {
			cssIndex := strings.Index(tmpl, placeHolderCss)
			part1 := tmpl[0:cssIndex]
			cssIndexEnd := cssIndex + 8
			jsIndex := strings.Index(tmpl, placeHolderJs)
			part2 := tmpl[cssIndexEnd:jsIndex]
			jsIndexEnd := jsIndex + 7
			bodyIndex := strings.Index(tmpl, placeHolderBody)
			part3 := tmpl[jsIndexEnd:bodyIndex]
			bodyIndexEnd := bodyIndex + 9
			part4 := tmpl[bodyIndexEnd:]

			inSize := len(inPath)
			initJs := "<script type=\"text/javascript\" src=\"/static/"
			endJs := "\"/>\n"

			cssBody := string(data)

			err = filepath.WalkDir(inPath, func(path string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					destPath := outPath + path[inSize:]
					if destPath[len(destPath)-5:] == ".html" {
						cssRef := ""
						var jsRefs []string
						body := ""

						cssRef, jsRefs, body, err = parseHtmlFragment(path)
						if err == nil {
							err = makeDirectory(destPath, len(d.Name()))
							if err == nil {
								var bodyBuilder strings.Builder
								bodyBuilder.WriteString(part1)
								bodyBuilder.WriteString(cssRef)
								bodyBuilder.WriteString(part2)
								for _, jsRef := range jsRefs {
									bodyBuilder.WriteString(initJs)
									bodyBuilder.WriteString(jsRef)
									bodyBuilder.WriteString(endJs)
								}
								bodyBuilder.WriteString(part3)
								bodyBuilder.WriteString(body)
								bodyBuilder.WriteString(part4)
								body := []byte(bodyBuilder.String())

								err = os.WriteFile(destPath, body, 0644)
							}
						}
					} else if destPath[len(destPath)-4:] == ".css" {
						var dataSup []byte
						dataSup, err = os.ReadFile(path)
						if err == nil {
							err = makeDirectory(destPath, len(d.Name()))
							if err == nil {
								body := []byte(cssBody + string(dataSup))
								err = os.WriteFile(destPath, body, 0644)
							}
						}
					}
				}
				return err
			})
		}
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

func parseHtmlFragment(path string) (string, []string, string, error) {
	cssRef := ""
	var jsRefs []string
	body := ""

	file, err := os.Open(path)
	if err == nil {
		jsRefs = make([]string, 0)
		var bodyBuilder strings.Builder

		cssStarted := false
		jsStarted := false
		bodyStarted := false

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && trimmed[0] != '#' {
				if bodyStarted {
					bodyBuilder.WriteString(line)
					bodyBuilder.WriteString("\n")
				} else if jsStarted {
					if trimmed == "Body:" {
						bodyStarted = true
					} else {
						jsRefs = append(jsRefs, trimmed)
					}
				} else if cssStarted {
					cssRef = trimmed
					cssStarted = false
				} else if trimmed == "Js:" {
					jsStarted = true
				} else if trimmed == "Body:" {
					bodyStarted = true
				} else {
					// Should be the start of css zone.
					cssStarted = true
				}
			}
		}

		if err = scanner.Err(); err == nil {
			body = bodyBuilder.String()
		}
	}
	return cssRef, jsRefs, body, err
}

func makeDirectory(path string, size int) error {
	i := len(path) - size
	return os.MkdirAll(path[:i], 0755)
}
