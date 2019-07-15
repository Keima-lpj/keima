package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/keima/util"
	"strconv"
	"strings"
	"sync"
)

var (
	moeimgUrl        = "http://moeimg.net/page/%s"
	moeimgWg         = sync.WaitGroup{}
	moeimgWg2        = sync.WaitGroup{}
	moeimgSaveFileWg = sync.WaitGroup{}
	moeimgRootDir    = "E:\\keke\\moeimg\\"
)

func MoeimgRun(page int) {
	moeimgWg.Add(page)
	//要爬取的页面
	for i := 1; i <= page; i++ {
		url := fmt.Sprintf(moeimgUrl, strconv.Itoa(i))
		fmt.Println("开始采集页面：", url)
		go moeimgSpiderRun(url)
	}
	moeimgWg.Wait()
	moeimgWg2.Wait()
	moeimgSaveFileWg.Wait()
}

func moeimgGetCollector() *colly.Collector {
	//获取一个收集器
	c := colly.NewCollector()
	//设置代理和请求头
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	//请求之前的设置
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Connection", "moeimg.net")
		r.Headers.Set("Host", "keep-alive")
		r.Headers.Set("If-None-Match", "3bed-58bbc967940f8")
		r.Headers.Set("Cookie", "__cfduid=df97a01336e673c1cdc31c97929055bdf1561017278; _ga=GA1.2.1164906477.1561017281; _gid=GA1.2.11570314.1561017281")
	})
	return c
}

/**
第一次抓取页面内容
*/
func moeimgSpiderRun(url string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			moeimgWg.Done()
		}
	}()

	//获取一个收集器
	c := moeimgGetCollector()

	//收到相应之后的处理
	c.OnResponse(func(resp *colly.Response) {
		/*response := string(resp.Body)
		fmt.Println(response)*/
	})

	//爬取到html后
	c.OnHTML(".box > a, .box > h2 > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")
		fmt.Println("获取外部连接：", link)
		//如果取到了html的字符串，则往里进一层
		if -1 != strings.Index(link, "html") {
			go moeimgSpiderImageRun(link, title)
		}
	})

	//错误时的报错信息
	c.OnError(func(resp *colly.Response, errHttp error) {
		fmt.Println(errHttp)
	})

	c.OnScraped(func(r *colly.Response) {
		moeimgWg.Done()
	})

	err := c.Visit(url)

	if err != nil {
		fmt.Println("报错啦！", err)
	}
}

/**
第二次抓取页面内容
*/
func moeimgSpiderImageRun(link, title string) {
	moeimgWg2.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			moeimgWg2.Done()
		}
	}()

	//最终保存路径
	saveDir := fmt.Sprintf("%s%s\\", moeimgRootDir, title)

	//这里在对应的文件夹下新建对应标题的文件夹
	makeDirErr := util.MakeDir(saveDir)
	if makeDirErr != nil {
		fmt.Println(fmt.Sprintf("目录%s创建失败，爬取%s页面失败", saveDir, link))
		return
	}

	fmt.Println("我准备再次访问链接了！", link)
	c2 := moeimgGetCollector()

	c2.OnHTML(".box > a", func(e2 *colly.HTMLElement) {
		imageSrc := e2.Attr("href")
		fmt.Println("获取图片连接：", imageSrc)
		//保存图片
		moeimgSaveFileWg.Add(1)
		go util.SaveFile(imageSrc, saveDir, "", &moeimgSaveFileWg)
	})

	c2.OnScraped(func(r *colly.Response) {
		moeimgWg2.Done()
	})

	err := c2.Visit(link)

	if err != nil {
		fmt.Println("2报错啦！", err)
		moeimgWg2.Done()
	}
	return
}
