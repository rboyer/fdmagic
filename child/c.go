package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var vanitymap = map[string]string{
	"0": "stdin",
	"1": "stdout",
	"2": "stderr",
	"3": "extrafile-parent",
}

func main() {
	log.Println("===== list /dev/fd =====")
	list, err := ioutil.ReadDir("/dev/fd")
	if err != nil {
		log.Fatal(err)
	}
	var has3 bool
	for _, f := range list {
		full := filepath.Join("/dev/fd", f.Name())

		if f.Name() == "3" {
			has3 = true
		}

		vanity := vanitymap[f.Name()]

		switch {
		case f.Mode()&os.ModeSymlink != 0:
			path, err := os.Readlink(full)
			if os.IsNotExist(err) {
				path = "[dangling symlink]"
			} else if err != nil {
				log.Fatal(err)
			}
			log.Printf("link: %q -> %q [vanity:%q]", f.Name(), path, vanity)
		default:
			log.Printf("file: %q [vanity:%q]", f.Name(), vanity)
		}
	}

	if has3 {
		log.Println("===== read /dev/fd/3 by FD and NewFile =====")
		f := os.NewFile(3, "fd3")
		d, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		if _, err = f.Seek(0, io.SeekStart); err != nil {
			log.Fatal(err)
		}
		// do not close

		log.Println(string(d))

		log.Println("===== read /dev/fd/3 by reading it as a new file =====")
		f, err = os.Open("/dev/fd/3")

		d, err = ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(d))
	}
}
