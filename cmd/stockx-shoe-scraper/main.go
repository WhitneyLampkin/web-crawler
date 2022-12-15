package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Shoe struct {
	Name  string
	Price string
}

func main() {
	file, err := os.Create("export.csv")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"Name", "Price"}
	writer.Write(headers)

	c := colly.NewCollector(
		colly.AllowedDomains("stockx.com"),
	)

	c.OnHTML(".css-1ibvugw-GridProductTileContainer", func(e *colly.HTMLElement) {
		shoe := Shoe{}
		shoe.Name = e.ChildText(".css-3lpefb")
		shoe.Price = e.ChildText(".css-9ryi0c")
		row := []string{shoe.Name, shoe.Price}
		writer.Write(row)
		fmt.Println(shoe.Name, "|", shoe.Price)
	})

	c.OnHTML(".css-12da55z-PaginationButton", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(nextPage)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.StatusCode)
	})

	c.Visit("https://stockx.com/sneakers")

	fmt.Println("Done")
}
