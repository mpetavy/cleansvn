package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	path    = flag.String("p", "", "working dir")
	test    = flag.Bool("t", false, "only check")
	exclude = flag.String("x", ".idea", "exclude files")
)

func init() {
	common.Init(false, "1.0.0", "", "2018", "Tool to clean a subversion controlled dir", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, run, 0)
}

func run() error {
	cmd := exec.Command("svn", "info")
	if *path == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil
		}

		path = &dir
	}

	cmd.Dir = *path

	b, err := cmd.Output()
	if err != nil {
		return err
	}

	output, err := common.ToUTF8String(string(b[:]), common.DefaultEncoding())
	if err != nil {
		return err
	}

	if strings.Index(strings.ToLower(output), "revision:") == -1 {
		return fmt.Errorf("the path %s is not a working copy", *path)
	}

	cmd = exec.Command("svn", "status")

	cmd.Dir = *path

	b, err = cmd.Output()
	if err != nil {
		return err
	}

	output, err = common.ToUTF8String(string(b[:]), common.DefaultEncoding())
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "?") {
			file := strings.TrimSpace(line[1:])

			for _, item := range strings.Split(*exclude, ";") {
				if item == file {
					common.Info("skip %s", file)
					continue
				}
			}

			fp := filepath.Join(*path, file)

			if common.FileExists(fp) {
				if *test {
					common.Info("delete %s [simulate]", fp)
				} else {
					common.Info("delete %s", fp)
					err := os.RemoveAll(fp)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
