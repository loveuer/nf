package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/loveuer/nf/nft/ntl/internal/opt"
	"github.com/loveuer/nf/nft/ntl/pkg/loading"
	"github.com/loveuer/nf/nft/tool"
	"github.com/loveuer/nf/pkg/log"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:     "new",
	Short:   "new a nf project",
	Example: "nfctl new <project> -t ultone [options]",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("必须提供 project 名称")
		}

		if strings.HasSuffix(args[0], "/") {
			return errors.New("project 名称不能以 / 结尾")
		}

		base := path.Base(args[0])
		if strings.HasPrefix(base, ".") {
			return errors.New("project 名称不能以 . 开头")
		}

		ch := make(chan *loading.Loading)
		defer close(ch)

		go loading.Print(cmd.Context(), ch)
		ch <- &loading.Loading{Content: "开始新建项目: " + args[0], Type: loading.TypeInfo}

		pwd, err := os.Getwd()
		if err != nil {
			ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeError}
			return err
		}

		moduleName := args[0]
		pwd = path.Join(filepath.ToSlash(pwd), base)

		log.Debug("cmd.new: new project, pwd = %s, name = %s, template = %s", pwd, moduleName, opt.Cfg.New.Template)

		ch <- &loading.Loading{Content: "开始下载模板: " + opt.Cfg.New.Template, Type: loading.TypeProcessing}

		repo := opt.Cfg.New.Template
		if v, ok := opt.TemplateMap[repo]; ok {
			repo = v
		}

		if err = tool.Clone(pwd, repo); err != nil {
			ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeError}
			return err
		}

		ch <- &loading.Loading{Content: "下载完成: " + opt.Cfg.New.Template, Type: loading.TypeSuccess}

		if err = os.RemoveAll(path.Join(pwd, ".git")); err != nil {
			ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeWarning}
		}

		if opt.Cfg.New.DisableInitScript {
			ch <- &loading.Loading{Content: fmt.Sprintf("创建项目 %s 成功", args[0]), Type: loading.TypeSuccess}
			return nil
		}

		var info os.FileInfo
		if info, err = os.Stat(path.Join(pwd, ".nfctl")); err != nil {
			log.Debug("cmd.new: stat .nfctl err, err = %v", err)
			if errors.Is(err, os.ErrNotExist) {
				ch <- &loading.Loading{Content: fmt.Sprintf("创建项目 %s 成功", args[0]), Type: loading.TypeSuccess}
				return nil
			}

			ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeWarning}
			ch <- &loading.Loading{Content: fmt.Sprintf("创建项目 %s 成功", args[0]), Type: loading.TypeSuccess}

			return nil
		}

		if info.IsDir() {
			ch <- &loading.Loading{Content: "错误的初始化脚本(is_dir)", Type: loading.TypeWarning}
			ch <- &loading.Loading{Content: fmt.Sprintf("创建项目 %s 成功", args[0]), Type: loading.TypeSuccess}
			return nil
		}

		ch <- &loading.Loading{Content: "开始初始化项目: " + args[0], Type: loading.TypeProcessing}

		return nil
	},
	SilenceErrors: true,
}

func initNew() *cobra.Command {
	newCmd.Flags().StringVarP(&opt.Cfg.New.Template, "template", "t", "ultone", "template name/url[example:ultone, https://gitea.loveuer.com/loveuer/ultone.git]")
	newCmd.Flags().BoolVar(&opt.Cfg.New.DisableInitScript, "disable-init-script", false, "disable init script(.nfctl)")
	return newCmd
}
