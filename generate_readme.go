package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func getHotTopic(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("http get error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("status code error: %d %s\n", resp.StatusCode, resp.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var topics []string
	// Find the review items
	doc.Find("table > tbody td.ranktop").Each(func(i int, selection *goquery.Selection) {
		topic := selection.Siblings().Find("a").Text()
		fmt.Println(i, topic)
		topics = append(topics, topic)
	})

	return topics, nil
}

func generateREADME(topics []string) {
	readme := ""
	for index, topic := range topics {
		// n. XXX
		temp := fmt.Sprintf("%d. %s", index+1, topic)
		topic = strings.Replace(topic, " ", "%20", -1)
		temp = temp + "[:link:](https://s.weibo.com/weibo?q=%23" + topic + "%23&Refer=top)\n"

		if index == 10 {
			temp = fmt.Sprintf("<details>\n<summary>%d ~ %d</summary>\n\n%s", index+1, Min(index+10, len(topics)), temp)
		} else if index == len(topics)-1 {
			temp = fmt.Sprintf("%s</details>", temp)
		} else if index >= 11 && index%10 == 0 {
			temp = fmt.Sprintf("</details>\n<details>\n<summary>%d ~ %d</summary>\n\n%s", index+1, Min(index+10, len(topics)), temp)
		} else {
		}

		readme += temp
	}

	writeStringToFile(readme)
}

func writeStringToFile(text string) {
	// write the whole body at once
	err := ioutil.WriteFile("README.md", []byte(text), 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	topics, err := getHotTopic("https://s.weibo.com/top/summary?cate=realtimehot")
	if err != nil {
		return
	}
	generateREADME(topics)
}
