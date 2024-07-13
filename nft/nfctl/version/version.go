package version

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/loveuer/nf/nft/log"
	"github.com/savioxavier/termlink"
	"net/http"
	"strings"
	"sync"
)

const Version = "v24.07.13-r1"

var (
	lk      = &sync.Mutex{}
	empty   = func() {}
	upgrade = func(v string) func() {
		return func() {
			color.Green("\n🎉 🎉 🎉 [nfctl] New Version Found: %s", v)
			color.Cyan("Upgrade it with: [go install github.com/loveuer/nf/nft/nfctl@master]")
			fmt.Print("Or Download by: ")
			color.Cyan(termlink.Link("Releases", "https://github.com/loveuer/nf/releases"))
			fmt.Println()
		}
	}
	Fn   = empty
	OkCh = make(chan struct{}, 1)
)

func Check() {
	ready := make(chan struct{})
	go func() {
		ready <- struct{}{}
		uri := "https://raw.gitcode.com/loveuer/nf/raw/master/nft/nfctl/version/version.go"
		prefix := "const Version = "
		resp, err := http.Get(uri)
		if err != nil {
			log.Debug("[Check] http get[%s] err: %v", uri, err.Error())
			return
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 16*1024), 1024*1024)

		for scanner.Scan() {
			line := scanner.Text()
			log.Debug("[Check] version.go line: %s", line)
			if strings.HasPrefix(line, prefix) {
				v := strings.TrimPrefix(line, prefix)
				if len(v) > 2 {
					v = v[1 : len(v)-1]
				}

				if v != "" && v > Version {
					lk.Lock()
					Fn = upgrade(v)
					lk.Unlock()
					OkCh <- struct{}{}
					return
				}
			}
		}
	}()
	<-ready
}
