package upload

import (
	"fmt"
	"os"
	"path/filepath"
)

type FileInfo struct {
	os.FileInfo
	Path string
	Dir  string
}

func uploadDir(opts *Options) error {
	files, err := AllFiles(opts.File)
	if err != nil {
		return err
	}
	for _, file := range files {
		err = uploadFile(opts, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func AllFiles(root string) (infos []FileInfo, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("dir: %v: name: %s size:%d dir:%s \n", info.IsDir(), path, info.Size(), filepath.Dir(path))

		fileInfo := FileInfo{FileInfo: info, Path: path, Dir: filepath.Dir(path)}
		infos = append(infos, fileInfo)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return
}
