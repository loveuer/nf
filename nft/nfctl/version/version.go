package version

import (
	"crypto/tls"
	"fmt"
	"github.com/fatih/color"
	"github.com/loveuer/nf/nft/log"
	"github.com/savioxavier/termlink"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	uri    = "https://raw.gitcode.com/loveuer/nf/raw/master/nft/nfctl/version/var.go"
	prefix = "const Version = "
)

func UpgradePrint(newVersion string) {
	fmt.Printf(`+----------------------------------------------------------------------+
|          🎉 🎉 🎉 %s 🎉 🎉 🎉          |
| %s |
| Or Download by:                                                      |
| %s                    |
| %s                   |
+----------------------------------------------------------------------+
`,
		color.GreenString("New Version Found: %s", newVersion),
		color.CyanString("Upgrade it with: [go install github.com/loveuer/nf/nft/nfctl@master]"),
		color.CyanString(termlink.Link("Releases", "https://github.com/loveuer/nf/releases")),
		color.CyanString(termlink.Link("Releases", "https://gitcode.com/loveuer/nf/releases")),
	)
}

func Check(printUpgradable bool, printNoNeedUpgrade bool, timeouts ...int) string {
	var (
		v string
	)

	defer func() {
		if printUpgradable {
			if v > Version {
				UpgradePrint(v)
			}
		}

		if printNoNeedUpgrade {
			if v == Version {
				color.Cyan("Your Version: %s is Newest", Version)
			}
		}
	}()

	timeout := time.Duration(30) * time.Second
	if len(timeouts) > 0 && timeouts[0] > 0 {
		timeout = time.Duration(timeouts[0]) * time.Second
	}

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	resp, err := (&http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}).Do(req)
	if err != nil {
		log.Debug("[Check] http get[%s] err: %v", uri, err.Error())
		return ""
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Debug("[Check] http read all body err: %v", err)
	}

	log.Debug("[Check] http get[%s] body:\n%s", uri, string(content))

	for _, line := range strings.Split(string(content), "\n") {
		log.Debug("[Check] version.go line: %s", line)
		if strings.HasPrefix(line, prefix) {
			may := strings.TrimPrefix(line, prefix)
			if len(may) > 2 {
				v = may[1 : len(may)-1]
			}

			return v
		}
	}

	return ""
}
