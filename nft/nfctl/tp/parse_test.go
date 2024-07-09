package tp

import (
	"github.com/loveuer/nf/nft/log"
	"os"
	"testing"
)

func TestParseInitFile(t *testing.T) {
	const init_bs = `
!replace
content
reg
*.go
ultone => {{.PROJECT_NAME}}
EOF

!replace
content
exact
go.mod
module ultone => module {{.PROJECT_NAME}}
EOF

!replace
name
main.go => loveuer.go
EOF

!generate
readme.md
# {{.PROJECT_NAME}}

### run
- ` + "`" + `go run . --help` + "`" + `
- ` + "`" + `go run .` + "`" + `

### build
- ` + "`" + `docker build -t {repo:tag} -f Dockerfile .` + "`" + `
EOF
`
	data := map[string]any{
		"PROJECT_NAME": "loveuer",
	}

	result, err := RenderVar([]byte(init_bs), data)
	if err != nil {
		log.Fatal(err.Error())
	}

	pwd, _ := os.Getwd()

	cmds, err := ParseCmd(pwd, result)
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, item := range cmds {
		log.Info("one cmd => %s\n\n", item.String())
	}
}
