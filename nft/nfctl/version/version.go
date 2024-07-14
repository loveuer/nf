package version

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/loveuer/nf/nft/log"
	"github.com/savioxavier/termlink"
	"strings"
	"time"
)

var (
	client = resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	uri    = "https://raw.gitcode.com/loveuer/nf/raw/master/nft/nfctl/version/var.go"
	prefix = "const Version = "
)

func UpgradePrint(newVersion string) {
	t := table.NewWriter()
	t.AppendRows([]table.Row{
		{color.GreenString("New Version Found: %s", newVersion)},
		{color.CyanString("Upgrade it with: [go install github.com/loveuer/nf/nft/nfctl@master]")},
		{fmt.Sprint("Or Download by: ")},
		{color.CyanString(termlink.Link("Releases", "https://github.com/loveuer/nf/releases"))},
		{color.CyanString(termlink.Link("Releases", "https://gitcode.com/loveuer/nf/releases"))},
	})

	fmt.Println(t.Render())
}

func Check(printUpgradable bool, printNoNeedUpgrade bool, timeout ...int) string {
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

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	if len(timeout) > 0 && timeout[0] > 0 {
		ctx, _ = context.WithTimeout(context.Background(), time.Duration(timeout[0])*time.Second)
	}

	resp, err := client.R().SetContext(ctx).
		Get(uri)
	if err != nil {
		log.Debug("[Check] http get[%s] err: %v", uri, err.Error())
		return ""
	}

	log.Debug("[Check] http get[%s] body:\n%s", uri, resp.String())

	for _, line := range strings.Split(resp.String(), "\n") {
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
