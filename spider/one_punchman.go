package spider

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/keima/util"
	"path"
	"strconv"
	"strings"
	"sync"
)

var (
	rootUrl      = "https://manhua.fzdm.com/132/"
	saveDir      = "E:\\keke\\one_punchMan\\"
	chapterNumWg = sync.WaitGroup{}
	page         = 50
)

/**
一拳超人的爬虫程序
@param chapter int 章节信息，如果传0，则默认爬取所有章节
*/
func OnePunchRun(chapter interface{}) {
	//如果是int，则默认爬取固定章节
	if value, ok := chapter.(int); ok {
		if value != 0 {
			var valString string
			if value < 10 {
				valString = "00" + strconv.Itoa(value)
			} else if value >= 10 && value <= 20 {
				valString = "0" + strconv.Itoa(value)
			} else {
				valString = strconv.Itoa(value)
			}
			chapterWg := sync.WaitGroup{}
			chapterWg.Add(page)
			for i := 0; i <= page; i++ {
				go onePunchChapterGet(&chapterWg, valString, i)
			}
			chapterWg.Wait()
		} else {
			//下载所有的章节
			//获取页面的章节信息(一共有多少章节)
			chapterNum := getChapterNum()
			chapterNumWg.Add(chapterNum)
			//循环章节
			for x := 1; x <= chapterNum; x++ {
				var valString string
				if x < 10 {
					valString = "00" + strconv.Itoa(x)
				} else if x >= 10 && x <= 20 {
					valString = "0" + strconv.Itoa(x)
				} else {
					valString = strconv.Itoa(x)
				}
				//开启goroutine读取该章节的页
				go func(x string) {
					chapterWg := sync.WaitGroup{}
					chapterWg.Add(page)
					for i := 0; i < page; i++ {
						go onePunchChapterGet(&chapterWg, x, i)
					}
					chapterWg.Wait()
					chapterNumWg.Done()
				}(valString)
			}
			chapterNumWg.Wait()
		}
	}
	//如果是slice，则爬取slice传递的章节
	if value, ok := chapter.([]int); ok {
		start := value[0]
		end := value[1]
		chapterNumWg.Add(end - start + 1)
		//循环章节
		for x := start; x <= end; x++ {
			var valString string
			if x < 10 {
				valString = "00" + strconv.Itoa(x)
			} else if x >= 10 && x <= 20 {
				valString = "0" + strconv.Itoa(x)
			} else {
				valString = strconv.Itoa(x)
			}
			//开启goroutine读取该章节的页
			go func(x int) {
				chapterWg := sync.WaitGroup{}
				chapterWg.Add(page)
				for i := 0; i < page; i++ {
					go onePunchChapterGet(&chapterWg, valString, i)
				}
				chapterWg.Wait()
				chapterNumWg.Done()
			}(x)
		}
		chapterNumWg.Wait()
	}

}

func getChapterNum() (num int) {
	num = 0
	c := onePunchCollector()
	//请求之前的设置
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})
	//爬取到html后
	c.OnHTML("meta", func(e *colly.HTMLElement) {
		if property := e.Attr("property"); property == "og:novel:latest_chapter_name" {
			chapterNum, err := strconv.Atoi(strings.TrimRight(strings.TrimLeft(e.Attr("content"), "一拳超人"), "话"))
			if err != nil {
				fmt.Println("截取转换章节数报错:", err)
			} else {
				num = chapterNum
			}
		}
	})
	//处理完html后
	c.OnScraped(func(r *colly.Response) {

	})
	err := c.Visit(rootUrl)
	if err != nil {
		fmt.Println(err)
	}
	return
}

func onePunchCollector() *colly.Collector {
	//获取一个收集器
	c := colly.NewCollector()
	//设置代理和请求头
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36"
	//请求之前的设置
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("Accept-Encoding", "gzip, deflate, br")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Host", "keep-alive")
		r.Headers.Set("If-None-Match", "3bed-58bbc967940f8")
		r.Headers.Set("Cookie", "Hm_lvt_cb51090e9c10cda176f81a7fa92c3dfc=1561009934,1561527502; picHost=www-mipengine-org.mipcdn.com/i/p1.manhuapan.com; Hm_lpvt_cb51090e9c10cda176f81a7fa92c3dfc=1561535595")
	})
	return c
}

/**
进入到章节页面中，爬取固定章节
*/
func onePunchChapterGet(chapterWg *sync.WaitGroup, chapter string, i int) {
	/*defer func() {
		err := recover()
		if err != nil {
			fmt.Println("panic:", err)
		}
	}()*/

	fmt.Println("准备爬取章节：", chapter, "的第", i, "页")

	defer chapterWg.Done()

	url := fmt.Sprintf("%s%s/index_%v.html", rootUrl, chapter, i)
	//这里在对应的文件夹下新建对应标题的文件夹
	chapterSaveDir := fmt.Sprintf("%s%s/", saveDir, chapter)
	makeDirErr := util.MakeDir(chapterSaveDir)
	if makeDirErr != nil {
		fmt.Println(fmt.Sprintf("目录%s创建失败", chapterSaveDir))
		return
	}

	//获取一个新的收集器
	c := onePunchCollector()
	//请求之前的设置
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})
	//收到相应之后的处理
	c.OnResponse(func(resp *colly.Response) {
		/*response := string(resp.Body)
		fmt.Println(response)*/
	})
	//爬取到html后
	c.OnHTML("#header + script + br + script", func(e *colly.HTMLElement) {
		arr := strings.Split(e.Text, "\"")
		url := "http://www-mipengine-org.mipcdn.com/i/" + arr[17] + "/" + arr[5]
		fmt.Println("开始爬取图片：", url)
		fileName := fmt.Sprintf("%v_%s", i, path.Base(url))
		util.SaveFile(url, chapterSaveDir, fileName, 0)
	})
	//处理完html后
	c.OnScraped(func(r *colly.Response) {

	})
	//错误时的报错信息
	/*c.OnError(func(resp *colly.Response, errHttp error) {
		fmt.Println(errHttp)
		chapterWg.Done()
		return
	})*/
	err := c.Visit(url)
	if err != nil {
		fmt.Println(url, "报错了，获取不到信息:", err)
		return
	}
}
