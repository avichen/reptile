package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//Item is a struct for json response
type Item struct {
	ID string
	P  string
	M  string
}

func getPrice(client *http.Client, url string) string {
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml; q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Encoding", "text/html")
	// request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4,zh-TW;q=0.2")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")

	resp, _ := client.Do(request)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var p []Item
	p = make([]Item, 0)
	err := json.Unmarshal(body, &p)
	if err != nil {
		fmt.Printf("%T\n%s\n%#v\n", err, err, err)
	}

	// fmt.Printf("id: %s ,price: %s\n", p[0].ID, p[0].P)
	time.Sleep(time.Second * 10)
	return p[0].P
}

func getItemsByPage(client *http.Client, url string) {
	//http://list.jd.com/list.html?cat=670%2C671%2C672&go=0
	//http://list.jd.com/list.html?cat=1316%2C1383%2C1401&go=0
	request, _ := http.NewRequest("GET", url, nil)

	request.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml; q=0.9,image/webp,*/*;q=0.8")
	request.Header.Set("Accept-Encoding", "text/html")
	// request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4,zh-TW;q=0.2")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36")

	resp, _ := client.Do(request)

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	var cnt int32
	var itemName string
	var price string
	doc.Find("li.gl-item").Each(func(i int, s *goquery.Selection) {
		// fmt.Printf("gl-item: %d  ", i)
		// Single item finder
		s.Find("div.gl-i-wrap.j-sku-item").Each(func(m int, b *goquery.Selection) {
			// fmt.Printf("tab-content-item: %d\n", k)
			skuCode, _ := b.Find("div.p-operate>a.p-o-btn.contrast.J_contrast").Attr("data-sku")
			// fmt.Println(skuCode)
			itemName = b.Find("div.p-name>a>em").Text()
			itemPrice := "http://p.3.cn/prices/mgets?skuIds=J_" + skuCode + "&type=1"
			price = getPrice(client, itemPrice)
			// fmt.Println(itemPrice)
			fmt.Printf("The item %s name is '%s', price is %s \n", skuCode, itemName, price)
			cnt++
		})
		// Tab Items finder
		s.Find("div.gl-i-tab-content>div.tab-content-item.j-sku-item").Each(func(k int, a *goquery.Selection) {
			// fmt.Printf("tab-content-item: %d\n", k)
			skuCode, _ := a.Find("div.p-operate>a.p-o-btn.contrast.J_contrast").Attr("data-sku")
			// fmt.Println(skuCode)
			itemName = a.Find("div.p-name>a>em").Text()
			itemPrice := "http://p.3.cn/prices/mgets?skuIds=J_" + skuCode + "&type=1"
			price = getPrice(client, itemPrice)
			fmt.Printf("The item %s name is '%s', price is %s \n", skuCode, itemName, price)
			cnt++
		})

	})

	fmt.Println(cnt)

	node := doc.Find("span.p-num>a").Last()
	href, exists := node.Attr("href")
	if exists {
		nextPage := "http://list.jd.com" + href
		fmt.Println(nextPage)
		// time.Sleep(time.Second * 10)
		// getItemsByPage(client, nextPage)

	}

}

func main() {
	client := &http.Client{}
	startURL := "http://list.jd.com/list.html?cat=1316%2C1383%2C1401&go=0"
	getItemsByPage(client, startURL)

	//https://coderwall.com/p/4c2zig/decode-top-level-json-array-into-a-slice-of-structs-in-golang
}
