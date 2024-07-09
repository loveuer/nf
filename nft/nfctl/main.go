package main

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	_ "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/loveuer/nf/nft/log"
	"io"
)

func main() {
	memo := memory.NewStorage()
	fs := memfs.New()
	repo, err := git.Clone(memo, fs, &git.CloneOptions{
		URL:             "http://10.220.10.35/dev/template/ultone.git",
		Auth:            &http.BasicAuth{Username: "loveuer", Password: "uu_L6neSDseoWx55babJ"},
		Depth:           1,
		SingleBranch:    true,
		InsecureSkipTLS: true,
	})
	if err != nil {
		panic(err)
	}

	infos, err := fs.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, item := range infos {
		log.Info("[fs.info] %s", item.Name())
		if item.Name() == "main.go" {
			file, err := fs.Open(item.Name())
			if err != nil {
				panic(err)
			}

			bs, err := io.ReadAll(file)
			if err != nil {
				panic(err)
			}

			log.Info("[fs.main]\n%s", string(bs))
		}
	}

	_ = repo
}
