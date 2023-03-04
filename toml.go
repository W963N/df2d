package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/naoina/toml"
	"golang.org/x/sync/errgroup"
)

type envInfo struct {
	Root                   string
	Excluded_ext           []string
	IsExcluded_hidden_file bool
}

func (cfg *envInfo) loadConf(file io.Reader) error {
	return toml.NewDecoder(file).Decode(&cfg)
}

func (cfg *envInfo) getRoot() string {
	return cfg.Root
}

func (cfg *envInfo) getExcludedExt() []string {
	return cfg.Excluded_ext
}

func (cfg *envInfo) getIsExcludedHiddenFile() bool {
	return cfg.IsExcluded_hidden_file
}

func (cfg *envInfo) chkFormat() {
	eg := new(errgroup.Group)

	eg.Go(func() error {
		if !filepath.IsAbs(cfg.getRoot()) {
			return errors.New("Abs is false")
		}
		return nil
	})

	eg.Go(func() error {
		if f, err := os.Stat(cfg.getRoot()); os.IsNotExist(err) || !f.IsDir() {
			return errors.New("Don't Exist file")
		}
		return nil
	})

	eg.Go(func() error {
		for _, ext := range cfg.getExcludedExt() {
			if filepath.Ext(ext) == "" {
				return errors.New("Don't Exist Ext: " + ext)
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		die("chk envInfo: %v\n", err)
	}
}
