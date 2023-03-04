package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type cp struct {
	src, dst                        string
	isMkdir, isOverwrite, isSymlink bool
}

func isContains(dsts []string) (string, bool) {
	encountered := map[string]bool{}
	for _, dst := range dsts {
		if !encountered[dst] {
			encountered[dst] = true
		} else {
			return dst, true
		}
	}
	return "", false
}

func chkUnique(cps []cp) error {
	dsts := []string{}
	for _, d := range cps {
		dsts = append(dsts, d.dst)
	}

	if dup, err := isContains(dsts); err {
		return errors.New("Same dest Path:" + dup + "\n")
	}

	return nil
}

func f2d(envs *[]cpEnv, hidden bool) error {
	cps := []cp{}
	for _, env := range *envs {
		envdir := filepath.Dir(env.path)
		paths, err := fileWalk(env.path, envdir, hidden)
		if err != nil {
			return err
		}

		cpPaths, err := reNewPaths(envdir, env.tml.Dest_dir, paths)
		if err != nil {
			return err
		}

		for k, v := range cpPaths {
			cps = append(cps, cp{
				src:         k,
				dst:         v,
				isMkdir:     env.tml.IsMkdir,
				isOverwrite: env.tml.IsOverwrite,
				isSymlink:   env.tml.IsSymlink,
			})
		}
	}

	if err := chkUnique(cps); err != nil {
		return err
	}

	if err := file2dir(cps); err != nil {
		return err
	}

	return nil
}

func file2dir(cps []cp) error {
	for _, cp := range cps {
		if err := createDir(cp); err != nil {
			return err
		}
	}

	eg := new(errgroup.Group)
	for _, cp := range cps {
		cp := cp
		eg.Go(func() error {
			if err := copyFile(cp, cp.isOverwrite); err != nil {
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

func copyFile(cp cp, overwrite bool) error {
	_, err := os.Stat(cp.dst)
	if os.IsNotExist(err) {
		if err := cpExec(cp.src, cp.dst, cp.isSymlink); err != nil {
			return err
		}
	} else {
		if overwrite {
			if err := cpExec(cp.src, cp.dst, cp.isSymlink); err != nil {
				return err
			}
		} else {
			fmt.Fprintf(
				os.Stderr,
				"[warinig]: Overwrite skip from %v to %v\n",
				cp.src, cp.dst,
			)
		}
	}

	return nil
}

type symOrg struct {
	orgPath string
	err     error
}

func getRealPath(sympath string) symOrg {
	realPath := sympath
	info, err := os.Lstat(sympath)
	if err != nil {
		return symOrg{realPath, nil}
	}

	if info.Mode()&os.ModeSymlink == os.ModeSymlink {
		realp, err := os.Readlink(sympath)
		if err != nil {
			return symOrg{realPath, err}
		}
		realPath = filepath.Dir(sympath) + "/" + realp
	}

	return symOrg{realPath, nil}
}

func cpExec(src, dst string, isSymlink bool) error {
	sRealPath, dRealPath := src, dst

	symOrg := getRealPath(src)
	if symOrg.err != nil {
		return symOrg.err
	}
	sRealPath = symOrg.orgPath

	symOrg = getRealPath(dst)
	if symOrg.err != nil {
		return symOrg.err
	}
	dRealPath = symOrg.orgPath

	if isSymlink {
		_, err := os.Stat(dRealPath)
		if !os.IsNotExist(err) {
			if err := os.Remove(dRealPath); err != nil {
				return err
			}
		}
		if err := os.Symlink(sRealPath, dRealPath); err != nil {
			return err
		}
	} else {
		fsrc, err := os.Open(sRealPath)
		if err != nil {
			return err
		}
		defer fsrc.Close()

		fdst, err := os.Create(dRealPath)
		if err != nil {
			return err
		}
		defer fdst.Close()

		_, err = io.Copy(fdst, fsrc)
		if err != nil {
			return err
		}
	}

	return nil
}

func createDir(cp cp) error {
	newdir := filepath.Dir(cp.dst)
	_, err := os.Stat(newdir)
	if os.IsNotExist(err) {
		if cp.isMkdir {
			if err := os.MkdirAll(newdir, 0755); err != nil {
				return err
			}
		} else {
			return errors.New("Isn't exist dir: " + newdir)
		}
	}

	return nil
}

func reNewPaths(src, dest string, paths []string) (map[string]string, error) {
	newPaths := map[string]string{}
	for _, path := range paths {
		newPath := filepath.Clean(strings.Replace(path, src, dest, 1))
		newPaths[path] = newPath
	}

	return newPaths, nil
}

func fileWalk(root, dir string, hidden bool) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		if hidden && strings.HasPrefix(filepath.Base(file.Name()), ".") {
			continue
		}

		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			searchPaths, err := fileWalk(root, path, hidden)
			if err != nil {
				return nil, err
			}
			paths = append(paths, searchPaths...)
		} else {
			if file.Name() == envfile {
				if path != root {
					return nil, nil
				}
				continue
			}
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}

	return paths, nil
}
