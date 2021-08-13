package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Llongfile)

	log.Println("start")

	s := flag.String("s", "", "site zenn|dev")
	f := flag.String("f", "./data.json", "file path")
	flag.Parse()

	var url string
	if *s == "zenn" {
		url = "https://zenn.dev/api/articles"
	} else if *s == "dev" {
		url = "https://dev.to/api/articles?top=1"
	} else {
		log.Fatalf("site error: %s", *s)
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	file, err := os.Create(*f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err = io.Copy(file, res.Body); err != nil {
		log.Fatal(err)
	}

	log.Println("finish")
}
