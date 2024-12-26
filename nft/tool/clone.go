package tool

import (
	"fmt"
	"net/url"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func Clone(projectDir string, repoURL string) error {
	ins, err := url.Parse(repoURL)
	if err != nil {
		return err
	}

	uri := fmt.Sprintf("%s://%s%s", ins.Scheme, ins.Host, ins.Path)

	opt := &git.CloneOptions{
		URL:             uri,
		Depth:           1,
		InsecureSkipTLS: true,
		SingleBranch:    true,
	}

	if ins.User != nil {
		password, _ := ins.User.Password()
		opt.Auth = &http.BasicAuth{
			Username: ins.User.Username(),
			Password: password,
		}
	}

	_, err = git.PlainClone(projectDir, false, opt)
	if err != nil {
		return err
	}

	return nil
}
