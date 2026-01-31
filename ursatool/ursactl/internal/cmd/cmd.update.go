package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"regexp"
	"strings"
	"time"

	resty "github.com/go-resty/resty/v2"
	"github.com/loveuer/ursa/ursatool/loading"
	"github.com/loveuer/ursa/ursatool/log"
	"github.com/loveuer/ursa/ursatool/ursactl/internal/opt"
	"github.com/loveuer/ursa/ursatool/tool"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update ursactl self",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func initUpdate() *cobra.Command {
	return updateCmd
}

func doUpdate(ctx context.Context) (err error) {
	ctxWithTimeout, cancel := tool.TimeoutCtx(ctx, 30)
	defer cancel()
	return loading.Do(ctxWithTimeout, func(ctx context.Context, print func(msg string, types ...loading.Type)) error {
		print("正在检查更新...")
		tip := "❗ 请尝试手动更新: go install github.com/loveuer/ursa/ursatool/ursactl@master"
		version := ""

		var rr *resty.Response
		if rr, err = resty.New().SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).R().
			SetContext(ctx).
			Get(opt.VersionURL); err != nil {
			err = fmt.Errorf("检查更新失败: %s\n%s", err.Error(), tip)
			return err
		}

		log.Debug("cmd.update: url = %s, raw_response = %s", opt.VersionURL, rr.String())

		if rr.StatusCode() != 200 {
			err = fmt.Errorf("检查更新失败: %s\n%s", rr.Status(), tip)
			return err
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
			return err
		}

		log.Debug("cmd.update: find version = %s, now_version = %s", version, opt.Version)

		if version <= opt.Version {
			print(fmt.Sprintf("已是最新版本: %s", opt.Version), loading.TypeSuccess)
			return nil
		}

		print(fmt.Sprintf("发现新版本: %s", version), loading.TypeInfo)

		print(fmt.Sprintf("正在更新到 %s ...", version))

		time.Sleep(2 * time.Second)

		print("暂时无法自动更新, 请尝试手动更新: go install github.com/loveuer/ursa/ursatool/ursactl@master", loading.TypeWarning)

		return nil
	})
}
