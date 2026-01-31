package cmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/loveuer/ursa/ursatool/loading"
	"github.com/loveuer/ursa/ursatool/log"
	"github.com/loveuer/ursa/ursatool/ursactl/internal/opt"
	"github.com/loveuer/ursa/ursatool/tool"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:           "new",
	Short:         "new a nf project",
	Example:       "ursactl new <project> -t ultone [options]",
	RunE:          doNew,
	SilenceErrors: true,
}

func initNew() *cobra.Command {
	newCmd.Flags().StringVarP(&opt.Cfg.New.Template, "template", "t", "ultone", "template name/url[example:ultone, https://gitea.loveuer.com/loveuer/ultone.git]")
	newCmd.Flags().BoolVar(&opt.Cfg.New.DisableInitScript, "disable-init-script", false, "disable init script(.ursactl)")
	return newCmd
}

func doNew(cmd *cobra.Command, args []string) (err error) {
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

	return loading.Do(cmd.Context(), func(ctx context.Context, print func(msg string, types ...loading.Type)) error {
		print("开始新建项目: "+args[0], loading.TypeInfo)

		pwd, err := os.Getwd()
		if err != nil {
			return err
		}

		moduleName := args[0]
		pwd = path.Join(filepath.ToSlash(pwd), base)

		log.Debug("cmd.new: new project, pwd = %s, name = %s, template = %s", pwd, moduleName, opt.Cfg.New.Template)

		print("开始下载模板: "+opt.Cfg.New.Template, loading.TypeProcessing)

		repo := opt.Cfg.New.Template
		if v, ok := opt.TemplateMap[repo]; ok {
			repo = v
		}

		if err = tool.Clone(pwd, repo); err != nil {
			return err
		}

		print("下载模板完成: "+opt.Cfg.New.Template, loading.TypeSuccess)

		if err = os.RemoveAll(path.Join(pwd, ".git")); err != nil {
			print(err.Error(), loading.TypeWarning)
		}

		print("开始初始化项目: "+args[0], loading.TypeProcessing)

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
					print("初始化文件失败: "+err.Error(), loading.TypeWarning)
					print("开始初始化项目: "+args[0], loading.TypeProcessing)
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
			return err
		}

		var (
			render *template.Template
			rf     *os.File
		)

		if render, err = template.New(base).Parse(opt.README); err != nil {
			log.Debug("cmd.new: new text template err, err = %s", err.Error())
			print("生成 readme 失败", loading.TypeWarning)
			goto END
		}

		if rf, err = os.OpenFile(path.Join(pwd, "readme.md"), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0o644); err != nil {
			log.Debug("cmd.new: new readme file err, err = %s", err.Error())
			print("生成 readme 失败", loading.TypeWarning)
			goto END
		}
		defer rf.Close()

		if err = render.Execute(rf, map[string]any{
			"project_name": base,
		}); err != nil {
			log.Debug("cmd.new: template execute err, err = %s", err.Error())
			print("生成 readme 失败", loading.TypeWarning)
		}

	END:
		print(fmt.Sprintf("项目: %s 初始化成功", args[0]), loading.TypeSuccess)

		return nil
	})
}
