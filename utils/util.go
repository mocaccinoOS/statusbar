package util

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

func Run(s string) (string, error) {
	cmd := exec.Command("/bin/sh", "-ce", s)
	fmt.Println("executing ", s)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("ret (fail)", string(out))

		return "", errors.Wrap(err, string(out))
		//log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Println("ret ", string(out))

	return string(out), nil
}

func Sudo(s string) (string, error) {
	return Run(fmt.Sprintf("pkexec /bin/sh -c \"%s\"", s))
}

func Touch(p string) error {
	err := os.MkdirAll(filepath.Dir(p), os.ModePerm)
	if err != nil {
		return err
	}
	emptyFile, err := os.Create(p)
	if err != nil {
		return err
	}
	emptyFile.Close()
	return nil
}
