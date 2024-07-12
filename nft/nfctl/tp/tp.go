package tp

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/loveuer/nf/nft/log"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type Cmd interface {
	String() string
	Execute() error
}

var (
	_ Cmd = (*Generate)(nil)
	_ Cmd = (*ReplaceContent)(nil)
	_ Cmd = (*ReplaceName)(nil)
)

type Generate struct {
	pwd      string
	filename string
	content  []string
}

func (t *Generate) String() string {
	return fmt.Sprintf("!generate\n%s\n%s\n", t.filename, strings.Join(t.content, "\n"))
}

func (t *Generate) Execute() error {
	var (
		err      error
		location = t.filename
		input    *os.File
	)

	log.Debug("[Generate] generate[%s]", t.filename)

	if !path.IsAbs(t.filename) {
		location = path.Join(t.pwd, t.filename)
	}

	if err = os.MkdirAll(path.Dir(location), 0644); err != nil {
		return err
	}

	if !strings.HasSuffix(location, "/") {
		if input, err = os.OpenFile(location, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0744); err != nil {
			return err
		}

		if len(t.content) > 0 {
			content := strings.Join(t.content, "\n")
			_, err = input.WriteString(content)
			return err
		}
	}

	return nil
}

func newGenerate(pwd string, lines []string) (*Generate, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("generate cmd require file/folder name")
	}

	return &Generate{
		pwd:      pwd,
		filename: lines[0],
		content:  lines[1:],
	}, nil
}

type replaceNameMatchType int

const (
	replaceNameMatchReg replaceNameMatchType = iota + 1
	replaceNameMatchExact
	replaceNameMatchPrefix
	replaceNameMatchSuffix
)

func (rm replaceNameMatchType) Label() string {
	switch rm {
	case replaceNameMatchReg:
		return "reg"
	case replaceNameMatchExact:
		return "exact"
	case replaceNameMatchPrefix:
		return "prefix"
	case replaceNameMatchSuffix:
		return "suffix"
	}

	log.Panic("unknown replace match type: %v", rm)
	return ""
}

type ReplaceContent struct {
	pwd     string
	name    string
	content string

	targetName    string
	matchType     replaceNameMatchType
	fromContent   string
	targetEmpty   bool
	targetContent string
}

func (t *ReplaceContent) String() string {
	return fmt.Sprintf("!replace content\n%s\n%s\n", t.name, t.content)
}

func (t *ReplaceContent) Execute() error {
	var (
		fn filepath.WalkFunc

		handler = func(location string) error {
			bs, err := os.ReadFile(location)
			if err != nil {
				return err
			}

			log.Debug("[ReplaceContent] handle[%s] replace [%s] => [%s]", location, t.fromContent, t.targetContent)
			newbs, err := t.executeFile(bs)
			if err != nil {
				return err
			}

			return os.WriteFile(location, newbs, 0644)
		}
	)

	switch t.matchType {
	case replaceNameMatchExact:
		fn = func(location string, info fs.FileInfo, err error) error {
			if location == path.Join(t.pwd, t.targetName) {
				log.Debug("[ReplaceContent] exact match: %s", location)
				return handler(location)
			}

			return nil
		}
	case replaceNameMatchPrefix:
		fn = func(location string, info fs.FileInfo, err error) error {
			if strings.HasPrefix(path.Base(location), t.targetName) {
				log.Debug("[ReplaceContent] prefix match: %s", location)
				return handler(location)
			}

			return nil
		}
	case replaceNameMatchSuffix:
		fn = func(location string, info fs.FileInfo, err error) error {
			if strings.HasSuffix(location, t.targetName) {
				log.Debug("[ReplaceContent] suffix match: %s", location)
				return handler(location)
			}

			return nil
		}
	case replaceNameMatchReg:
		fn = func(location string, info fs.FileInfo, err error) error {
			if match, err := regexp.MatchString(t.targetName, location); err == nil && match {
				log.Debug("[ReplaceContent] reg match: %s", location)
				return handler(location)
			}

			return nil
		}
	}

	return filepath.Walk(t.pwd, fn)
}

