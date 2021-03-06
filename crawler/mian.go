package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/tidwall/gjson"
)

type resp struct {
	channel_id  int
	view_count  string
	like_count  string
	author_name string
	author_id   string
	bvid        string
	card_type   string
}

var DB *gorm.DB
var err error

//InitDB 是一个初始化的函数
func initDB() *gorm.DB {
	DB, err = gorm.Open("mysql", "root:zk2824895143@(localhost)/user?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	//defer DB.Close()
	return DB
}

func httpget(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()

	//循环读取网页，传出给调用者
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			fmt.Println("读取网页完成")
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		} //

		result += string(buf[:n])
	}
	return
}

func work() {
	fmt.Println("正在爬取目标网页。。。")

	for k := 100; k <= 200; k++ {
		fmt.Printf("这是第%d个频道", k)
		url2 := strconv.Itoa(k)
		url := "https://api.bilibili.com/x/web-interface/web/channel/multiple/list?channel_id=" + url2 + "&sort_type=hot&page_size=30"
		result, err := httpget(url)
		if err != nil {
			fmt.Println("http get err", err)

		}

		for i := 1; i <= 30; i++ {

			num := strconv.Itoa(i)
			// fmt.Println(num)
			temp := resp{

				channel_id:  k,
				view_count:  gjson.Get(result, "data.list."+num+".view_count").String(),
				like_count:  gjson.Get(result, "data.list."+num+".like_count").String(),
				author_name: gjson.Get(result, "data.list."+num+".author_name").String(),
				author_id:   gjson.Get(result, "data.list."+num+".author_id").String(),
				bvid:        gjson.Get(result, "data.list."+num+".bvid").String(),
				card_type:   gjson.Get(result, "data.list."+num+".card_type").String(),
			}
			if temp.card_type != "" {
				// fmt.Println(temp)
				if err := DB.Table("bilibili_copy1").Create(&temp).Error; err != nil {
					fmt.Println(err)
				}
			}

		}

	}
}

func main() {
	initDB()
	work()
}
