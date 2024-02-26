/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-09-23 22:29:13
 * @LastEditTime: 2023-09-24 00:08:55
 * @LastEditors: yuanshisan
 */
package util

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func PathCreate(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func PathExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func FileExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsExist(err)
}

func FileCreate(content bytes.Buffer, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		return err
	}
	return file.Close()
}

func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		result = info.Size()
		return nil
	})
	return result
}

func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func GetFileType(src string) (string, error) {
	file, err := os.Open(src)

	if err != nil {
		return "", err
	}

	buff := make([]byte, 512)

	_, err = file.Read(buff)
	if err != nil {
		return "", err
	}

	filetype := http.DetectContentType(buff)

	return filetype, nil
}

func GetDirFiles(dir string) ([]string, error) {
	dirList, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var res []string

	for _, f := range dirList {
		if f.IsDir() {
			files, err := GetDirFiles(dir + string(os.PathSeparator) + f.Name())
			if err != nil {
				return nil, err
			}
			res = append(res, files...)
		} else {
			res = append(res, dir+string(os.PathSeparator)+f.Name())
		}
	}

	return res, nil
}
