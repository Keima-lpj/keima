package util

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

//从网页上下载图片，保存到对应的文件夹中
func SaveFile(fileSrc, fileDir, fileName string, filter int) {
	if fileName == "" {
		//获取文件基础名称
		fileName = path.Base(fileSrc)
	}
	if filter == 1 {
		//过滤特殊字符
		fileName = filterSpec(fileName)
	}
	//最终保存的文件名
	saveFileDir := fileDir + fileName

	//判断如果目标文件存在，则不继续操作
	if has, _ := PathExists(saveFileDir); !has {
		client := http.Client{}
		client.Timeout = time.Second * 300
		res, err := client.Get(fileSrc)

		//res, err := http.Get(fileSrc)
		if err != nil {
			fmt.Println("保存图片时报错了！", err)
			return
		}
		// defer后的为延时操作，通常用来释放相关变量
		defer res.Body.Close()
		// 获得get请求响应的reader对象
		reader := bufio.NewReaderSize(res.Body, 32*1024)

		file, err := os.Create(saveFileDir)
		if err != nil {
			return
		}
		// 获得文件的writer对象
		writer := bufio.NewWriter(file)
		n, copyErr := io.Copy(writer, reader)
		if copyErr != nil {
			fmt.Println("复制文件错误:", copyErr)
			return
		}
		fmt.Println("打印copy written", n)
	} else {
		fmt.Println("文件", saveFileDir, "已存在")
	}
}

//判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		//文件存在
		return true, nil
	}
	if os.IsNotExist(err) {
		//文件不存在
		return false, nil
	}
	//不确定是否存在
	return false, err
}

//新建目录
func MakeDir(dir string) (err error) {
	//可以使用Mkdir或者MkdirAll，前者父级目录必须存在，不然会失败，后者可以递归创建目录
	return os.MkdirAll(dir, 0766)
}

//过滤特殊字符
func filterSpec(fileName string) string {
	fileName = strings.Replace(fileName, "%", "", 10)
	fileName = strings.Replace(fileName, "yande.re", "", 10)
	return fileName
}
