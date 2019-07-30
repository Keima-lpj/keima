package spider

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/keima/util"
	"strconv"
	"strings"
	"sync"
)

var (
	catRootDir       = "E:\\keke\\cat\\"
	catBaseTitle     = ""
	baseUrl          = "https://www.969uy.com/"
	catPageWaitGroup = sync.WaitGroup{}
	catWaitGroup     = sync.WaitGroup{}
	catSaveFileGroup = sync.WaitGroup{}
)

//程序开始
func CatRun(getType, startPage, endPage int) {
	for i := startPage; i <= endPage; i++ {
		//根据要爬取的type拼接url
		url, err := getUrl(getType, strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
			continue
		}
		catPageWaitGroup.Add(1)
		go catSpiderRun(url, i)
	}
	catPageWaitGroup.Wait()
	catWaitGroup.Wait()
}

//根据选择的数字进入对应的页面
func getUrl(getType int, page string) (string, error) {
	var url string

	if page == "1" {
		page = ""
	} else {
		page = "-" + page
	}
	switch getType {
	case 0:
		url = fmt.Sprintf("%s%s%s.%s", baseUrl, "tupian/list-%E8%87%AA%E6%8B%8D%E5%81%B7%E6%8B%8D", page, "html")
		catBaseTitle = "自拍偷拍"
	case 1:
		url = fmt.Sprintf("%s%s%s.%s", baseUrl, "tupian/list-%E4%BA%9A%E6%B4%B2%E8%89%B2%E5%9B%BE", page, "html")
		catBaseTitle = "亚洲色图"
	case 2:
		url = fmt.Sprintf("%s%s%s.%s", baseUrl, "tupian/list-%E6%AC%A7%E7%BE%8E%E8%89%B2%E5%9B%BE", page, "html")
		catBaseTitle = "欧美色图"
	case 3:
		url = fmt.Sprintf("%s%s%s.%s", baseUrl, "tupian/list-%E7%BE%8E%E8%85%BF%E4%B8%9D%E8%A2%9C", page, "html")
		catBaseTitle = "美腿丝袜"
	case 4:
		url = fmt.Sprintf("%s%s%s.%s", baseUrl, "tupian/list-%E6%B8%85%E7%BA%AF%E5%94%AF%E7%BE%8E", page, "html")
		catBaseTitle = "清纯唯美"
	case 5:
		url = fmt.Sprintf("%s%s%s.%s", baseUrl, "tupian/list-%E4%B9%B1%E4%BC%A6%E7%86%9F%E5%A5%B3", page, "html")
		catBaseTitle = "乱伦熟女"
	case 6:
		url = fmt.Sprintf("%s%s%s.%s", baseUrl, "tupian/list-%E5%8D%A1%E9%80%9A%E5%8A%A8%E6%BC%AB", page, "html")
		catBaseTitle = "卡通动漫"
	case 7:
		url = fmt.Sprintf("%s%s", baseUrl, "meinv/index.html")
		catBaseTitle = "极品美女"
	default:
		return "", errors.New("输入的数字有误！")
	}
	return url, nil
}

//获取一个收集器
func catGetCollector() *colly.Collector {
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
		r.Headers.Set("cache-control", "no-cache")
		r.Headers.Set("pragma", "no-cache")
		//r.Headers.Set("referer", "https://www.969uy.com/index/home.html")
		r.Headers.Set("cookie", "__cfduid=da6ec18fadae83b3f64db5067bd0620351564487560; _ga=GA1.2.940297954.1564487564; _gid=GA1.2.1762550513.1564487564; _gat_gtag_UA_126205200_1=1; _gat_gtag_UA_126205200_2=1; Hm_lvt_427f72ce75b0677eb10f24419484eb80=1564487566; Hm_lpvt_427f72ce75b0677eb10f24419484eb80=1564487594; _gat_gtag_UA_138595290_2=1; playss=1")
	})
	return c
}

//主逻辑
func catSpiderRun(url string, page int) {
	//获取一个收集器
	c := catGetCollector()
	//收到相应之后的处理
	c.OnResponse(func(resp *colly.Response) {
		/*response := string(resp.Body)
		fmt.Println(response)*/
	})

	//爬取到html后
	c.OnHTML("#tpl-img-content > li > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")
		link = baseUrl + link
		fmt.Println("获取外部连接：", link)

		//判断，如果link中农有meinv，则继续下一层
		if strings.Index(link, "meinv") != -1 {
			catMeinvSpiderImageRun(link)
		} else {
			catWaitGroup.Add(1)
			//第二次爬取页面
			go catSpiderImageRun(link, title)
		}
	})

	//错误时的报错信息
	c.OnError(func(resp *colly.Response, errHttp error) {
		fmt.Println(errHttp)
	})

	c.OnScraped(func(r *colly.Response) {
		catPageWaitGroup.Done()
	})

	err := c.Visit(url)

	if err != nil {
		fmt.Println("报错啦！", err)
		catPageWaitGroup.Done()
	}
}

//保存图片
func catSpiderImageRun(link, title string) {
	//最终保存路径
	saveDir := fmt.Sprintf("%s\\%s\\%s\\", catRootDir, catBaseTitle, title)
	//这里在对应的文件夹下新建对应标题的文件夹
	makeDirErr := util.MakeDir(saveDir)
	if makeDirErr != nil {
		fmt.Println(fmt.Sprintf("目录%s创建失败，爬取%s页面失败", saveDir, link))
		return
	}

	fmt.Println("我准备再次访问链接了！", link)
	c2 := catGetCollector()

	c2.OnHTML(".videopic", func(e2 *colly.HTMLElement) {
		imageSrc := e2.Attr("data-original")
		fmt.Println("获取图片连接：", imageSrc)
		//保存图片
		catSaveFileGroup.Add(1)
		util.SaveFile(imageSrc, saveDir, "", &catSaveFileGroup)
	})

	c2.OnScraped(func(r *colly.Response) {
		catWaitGroup.Done()
	})

	err := c2.Visit(link)

	if err != nil {
		fmt.Println("2报错啦！", err)
		catWaitGroup.Done()
	}
	return
}

//访问美女页面
func catMeinvSpiderImageRun(url string) {
	//获取一个收集器
	c := catGetCollector()
	//收到相应之后的处理
	c.OnResponse(func(resp *colly.Response) {
		/*response := string(resp.Body)
		fmt.Println(response)*/
	})

	//爬取到html后
	c.OnHTML("#tpl-img-content > li > a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Attr("title")
		link = baseUrl + link
		fmt.Println("再次获取外部连接：", link)
		catWaitGroup.Add(1)
		//第二次爬取页面
		go catSpiderImageRun(link, title)
	})

	//错误时的报错信息
	c.OnError(func(resp *colly.Response, errHttp error) {
		fmt.Println(errHttp)
	})

	c.OnScraped(func(r *colly.Response) {
	})

	err := c.Visit(url)

	if err != nil {
		fmt.Println("报错啦！", err)
	}
}
