package main

import (
	"flag"
	"fmt"
	"os"
)

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args)
	os.Exit(1)
}

func main() {
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		os.Exit(0)
	}

	file, err := os.Open(env)
	if err != nil {
		die("Open Env file error %v:", err)
	}
	defer file.Close()

	cfg := &envInfo{}
	if err := cfg.loadConf(file); err != nil {
		die("Struct Config error %v:", err)
	}
	cfg.chkFormat()

	envs := []cpEnv{}
	if err := getEnv(&envs, cfg.getRoot(), cfg.getIsExcludedHiddenFile()); err != nil {
		die("Get df2d.toml error %v:", err)
	}

	if err := f2d(&envs, cfg.getIsExcludedHiddenFile()); err != nil {
		die("f2d error %v:", err)
	}
}
