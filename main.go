package main

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"time"
)

type File struct {
	Path string
	PathList []string
}

func (file *File) Modify() {
	readByte := make([]byte, 4)
	writeByte := make([]byte, 4)
	loliByte := make([]byte, 4)
	var loli string
	for _, v := range file.PathList {
		// 打开文件
		oldName := file.Path + v
		f, err := os.OpenFile(oldName, os.O_RDWR, 0644)
		if err != nil {
			fmt.Println("error:", err)
		}
		// 读取文件头4个字节内容
		_, err = f.ReadAt(readByte, 0)
		if err != nil {
			fmt.Println("error:", err)
		}
		// 读取文件尾4个字节内容
		fileInfo, _ := f.Stat()
		_, err = f.ReadAt(loliByte, fileInfo.Size() - 4)
		if err != nil {
			fmt.Println("error:", err)
		}
		// loliByte转换成字符串
		for _, v := range loliByte {
			loli += string(v)
		}
		fmt.Println("loliByte = ", loliByte)
		// readByte异或0xF
		for i, v := range readByte {
			writeByte[i] = v ^ 0xF
		}
		// 文件头写入writeByte
		_, err = f.WriteAt(writeByte, 0)
		if err != nil {
			fmt.Println("error:", err)
		}

		// 取纳秒时间戳
		timestamp := fmt.Sprintf("%d%s", time.Now().UnixNano(), "loli")
		// 将时间戳转换成[]byte
		zeroByte := []byte{0,0,0,0,0,0,0,0,0}
		timestampByte := append(zeroByte, []byte(timestamp)...)

		// 当变量loli的值 == loli，直接重写时间戳，不再追加写入
		if loli != "loli" {
			// 写入timestampByte
			_, err = f.WriteAt(timestampByte, fileInfo.Size())
			if err != nil {
				fmt.Println("error:", err)
			}
		} else {
			// 重写timestampByte
			_, err = f.WriteAt(timestampByte, fileInfo.Size() - 32)
			if err != nil {
				fmt.Println("error:", err)
			}
		}
		f.Close()
		// 判断文件格式
		if strings.HasSuffix(oldName, ".png") {
			os.Rename(oldName, oldName[:len(oldName) - 4]) // 恢复原文件名
		} else {
			os.Rename(oldName, oldName + ".png")	// 原文件名 + .png
		}
		fmt.Println(oldName)
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
	args := os.Args
	// 如果软件启动参数>=2个，就以参数的形式获取要修改的文件或文件夹的路径
	if len(args) >= 2 {
		for i := 1; i < len(args); i++ {
			file.Path += args[i] + " "
		}
	} else {
		fmt.Print("Path:")
		for {
			n, _ := fmt.Scanf("%s", &tmp)
			if n == 0 {
				break
			}
			file.Path += tmp + " "
		}
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