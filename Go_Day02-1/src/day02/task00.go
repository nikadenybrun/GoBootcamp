package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	flags     string
	extension string
)

func walkFunc(path string, info os.FileInfo, err error) error {
	_, err1 := os.Stat(path)
	if err1 != nil {
		if os.IsPermission(err) {
			return nil
		}
	}
	if err != nil {
		fmt.Println("err")
		return err
	}
	if !strings.Contains(flags, "-d") && info.IsDir() {
		return nil
	}
	if !strings.Contains(flags, "-sl") {
		info1, _ := os.Lstat(path)
		if info1.Mode()&os.ModeSymlink == os.ModeSymlink {
			return nil
		}
	}
	if !strings.Contains(flags, "-f") {
		info1, err2 := os.Lstat(path)
		if err2 != nil {
			fmt.Println("err2")
			return err2
		}
		if info1.Mode()&os.ModeSymlink != os.ModeSymlink && !info.IsDir() {
			return nil
		}
	}
	if !strings.Contains(flags, "-f") && strings.Contains(flags, "-ext") {

		return nil
	}
	if strings.Contains(flags, "-f") && strings.Contains(flags, "-ext") {
		info1, _ := os.Lstat(path)
		if !strings.HasSuffix(path, extension) && info1.Mode()&os.ModeSymlink != os.ModeSymlink && !info.IsDir() {
			return nil
		}
	}
	info1, _ := os.Lstat(path)
	if strings.Contains(flags, "-sl") && info1.Mode()&os.ModeSymlink == os.ModeSymlink {
		targetPath, err3 := os.Readlink(path)
		if err3 != nil {
			fmt.Println("err3")
			return err3
		}
		targetFullPath := filepath.Join(filepath.Dir(path), targetPath)
		_, err4 := os.Stat(targetFullPath)
		if err4 != nil {
			fmt.Printf("%s -> [broken]\n", path)
		} else {
			fmt.Printf("%s -> %s\n", path, path[:len(path)-len(info.Name())]+targetPath)

		}

	} else {
		fmt.Printf("%s\n", path)
	}
	return nil
}

func main() {
	var directory string

	if len(os.Args) > 1 {
		directory = os.Args[len(os.Args)-1]
		for i := 1; i < len(os.Args)-1; i++ {
			if strings.HasSuffix(flags, "-ext") {
				extension = os.Args[i]
			}
			flags += os.Args[i]
		}
	} else {
		fmt.Println("No path provided")
	}
	if flags == "" {
		flags += "-f-d-sl"
	}
	if err := filepath.Walk(directory, walkFunc); err != nil {
		fmt.Println(err)
	}
}
