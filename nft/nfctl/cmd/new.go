package cmd

import (
	"errors"
	"fmt"
	"github.com/loveuer/nf/nft/log"
	"github.com/loveuer/nf/nft/nfctl/clone"
	"github.com/loveuer/nf/nft/nfctl/opt"
	"github.com/loveuer/nf/nft/nfctl/tp"
	"github.com/loveuer/nf/nft/nfctl/version"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"path"
)

var (
	cmdNew = &cobra.Command{
		Use:   "new",
		Short: "nfctl new: start new project",
		Example: `nfctl new {project} -t ultone [recommend]
nfctl new {project} -t https://github.com/loveuer/ultone.git
nfctl new {project} --template http://username:token@my.gitlab.com/my-zone/my-repo.git
`,
		SilenceUsage: true,
	}

	template    string
	disableInit bool

	preTemplateMap = map[string]string{
		"ultone": "https://gitcode.com/loveuer/ultone.git",
	}
)

func initNew() {
	cmdNew.Flags().StringVarP(&template, "template", "t", "", "template name/url[example:ultone, https://github.com/xxx/yyy.git]")
	cmdNew.Flags().BoolVar(&disableInit, "without-init", false, "don't run template init script")

	cmdNew.RunE = func(cmd *cobra.Command, args []string) error {
		version.Check(true, false, 5)

		var (
			err        error
			urlIns     *url.URL
			pwd        string
			projectDir string
			initBs     []byte
			renderBs   []byte
			scripts    []tp.Cmd
		)

		if len(args) == 0 {
			return fmt.Errorf("project name required")
		}

		if pwd, err = os.Getwd(); err != nil {
			return fmt.Errorf("get work dir err")
		}

		projectDir = path.Join(pwd, args[0])

		if _, err = os.Stat(projectDir); !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("project folder already exist")
		}

		if err = os.MkdirAll(projectDir, 0750); err != nil {
			return fmt.Errorf("create project dir err: %v", err)
		}

		defer func() {
			if err != nil {
				_ = os.RemoveAll(projectDir)
			}
		}()

		if template == "" {
			// todo no template new project
			return fmt.Errorf("😥create basic project(without template) comming soon...")
		}

		cloneUrl := template
		if ptUrl, ok := preTemplateMap[cloneUrl]; ok {
			cloneUrl = ptUrl
		}

		if urlIns, err = url.Parse(cloneUrl); err != nil {
			return fmt.Errorf("invalid clone url: %v", err)
		}

		if err = clone.Clone(projectDir, urlIns); err != nil {
			return fmt.Errorf("clone template err: %v", err)
		}

		if initBs, err = os.ReadFile(path.Join(projectDir, ".nfctl")); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}

			return fmt.Errorf("read nfctl script file err: %v", err)
		}

		if renderBs, err = tp.RenderVar(initBs, map[string]any{
			"PROJECT_NAME": args[0],
		}); err != nil {
			return fmt.Errorf("render template init script err: %v", err)
		}

		if scripts, err = tp.ParseCmd(projectDir, renderBs); err != nil {
			return fmt.Errorf("parse template init script err: %v", err)
		}

		for _, script := range scripts {
			if opt.Debug {
				log.Debug("start script:\n%s\n", script.String())
			}

			if err = script.Execute(); err != nil {
				return fmt.Errorf("execute template init script err: %v", err)
			}
		}

		if err = os.RemoveAll(path.Join(projectDir, ".git")); err != nil {
			log.Warn("remove .git folder err: %s", err.Error())
		}

		log.Info("🎉 create project [%s] 成功!!!", args[0])

		return nil
	}
}
