package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type Product struct {
	ProductName string `json:"productName"`
	ProductPrice string `json:"productPrice"`
	ProductLink string `json:"productLink"`
}

func main() {
	c := colly.NewCollector()

	games := []Product{}

	// check all pages
	c.OnHTML("li.pagination--next", func(h *colly.HTMLElement) {
		c.Visit(h.Request.AbsoluteURL(h.ChildAttr("a", "href")))
	})

	// check product details
	c.OnHTML("div.productitem", func(h *colly.HTMLElement) {
		i := Product {
			ProductName: strings.Replace(h.ChildText("div.productitem--info h2.productitem--title a"), "\u0026", "", -1),
			ProductPrice: h.ChildText("div.price--main span.money"),
			ProductLink: "https://gameline.ph" + h.ChildAttr("div.productitem a.productitem--image-link", "href"),
		}

		games = append(games, i);
	})

	// Display all url to scrape
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Visit main url
	c.Visit("https://gameline.ph/search?type=article%2Cpage%2Cproduct&q=ps5*+games*")

	// Parse to JSON
	data, err := json.MarshalIndent(games, " ", "")

	if err != nil {
		log.Fatal()
	}

	// Create file
	os.WriteFile("scraped_data", []byte(data), 0666)

	if err != nil {
		fmt.Println(err)
	}
}
