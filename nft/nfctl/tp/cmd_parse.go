package tp

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func ParseCmd(pwd string, content []byte) ([]TpCmd, error) {
	var (
		err   error
		cmds  = make([]TpCmd, 0)
		start = false
	)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 1024), 1024*1024*10)

	record := make([]string, 0)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "!") {
			if start {
				return nil, fmt.Errorf("invalid content: unEOF cmd block found")
			}

			start = true
			record = append(record, line)
			continue
		}

		if strings.HasPrefix(line, "EOF") {
			start = false
			if len(record) == 0 {
				continue
			}

			var cmd TpCmd
			if cmd, err = ParseBlock(pwd, record); err != nil {
				return nil, err
			}

			cmds = append(cmds, cmd)
			record = record[:0]
			continue
		}

		if start {
			record = append(record, line)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return cmds, err
}

func ParseBlock(pwd string, lines []string) (TpCmd, error) {
	switch lines[0] {
	case "!replace":
		return newReplace(pwd, lines[1:])
	case "!generate":
		return newGenerate(pwd, lines[1:])
	}

	return nil, fmt.Errorf("invalid cmd block: unknown type: %s", lines[0])
}
