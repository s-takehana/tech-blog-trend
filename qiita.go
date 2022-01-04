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

		section := s.Find("section")

		var qiita qiita

		// Section 1
		section1 := section.First()

		a := section1.Find("a")
		qiita.Username = a.Text()
		avatar, _ := a.Find("img").Attr("src")
		qiita.Avatar = avatar

		p := section1.Find("p").Last()
		v := strings.Split(p.Text(), " ")
		t, err := time.Parse("2006-01-02", v[2])
		if err != nil {
			log.Fatal(err)
		}
		qiita.Datetime = t

		// Section 2
		section2 := section1.Next()

		a = section2.Find("h2").Find("a")
		qiita.Title = a.Text()
		href, _ := a.Attr("href")
		qiita.URL = href

		// Section 3
		section3 := section2.Next()

		div := section3.Children().Find("div")
		div1 := div.First()

		var tags []string
		div1.Find("a").Each(func(ii int, ss *goquery.Selection) {
			tags = append(tags, ss.Text())
		})
		qiita.Tags = tags

		lgtm, _ := strconv.Atoi(div1.Next().Text())
		qiita.LGTM = lgtm

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
