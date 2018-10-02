package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func main() {
	f, err := ioutil.TempFile("/tmp", "aaa")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = f.WriteString("secret"); err != nil {
		log.Fatal(err)
	}
	if _, err = f.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}
	// do not close

	// unlink tempfile
	if err = os.Remove(f.Name()); err != nil {
		log.Fatal(err)
	}

	if true {
		log.Println("@@@@@ doing actual exec @@@@@")
		// exec

		if coe, err := isCloseOnExec(f.Fd()); err != nil {
			log.Fatal(err)
		} else if coe {
			fmt.Println("existing temp file has CLOEXEC set")

			fmt.Println("unsetting CLOEXEC on temp file")

			if err := unsetCloseOnExec(f.Fd()); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("existing temp file does NOT have CLOEXEC set")
		}

		if err = unix.Exec("./child/child", nil, nil); err != nil {
			log.Fatal(err)
		}

	} else {
		log.Println("@@@@@ doing supervision @@@@@")

		cmd := exec.Command("./child/child")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.ExtraFiles = []*os.File{f}

		cmd.Run()
	}

}

func isCloseOnExec(fd uintptr) (bool, error) {
	flags, err := getFdFlags(fd)
	if err != nil {
		return false, err
	}
	return flags&unix.FD_CLOEXEC != 0, nil
}

func unsetCloseOnExec(fd uintptr) error {
	flags, err := getFdFlags(fd)
	if err != nil {
		return err
	}

	if _, err = unix.FcntlInt(fd, unix.F_SETFD, flags&(^unix.FD_CLOEXEC)); err != nil {
		return err
	}
	return nil
}

func getFdFlags(fd uintptr) (int, error) {
	return unix.FcntlInt(fd, unix.F_GETFD, 0)
}
