package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/keima/util"
	"strconv"
	"sync"
)

var (
	nijieroChUrl     = "https://nijiero-ch.com/page/%s"
	nijieroChRootDir = "E:\\keke\\nijiero-ch\\"
	nijieroChWg      = sync.WaitGroup{}
	nijieroChWg2     = sync.WaitGroup{}
)

//爬取nijiero-ch.com的页面
func NijieroChRun(page int) {
	//如果输入的页码为0，则爬取所有的图片（可能会很慢哦）
	if page == 0 {
		//todo 这里通过爬取页面，获取总共的页码数

	} else {
		nijieroChWg.Add(page)
		//爬取第一页到输入的页码
		for i := 1; i <= page; i++ {
			go nijieroChSpiderRun(page)
		}
	}
	nijieroChWg.Wait()
}

//获取收集器
func nijieroChGetCollector() *colly.Collector {
	//获取一个收集器
	c := colly.NewCollector()
	//设置代理和请求头
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	//请求之前的设置
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
		r.Headers.Set("Accept-Language", "en-hk,zh;q=0.9,en;q=0.8")
		r.Headers.Set("cache-control", "max-age=0")
		r.Headers.Set("Cookie", "__cfduid=d5ff0516f6bee2e798cf2a30dbce201fb1562979989; 6666cd76f96956469e7be39d750cc7d9=1562979991; swpm_session=f8c73f9afa557dd3741f1edf309096f4; 516cb421a0b9e4c5876b936b3c266642=1562980006; 4635521ffb8828249a72cc8b1deda0d4=1562980018; b45d97180970a64d22c6cd45f5657c39=1562980137; 1dd124a9ff39b5799e121fc0c3e01577=1562980149")
	})
	return c
}

func nijieroChSpiderRun(page int) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		nijieroChWg.Done()
	}()

	//根据传递的page得到爬取的url
	url := fmt.Sprintf(nijieroChUrl, strconv.Itoa(page))
	//得到一个收集器
	c := nijieroChGetCollector()
	//收到相应之后的处理
	c.OnResponse(func(resp *colly.Response) {
		//response := string(resp.Body)
		response := resp.Body
		fmt.Println(response)
	})

	//爬取到html后
	c.OnHTML(".box > a, .box > h2 > a", func(e *colly.HTMLElement) {
		/*link := e.Attr("href")
		title := e.Attr("title")
		fmt.Println("获取外部连接：", link)
		//如果取到了html的字符串，则往里进一层
		if -1 != strings.Index(link, "html") {
			go nijieroChSpiderImageRun(link, title)
		}*/
	})

	//错误时的报错信息
	c.OnError(func(resp *colly.Response, errHttp error) {
		fmt.Println(errHttp)
	})

	c.OnScraped(func(r *colly.Response) {
		//moeimgWg.Done()
	})

	err := c.Visit(url)

	if err != nil {
		fmt.Println("报错啦！", err)
	}
}

//第二次抓取页面
func nijieroChSpiderImageRun(link, title string) {
	nijieroChWg2.Add(1)
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		nijieroChWg2.Done()
	}()

	//最终保存路径
	saveDir := fmt.Sprintf("%s%s\\", nijieroChRootDir, title)
	//这里在对应的文件夹下新建对应标题的文件夹
	makeDirErr := util.MakeDir(saveDir)
	if makeDirErr != nil {
		fmt.Println(fmt.Sprintf("目录%s创建失败，爬取%s页面失败", saveDir, link))
		return
	}

	fmt.Println("我准备再次访问链接了！", link)
	c2 := nijieroChGetCollector()

	c2.OnHTML(".box > a", func(e2 *colly.HTMLElement) {
		imageSrc := e2.Attr("href")
		fmt.Println("获取图片连接：", imageSrc)
		//保存图片
		util.SaveFile(imageSrc, saveDir, "", 0)
	})

	c2.OnScraped(func(r *colly.Response) {
		//nijieroChWg2.Done()
	})

	err := c2.Visit(link)

	if err != nil {
		fmt.Println("2报错啦！", err)
	}
	return
}
