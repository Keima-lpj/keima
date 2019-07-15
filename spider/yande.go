package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/keima/util"
	"strconv"
	"sync"
)

var (
	yandeUrl        = "https://yande.re/post?page=%s"
	yandeWg         = sync.WaitGroup{}
	yandeWg2        = sync.WaitGroup{}
	yandeSaveFileWg = sync.WaitGroup{}
	yandeRootDir    = "E:\\keke\\yande\\"
)

func YandeRun(page int) {
	yandeWg.Add(page)
	//要爬取的页面
	for i := 1; i <= page; i++ {
		url := fmt.Sprintf(yandeUrl, strconv.Itoa(i))
		fmt.Println("开始采集页面：", url)
		go yandeSpiderRun(url)
	}
	yandeWg.Wait()
	yandeWg2.Wait()
	yandeSaveFileWg.Wait()
}

func yandeGetCollector() *colly.Collector {
	//获取一个收集器
	c := colly.NewCollector()
	//设置代理和请求头
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	//请求之前的设置
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("cache-control", "max-age=0")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Host", "keep-alive")
		r.Headers.Set("If-None-Match", "W/\"6dbf0c32b3169d7d35c74d8c76e9824d\"")
		r.Headers.Set("Cookie", "country=CN; vote=1; __utmc=5621947; __utmz=5621947.1561082582.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none); __utma=5621947.37221620.1561082582.1561085481.1561088170.3; __utmt=1; forum_post_last_read_at=%222019-06-21T05%3A36%3A38.083%2B02%3A00%22; yande.re=MzVlYiswblI3SGY1dWVTZ3JjR3ZObExhRHBseXFYOU5Pbm1MTmtzSnFtak9yZE16MDFwSnRPa3dTTHV0TXFzRkx1RklRVGFvZnVKSHNLUEFJR2U3QkluLzc1ZDlZWTgyUlFrQXJibmsvc3dKZGRSSzNVOHV0d3RMYUZlT1VBdDRoSkZlZ3A2QnBFL1JSUEY2c1Q4WHJnPT0tLU44SXZVK3M4VzlFeDZzMXdrQ3AxZnc9PQ%3D%3D--2fa3210125c81d80e8a1c940ebd053a16b56db34; __utmb=5621947.4.10.1561088170")
	})
	return c
}

/**
第一次抓取页面内容
*/
func yandeSpiderRun(url string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			yandeWg.Done()
		}
	}()
	//获取一个收集器
	c := yandeGetCollector()
	//收到相应之后的处理
	c.OnResponse(func(resp *colly.Response) {
		/*response := string(resp.Body)
		fmt.Println(response)*/
	})
	//爬取到html后
	c.OnHTML("#post-list-posts > li > a", func(e *colly.HTMLElement) {
		yandeWg2.Add(1)
		link := e.Attr("href")
		fmt.Println("获取图片链接：", link)
		go yandeSpiderImageRun(link)
	})
	//错误时的报错信息
	c.OnError(func(resp *colly.Response, errHttp error) {
		fmt.Println(errHttp)
	})
	//onhtml执行后，将计数器减1
	c.OnScraped(func(r *colly.Response) {
		yandeWg.Done()
	})

	err := c.Visit(url)

	if err != nil {
		fmt.Println("报错啦！", err)
	}
}

//爬取图片
func yandeSpiderImageRun(link string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//这里在对应的文件夹下新建对应标题的文件夹
	makeDirErr := util.MakeDir(yandeRootDir)
	if makeDirErr != nil {
		fmt.Println(fmt.Sprintf("目录%s创建失败，爬取%s页面失败", yandeRootDir, link))
		return
	}
	yandeSaveFileWg.Add(1)
	go util.SaveFile(link, yandeRootDir, "", &yandeSaveFileWg)
	yandeWg2.Done()
	return
}