func (t *ReplaceContent) executeFile(raw []byte) ([]byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	scanner.Buffer(make([]byte, 1024), 1024*1024)

	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(
			lines,
			strings.ReplaceAll(line, t.fromContent, t.targetContent),
		)
	}

	return []byte(strings.Join(lines, "\n")), nil
}

func newReplaceContent(pwd string, lines []string) (*ReplaceContent, error) {
	if len(lines) != 2 {
		return nil, fmt.Errorf("invalid replace_content cmd: required 2 lines params")
	}

	var (
		name      = lines[0]
		content   = lines[1]
		matchType replaceNameMatchType
	)

	names := strings.SplitN(name, " ", 2)
	if len(names) != 2 {
		return nil, fmt.Errorf("invalid replace_content cmd: name line, required: [reg/exact/prefix/shuffix] {filename}")
	}

	switch names[0] {
	case "exact":
		matchType = replaceNameMatchExact
	case "reg":
		matchType = replaceNameMatchReg
	case "prefix":
		matchType = replaceNameMatchPrefix
	case "suffix":
		matchType = replaceNameMatchSuffix
	default:
		return nil, fmt.Errorf("invalid replace_content name match type, example: [reg *.go] [exact go.mod]")
	}

	var (
		targetName    string = names[1]
		targetEmpty          = false
		targetContent string
	)
	contents := strings.SplitN(content, "=>", 2)
	fromContent := strings.TrimSpace(contents[0])
	if len(contents) == 1 {
		targetEmpty = true
	} else {
		if targetContent = strings.TrimSpace(contents[1]); targetContent == "" || targetContent == `""` || targetContent == `''` {
			targetEmpty = true
		}
	}

	return &ReplaceContent{
		pwd:     pwd,
		name:    name,
		content: content,

		matchType:     matchType,
		targetName:    targetName,
		fromContent:   fromContent,
		targetEmpty:   targetEmpty,
		targetContent: targetContent,
	}, nil
}

type ReplaceName struct {
	pwd  string
	line string

	targetEmpty   bool
	fromContent   string
	targetContent string
}

func (t *ReplaceName) String() string {
	return fmt.Sprintf("!replace name\n%s\n", t.line)
}

func (t *ReplaceName) Execute() error {
	fullpath := path.Join(t.pwd, t.fromContent)
	if t.targetEmpty {
		return os.RemoveAll(fullpath)
	}

	ftpath := path.Join(t.pwd, t.targetContent)
	return os.Rename(fullpath, ftpath)
}

func newReplaceName(pwd string, lines []string) (*ReplaceName, error) {
	if len(lines) != 1 {
		return nil, fmt.Errorf("replace_name need one line param, for example: mian.go => main.go")
	}

	var (
		content       = lines[0]
		targetEmpty   = false
		fromContent   string
		targetContent string
	)

	contents := strings.SplitN(content, "=>", 2)
	fromContent = strings.TrimSpace(contents[0])
	if len(contents) == 1 {
		targetEmpty = true
	} else {
		if targetContent = strings.TrimSpace(contents[1]); targetContent == "" || targetContent == `""` || targetContent == `''` {
			targetEmpty = true
		}
	}

	if !targetEmpty {
		if (strings.HasPrefix(targetContent, `"`) && strings.HasSuffix(targetContent, `"`)) || (strings.HasPrefix(targetContent, `'`) && strings.HasSuffix(targetContent, `'`)) {
			targetContent = targetContent[1 : len(targetContent)-1]
		}
	}

	return &ReplaceName{
		pwd:  pwd,
		line: content,

		targetEmpty:   targetEmpty,
		fromContent:   fromContent,
		targetContent: targetContent,
	}, nil
}
