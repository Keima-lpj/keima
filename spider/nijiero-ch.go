package spider

import (
	"fmt"
	"github.com/keima/util"
	"strconv"
	"strings"
	"sync"
)

var (
	nijieroChUrl        = "https://nijiero-ch.com/page/%s"
	nijieroChRootDir    = "E:\\keke\\nijiero-ch\\"
	nijieroChWg         = sync.WaitGroup{}
	nijieroChWg2        = sync.WaitGroup{}
	nijieroChSaveFileWg = sync.WaitGroup{}
)

//爬取nijiero-ch.com的页面
func NijieroChRun(startPage, endPage int) {
	count := endPage - startPage + 1
	nijieroChWg.Add(count)
	//要爬取的页面
	for i := startPage; i <= endPage; i++ {
		go nijieroChSpiderRun(i)
	}
	nijieroChWg.Wait()
	nijieroChWg2.Wait()
	//nijieroChSaveFileWg.Wait()
}

func nijieroChSpiderRun(page int) {
	defer nijieroChWg.Done()

	//根据传递的page得到爬取的url
	url := fmt.Sprintf(nijieroChUrl, strconv.Itoa(page))
	contents := util.ExplainUrl(url)
	//以“post hentry”来切分获取的连接
	contentsArr := strings.Split(contents, "post hentry")
	for k, v := range contentsArr {
		//跳过第一个
		if k == 0 {
			continue
		}
		//将v以"切分
		vArr := strings.Split(v, "\"")
		secondUrl := vArr[2]
		title := vArr[4]
		//开始爬取子页面
		go nijieroChSpiderImageRun(secondUrl, title)
	}
}

//第二次抓取页面
func nijieroChSpiderImageRun(link, title string) {
	defer nijieroChWg2.Done()

	nijieroChWg2.Add(1)

	//最终保存路径
	saveDir := fmt.Sprintf("%s%s\\", nijieroChRootDir, title)
	//这里在对应的文件夹下新建对应标题的文件夹
	makeDirErr := util.MakeDir(saveDir)
	if makeDirErr != nil {
		fmt.Println(fmt.Sprintf("目录%s创建失败，爬取%s页面失败", saveDir, link))
		return
	}
	fmt.Println("我准备再次访问链接了！", link)

	contents := util.ExplainUrl(link)
	//以"imageTitle"来切分获取的连接
	contentsArr := strings.Split(contents, "imageTitle")
	i := 0
	for k, v := range contentsArr {
		//跳过第一个
		if k <= 1 || (strings.Index(v, "script") != -1 && strings.Index(v, "jpg") == -1) {
			continue
		}
		i++
		//将v以"切分
		vArr := strings.Split(v, "\"")
		imageUrl := vArr[2]
		fmt.Println(imageUrl)
		nijieroChSaveFileWg.Add(1)
		util.SaveFile(imageUrl, saveDir, "", &nijieroChSaveFileWg)
	}
	fmt.Println("爬取页面成功，爬取图片数量：", i)
	return
}
