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

package preparetemplates

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	markdownextension "github.com/dvaumoron/puzzlemarkdownextension"
	"github.com/yuin/goldmark"
)

const htmlExt = "html"

const partSep = "---"

const headerPlaceHolder = "[[.WidgetHeader]]"
const headerPlaceHolderLen = len(headerPlaceHolder)
const bodyPlaceHolder = "[[.WidgetBody]]"
const bodyPlaceHolderLen = len(bodyPlaceHolder)

const initJs = "<script type=\"text/javascript\" src=\"/static/"
const endJs = "\"/>\n"

func PrepareTemplates(projectPath string) error {
	if last := len(projectPath) - 1; projectPath[last] == '/' {
		projectPath = projectPath[:last]
	}
	inPath := projectPath + "/fragments/"
	outPath := projectPath + "/templates/"

	data, err := os.ReadFile(outPath + "main.html")
	if err != nil {
		return err
	}

	tmpl := string(data)

	jsIndex := strings.Index(tmpl, headerPlaceHolder)
	part1 := tmpl[:jsIndex]
	jsIndexEnd := jsIndex + headerPlaceHolderLen
	bodyIndex := strings.Index(tmpl[jsIndexEnd:], bodyPlaceHolder) + jsIndexEnd
	part2 := tmpl[jsIndexEnd:bodyIndex]
	bodyIndexEnd := bodyIndex + bodyPlaceHolderLen
	part3 := tmpl[bodyIndexEnd:]

	inSize := len(inPath)

	engine := newMarkdownEngine()

	return filepath.WalkDir(inPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		destPath := outPath + path[inSize:]
		dotIndex := strings.LastIndexByte(destPath, '.')
		if dotIndex == -1 {
			return nil
		}
		dotIndex++
		computeBody := noCompute
		switch ext := destPath[dotIndex:]; ext {
		case htmlExt:
		case "md":
			destPath = destPath[:dotIndex] + htmlExt
			computeBody = engine.markdownCompute
		default:
			return nil
		}

		jsRefs, bodyLines, err := parseHtmlFragment(path)
		if err != nil {
			return err
		}

		err = makeDirectory(destPath, len(d.Name()))
		if err != nil {
			return err
		}

		var bodyBuilder strings.Builder
		bodyBuilder.WriteString(part1)
		for _, jsRef := range jsRefs {
			bodyBuilder.WriteString(initJs)
			bodyBuilder.WriteString(jsRef)
			bodyBuilder.WriteString(endJs)
		}
		bodyBuilder.WriteString(part2)

		bodyLines, err = computeBody(bodyLines)
		if err != nil {
			return err
		}

		for _, line := range bodyLines {
			bodyBuilder.WriteString(line)
			bodyBuilder.WriteByte('\n')
		}
		bodyBuilder.WriteString(part3)

		return os.WriteFile(destPath, []byte(bodyBuilder.String()), 0644)
	})
}

func parseHtmlFragment(path string) ([]string, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var jsRefs []string
	var bodyLines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if trimmed := strings.TrimSpace(scanner.Text()); trimmed != "" {
			if trimmed == partSep {
				break
			}
			jsRefs = append(jsRefs, trimmed)
		}
	}
	for scanner.Scan() {
		if trimmed := strings.TrimSpace(scanner.Text()); trimmed != "" {
			bodyLines = append(bodyLines, trimmed)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, nil, err
	}

	if len(bodyLines) == 0 {
		bodyLines = jsRefs
		jsRefs = nil
	}
	return jsRefs, bodyLines, nil
}

func makeDirectory(path string, nameSize int) error {
	i := len(path) - nameSize
	return os.MkdirAll(path[:i], 0755)
}

func noCompute(bodyLines []string) ([]string, error) {
	// nothing to do
	return bodyLines, nil
}

type markdownEngine struct {
	md goldmark.Markdown
}

func newMarkdownEngine() markdownEngine {
	return markdownEngine{md: markdownextension.NewDefaultEngine()}
}

func (e markdownEngine) markdownCompute(bodyLines []string) ([]string, error) {
	var paramBuilder strings.Builder
	for _, line := range bodyLines {
		paramBuilder.WriteString(line)
		paramBuilder.WriteByte('\n')
	}

	var resBuilder strings.Builder
	if err := e.md.Convert([]byte(paramBuilder.String()), &resBuilder); err != nil {
		return nil, err
	}
	return []string{resBuilder.String()}, nil
}
