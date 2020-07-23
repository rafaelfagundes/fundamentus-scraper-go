package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	// This library prints structured data on the console
	// Usage: spew.Dump(data)
	// "github.com/davecgh/go-spew/spew"
	"github.com/gocolly/colly"
)

func getStockData() []byte {
	result := make([]map[string]interface{}, 0)

	c := colly.NewCollector(colly.CacheDir("./cache"))

	// Table with stock data
	c.OnHTML("table", func(e *colly.HTMLElement) {

		// Get the table headers
		headers := make([]string, 0)
		e.ForEach("thead tr th a", func(_ int, el *colly.HTMLElement) {
			headers = append(headers, el.Text)
		})

		// Get the table rows
		e.ForEach("tbody tr", func(_ int, tr *colly.HTMLElement) {

			// { "key": dynamic value }
			row := make(map[string]interface{})

			// Table cells
			tr.ForEach("td", func(_ int, td *colly.HTMLElement) {

				// Convert from pt_BR currency/number format
				value := strings.ReplaceAll(td.Text, ".", "")
				value = strings.ReplaceAll(value, ",", ".")

				// Remove % sign and convert to float64
				if strings.Contains(value, "%") {
					value = strings.ReplaceAll(value, "%", "")
					f, _ := strconv.ParseFloat(value, 64)
					numberValue := f / 100
					value = fmt.Sprintf("%f", numberValue)
				}

				// If currency or numeric
				if strings.Contains(value, ".") {
					f, _ := strconv.ParseFloat(value, 64)
					row[headers[td.Index]] = f
				} else { // Else is the stock symbol
					row[headers[td.Index]] = value
				}
			})

			// Append row to result data
			result = append(result, row)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.fundamentus.com.br/resultado.php")

	jsonString, _ := json.Marshal(result)
	return jsonString
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func saveToDisk(json []byte) {
	err := ioutil.WriteFile("output/data.json", json, 0644)
	check(err)
}

func main() {
	fmt.Println("Wait while the application retrieves the stock data...")
	json := getStockData()
	saveToDisk(json)
	fmt.Println("Done! The file was written to ./output/data.json")
}
