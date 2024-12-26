package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/loveuer/nf/nft/log"
	"github.com/loveuer/nf/nft/nfctl/internal/opt"
	"github.com/loveuer/nf/nft/nfctl/pkg/loading"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update nfctl self",
	RunE:  func(cmd *cobra.Command, args []string) error { return nil },
}

func initUpdate() *cobra.Command {
	return updateCmd
}

func doUpdate(ctx context.Context) (err error) {
	ch := make(chan *loading.Loading)
	defer close(ch)

	go func() {
		loading.Print(ctx, ch)
	}()

	ch <- &loading.Loading{Content: "正在检查更新...", Type: loading.TypeProcessing}
	tip := "❗ 请尝试手动更新: go install github.com/loveuer/nf/nft/nfctl@latest"
	version := ""

	var rr *resty.Response
	if rr, err = resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
		SetContext(ctx).
		Get(opt.VersionURL); err != nil {
		err = fmt.Errorf("检查更新失败: %s\n%s", err.Error(), tip)
		ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeError}
		return err
	}

	log.Debug("cmd.update: url = %s, raw_response = %s", opt.VersionURL, rr.String())

	if rr.StatusCode() != 200 {
		err = fmt.Errorf("检查更新失败: %s\n%s", rr.Status(), tip)
		ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeError}
		return
	}

	reg := regexp.MustCompile(`const Version = "v\d{2}\.\d{2}\.\d{2}-r\d{1,2}"`)
	for _, line := range strings.Split(rr.String(), "\n") {
		if reg.MatchString(line) {
			version = strings.TrimSpace(strings.TrimPrefix(line, "const Version = "))
			version = version[1 : len(version)-1]
			break
		}
	}

	if version == "" {
		err = fmt.Errorf("检查更新失败: 未找到版本信息\n%s", tip)
		ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeError}
		return err
	}

	log.Debug("cmd.update: find version = %s, now_version = %s", version, opt.Version)

	if version <= opt.Version {
		ch <- &loading.Loading{Content: fmt.Sprintf("已是最新版本: %s", opt.Version), Type: loading.TypeSuccess}
		return nil
	}

	ch <- &loading.Loading{Content: fmt.Sprintf("发现新版本: %s", version), Type: loading.TypeInfo}

	ch <- &loading.Loading{Content: fmt.Sprintf("正在更新到 %s ...", version)}

	time.Sleep(2 * time.Second)
	ch <- &loading.Loading{Content: "暂时无法自动更新, 请尝试手动更新: go install github.com/loveuer/nf/nft/nfctl@latest", Type: loading.TypeWarning}
	return nil
}
