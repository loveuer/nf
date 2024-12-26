package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/loveuer/nf/nft/log"
	"github.com/loveuer/nf/nft/nfctl/internal/opt"
	"github.com/loveuer/nf/nft/nfctl/pkg/loading"
	"github.com/loveuer/nf/nft/tool"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:           "new",
	Short:         "new a nf project",
	Example:       "nfctl new <project> -t ultone [options]",
	RunE:          doNew,
	SilenceErrors: true,
}

func initNew() *cobra.Command {
	newCmd.Flags().StringVarP(&opt.Cfg.New.Template, "template", "t", "ultone", "template name/url[example:ultone, https://gitea.loveuer.com/loveuer/ultone.git]")
	newCmd.Flags().BoolVar(&opt.Cfg.New.DisableInitScript, "disable-init-script", false, "disable init script(.nfctl)")
	return newCmd
}

func doNew(cmd *cobra.Command, args []string) error {
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

	ch <- &loading.Loading{Content: "下载模板完成: " + opt.Cfg.New.Template, Type: loading.TypeSuccess}

	if err = os.RemoveAll(path.Join(pwd, ".git")); err != nil {
		ch <- &loading.Loading{Content: err.Error(), Type: loading.TypeWarning}
	}

	ch <- &loading.Loading{Content: "开始初始化项目: " + args[0], Type: loading.TypeProcessing}

	if err = filepath.Walk(pwd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "go.mod") {
			var content []byte
			if content, err = os.ReadFile(path); err != nil {
				ch <- &loading.Loading{Content: "初始化文件失败: " + err.Error(), Type: loading.TypeWarning}
				ch <- &loading.Loading{Content: "开始初始化项目: " + args[0], Type: loading.TypeProcessing}
				return nil
			}

			scanner := bufio.NewScanner(bytes.NewReader(content))
			replaced := make([]string, 0, 16)
			for scanner.Scan() {
				line := scanner.Text()
				// 操作 go.mod 文件时, 忽略 toolchain 行, 以更好的兼容 go1.20
				if strings.HasSuffix(path, "go.mod") && strings.HasPrefix(line, "toolchain") {
					continue
				}
				replaced = append(replaced, strings.ReplaceAll(line, opt.Cfg.New.Template, moduleName))
			}
			if err = os.WriteFile(path, []byte(strings.Join(replaced, "\n")), 0o644); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		ch <- &loading.Loading{Content: "初始化文件失败: " + err.Error(), Type: loading.TypeWarning}
		return err
	}

	var (
		render *template.Template
		rf     *os.File
	)

	if render, err = template.New(base).Parse(opt.README); err != nil {
		log.Debug("cmd.new: new text template err, err = %s", err.Error())
		ch <- &loading.Loading{Content: "生成 readme 失败", Type: loading.TypeWarning}
		goto END
	}

	if rf, err = os.OpenFile(path.Join(pwd, "readme.md"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o644); err != nil {
		log.Debug("cmd.new: new readme file err, err = %s", err.Error())
		ch <- &loading.Loading{Content: "生成 readme 失败", Type: loading.TypeWarning}
		goto END
	}
	defer rf.Close()

	if err = render.Execute(rf, map[string]any{
		"project_name": base,
	}); err != nil {
		log.Debug("cmd.new: template execute err, err = %s", err.Error())
		ch <- &loading.Loading{Content: "生成 readme 失败", Type: loading.TypeWarning}
	}

END:
	ch <- &loading.Loading{Content: fmt.Sprintf("项目: %s 初始化成功", args[0]), Type: loading.TypeSuccess}

	return nil
}
