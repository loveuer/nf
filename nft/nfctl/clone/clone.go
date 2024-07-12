package clone

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	_ "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/loveuer/nf/nft/log"
	"net/url"
)

func Clone(pwd string, ins *url.URL) error {
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

	log.Info("start clone %s", uri)
	_, err := git.PlainClone(pwd, false, opt)
	if err != nil {
		return err
	}

	return nil
}
