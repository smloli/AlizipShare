package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
)

type File struct {
	Path string
	PathList []string
}

func (file *File) Modify() {
	readByte := make([]byte, 4)
	writeByte := make([]byte, 4)
	for _, v := range file.PathList {
		// 打开文件
		f, err := os.OpenFile(file.Path + v, os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("error:", err)
		}
		// 读取文件头4个字节的内容，将内容赋值给readByte
		_, err = f.ReadAt(readByte, 0)
		if err != nil {
			fmt.Println("error:", err)
		}
		// // readByte 异或 0xF，将结果赋值给writeByte
		for i, v := range readByte {
			writeByte[i] = v ^ 0xF
		}
		// 向文件头写入writeByte，长度为4个字节
		_, err = f.WriteAt(writeByte, 0)
		if err != nil {
			fmt.Println("error:", err)
		}
		f.Close()
		oldName := file.Path + v
		fmt.Println(oldName)
		// 判断文件格式
		if strings.HasSuffix(oldName, ".png") {
			os.Rename(oldName, oldName[:len(oldName) - 4]) // 恢复原文件名
		} else {
			os.Rename(oldName, oldName + ".png")	// 原文件名 + .png
		}
	}
}

// 返回指定目录里的文件列表
func (file *File) getFileInfo() {
	f, err := ioutil.ReadDir(file.Path)
	if err != nil {
		fmt.Println("error:", err)
	}

	for _, v := range f {
		// 判断是否为文件夹
		if !v.IsDir() {
			if file.PathList[0] == "" {
				file.PathList[0] = v.Name()
				continue
			}
			file.PathList = append(file.PathList, v.Name())
		}
	}
}

// 处理带空格的路径
func (file *File) getPath() {
	var tmp string
	fmt.Print("Path:")
	for {
		n, _ := fmt.Scanf("%s", &tmp)
		if n == 0 {
			break
		}
		file.Path += tmp + " "
	}
	file.Path = strings.TrimSpace(file.Path)
}

func main() {
	var file File
	file.getPath()
	file.PathList = make([]string, 1)
	// 判断是否为文件夹
	if b, _ := os.Stat(file.Path); b.IsDir() {
		file.getFileInfo()
	} else {
		index := strings.LastIndex(file.Path, "\\")
		// 获取文件名
		file.PathList[0] = file.Path[index + 1:]
		// file.Path 去除文件名
		file.Path = file.Path[:index]
	}
	file.Path += "\\"
	file.Modify()
}