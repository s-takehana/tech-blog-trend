package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type qiita struct {
	Username string    `json:"username"`
	Avatar   string    `json:"avatar"`
	Datetime time.Time `json:"datetime"`
	Title    string    `json:"title"`
	URL      string    `json:"url"`
	Tags     []string  `json:"tags"`
	LGTM     int       `json:"lgtm"`
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	log.Println("start")

	f := flag.String("f", "./qiita.json", "file path")
	flag.Parse()

	res, err := http.Get("https://qiita.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var qiitas []qiita

	doc.Find("article").Each(func(i int, s *goquery.Selection) {

		header := s.Find("header")
		a := header.Find("a")

		qiita := qiita{}

		a.Each(func(ii int, ss *goquery.Selection) {
			href, _ := ss.Attr("href")
			text := ss.Text()

			if !strings.HasPrefix(href, "/organizations/") && len(text) > 0 {
				qiita.Username = text
			}
		})

		avatar, _ := a.Find("img").Attr("src")
		qiita.Avatar = avatar

		datetime, _ := header.Find("time").Attr("datetime")
		t, err := time.Parse(time.RFC3339, datetime)
		if err != nil {
			log.Fatal(err)
		}
		qiita.Datetime = t

		a = s.Find("h2").Find("a")
		qiita.Title = a.Text()
		href, _ := a.Attr("href")
		qiita.URL = href

		footer := s.Find("footer")

		var tags []string
		footer.Find("a").Each(func(ii int, ss *goquery.Selection) {
			href, _ := ss.Attr("href")

			if strings.HasPrefix(href, "/tags/") {
				tags = append(tags, ss.Text())
			}
		})
		qiita.Tags = tags

		footer.Find("div").EachWithBreak(func(ii int, ss *goquery.Selection) bool {
			if lgtm, err := strconv.Atoi(ss.Text()); err == nil {
				qiita.LGTM = lgtm
				return false
			}
			return true
		})

		qiitas = append(qiitas, qiita)
	})

	file, err := os.Create(*f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err = json.NewEncoder(file).Encode(qiitas); err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(qiitas)
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	log.Printf("\n%s\n", out.String())

	log.Println("finish")
}
