package main

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type Info struct {
	paths     []paths
	filenames []string
}
type paths struct {
	filename    string
	dirname     string
	new_dirName string
}

func (f *Info) parseFiles() *Info {

	for i, file := range f.filenames {
		count := strings.Count(file, "/")
		for j := 0; j < count; j++ {
			_, f.paths[i].filename, _ = strings.Cut(f.paths[i].filename, "/")

		}
		f.paths[i].dirname, _ = strings.CutSuffix(file, f.paths[i].filename)

		if f.paths[i].new_dirName == file {
			f.paths[i].new_dirName = f.paths[i].dirname

		}
		if string(f.paths[i].dirname[0]) == "/" {
			f.paths[i].dirname = f.paths[i].dirname[1:]
		}
		if string(f.paths[i].new_dirName[0]) == "/" {
			f.paths[i].new_dirName = f.paths[i].new_dirName[1:]
		}
	}
	mkdir(f.paths[0].new_dirName)
	return f
}
func (f *Info) getArgs() *Info {
	if len(os.Args) > 1 {
		var i int = 1
		if os.Args[1] == "-a" {
			// f.paths[0].new_dirName = os.Args[2]
			i = 3
		}
		for ; i < len(os.Args); i++ {
			f.filenames = append(f.filenames, os.Args[i])
			f.paths = append(f.paths, paths{
				filename:    os.Args[i],
				dirname:     os.Args[i],
				new_dirName: "", //f.paths[0].new_dirName
			})
		}

		for i := 0; i < len(f.filenames); i++ {
			if os.Args[1] == "-a" {
				f.paths[i].new_dirName = os.Args[2]
			} else {
				f.paths[i].new_dirName = os.Args[i+1]
			}
		}

	} else {
		log.Fatal("No path provided")
	}
	return f
}

func archivate(p paths, wg *sync.WaitGroup) {
	fmt.Println(p.filename)
	defer wg.Done()

	tarFile, err := os.Create(getFileName((p.new_dirName+"/")+p.filename) + ".tar.gz")
	if err != nil {
		log.Fatal(err)
	}
	tarWriter := tar.NewWriter(tarFile)
	defer tarFile.Close()
	fileInfo, err := os.Stat(p.dirname + p.filename)
	if err != nil {
		log.Fatal(err)
	}
	hdr := &tar.Header{
		Name: p.filename,
		Mode: 0600,
		Size: int64(fileInfo.Size()),
	}
	if err := tarWriter.WriteHeader(hdr); err != nil {
		log.Fatal(err)
	}
	src, err := os.Open(p.dirname + p.filename)
	if err != nil {
		log.Fatal(err)
	}
	defer src.Close()
	if _, err = io.Copy(tarWriter, src); err != nil {
		log.Fatal(err)
	}
	if err := tarWriter.Close(); err != nil {
		log.Fatal(err)
	}
}

func getFileName(filename string) string {
	if strings.HasSuffix(filename, ".log") {
		idx := strings.Index(filename, ".log")
		filename = filename[:idx]
		filename += "_" + fmt.Sprintf("%d", time.Now().Unix())

		return filename
	}
	return ""
}

func mkdir(dirname string) {
	_, err := os.Stat(dirname)

	if os.IsNotExist(err) {
		err := os.Mkdir(dirname, 0755)
		if err != nil {
			log.Fatal("didnt manage to create a dir")
		}

	}
}
func main() {
	info := Info{}
	info.getArgs()
	info.parseFiles()

	var wg sync.WaitGroup
	wg.Add(len(info.filenames))
	for i, filename := range info.paths {
		info.paths[i].new_dirName = info.paths[0].new_dirName
		go archivate(filename, &wg)
	}
	wg.Wait()

}
