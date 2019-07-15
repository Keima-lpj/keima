package util

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
)

// 获取当前页面的编码格式，自动判定网页编码格式并转码
func ExplainUrl(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("响应码错误:", resp.StatusCode)
	}
	// 转换编码格式
	bufReader := bufio.NewReader(resp.Body)
	encoding := determineEncoding(bufReader)
	//fmt.Println(encoding)
	reader := transform.NewReader(bufReader, encoding.NewDecoder())
	contents, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return string(contents)
}

// 获取当前页面的编码格式
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		// 如果没有获取到编码格式，则返回默认UTF-8编码格式
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
