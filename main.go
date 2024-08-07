package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Stock struct {
	Company string `json:"company"`
	Price   string `json:"price"`
	Change  string `json:"change"`
}

func main() {
	ticker := []string{
		"MSFT", "IBM", "GE", "UNP", "COST", "MCD", "V", "WMT", "DIS",
		"MMM", "INTC", "AXP", "AAPL", "BA", "CSCO", "GS", "JPM", "CRM", "VZ",
	}

	stocks := []Stock{}

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div#quote-header-info", func(e *colly.HTMLElement) {
		stock := Stock{
			Company: e.ChildText("h1"),
			Price:   e.ChildText("fin-streamer[data-field='regularMarketPrice']"),
			Change:  e.ChildText("fin-streamer[data-field='regularMarketChangePercent'] span"),
		}
		fmt.Println("Company:", stock.Company)
		fmt.Println("Price:", stock.Price)
		fmt.Println("Change:", stock.Change)

		stocks = append(stocks, stock)
	})

	for _, t := range ticker {
		c.Visit("https://finance.yahoo.com/quote/" + t + "/")
	}

	c.Wait()

	fmt.Println(stocks)

	file, err := os.Create("stocks.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	headers := []string{"company", "price", "change"}
	writer.Write(headers)
	for _, stock := range stocks {
		record := []string{stock.Company, stock.Price, stock.Change}
		writer.Write(record)
	}
	writer.Flush()
}
