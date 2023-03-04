package main

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/naoina/toml"
	"golang.org/x/sync/errgroup"
)

const envfile string = "df2d.toml"

type cpEnv struct {
	path string
	tml  envToml
	stat *syscall.Stat_t
}

type envToml struct {
	IsMkdir, IsOverwrite, IsSymlink bool
	Dest_dir                        string
}

func (cfg *envToml) loadConf(file io.Reader) error {
	return toml.NewDecoder(file).Decode(&cfg)
}

func getEnv(envs *[]cpEnv, dir string, hidden bool) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	eg := new(errgroup.Group)
	for _, file := range files {
		if hidden && strings.HasPrefix(filepath.Base(file.Name()), ".") {
			continue
		}

		file := file
		eg.Go(func() error {
			path := filepath.Join(dir, file.Name())
			err := envWalk(envs, path, file, hidden)
			if err != nil {
				return err
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}

func envWalk(envs *[]cpEnv, path string, file fs.DirEntry, hidden bool) error {
	if file.IsDir() {
		if err := getEnv(envs, path, hidden); err != nil {
			return err
		}
	} else {
		if file.Name() == envfile {
			if err := setMvEnv(envs, path, file); err != nil {
				return err
			}
		}
	}

	return nil
}

func setMvEnv(envs *[]cpEnv, path string, file fs.DirEntry) error {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	stat, ok := fileinfo.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New("Not a syscall.Stat_t")
	}

	tomlfile, err := os.Open(path)
	if err != nil {
		die("Open Env file error %v:", err)
	}
	defer tomlfile.Close()

	tml := &envToml{}
	if err := tml.loadConf(tomlfile); err != nil {
		return err
	}

	env := &cpEnv{path: path, tml: *tml, stat: stat}
	*envs = append(*envs, *env)

	return nil
}
