package tp

import (
	"github.com/loveuer/nf/nft/log"
	"os"
	"testing"
)

func TestParseInitFile(t *testing.T) {
	bs, err := os.ReadFile("xtest")
	if err != nil {
		log.Fatal(err.Error())
	}

	data := map[string]any{
		"PROJECT_NAME": "myproject",
	}

	result, err := RenderVar(bs, data)
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
		if err = item.Execute(); err != nil {
			log.Fatal(err.Error())
		}
	}
}
