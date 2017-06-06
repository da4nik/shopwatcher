package wildberries

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/da4nik/shopwatcher/types"
)

const (
	// ShopName shop identifier
	ShopName = "wildberries"

	// ShopRegexp regexp to check url
	ShopRegexp = "https://www.wildberries.ru/catalog/[0-9]+/detail.aspx"
)

// Parse - parses wildberries products
func Parse(url string) (types.Product, error) {
	prod := types.Product{}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return prod, err
	}

	// Name and price
	doc.Find(".good-header").Each(func(i int, s *goquery.Selection) {
		prod.Name = cleanupString(s.Find("h1").Text())

		s.Find(".price meta").Each(func(i int, meta *goquery.Selection) {
			prop, exists := meta.Attr("itemprop")
			if exists && prop == "price" {
				price, _ := meta.Attr("content")
				prod.Price = price
			}

			if exists && prop == "priceCurrency" {
				currency, _ := meta.Attr("content")
				prod.Currency = currency
			}
		})
	})

	// Sizes
	avaialble := false
	doc.Find("#pp-sizes label.j-size").Each(func(i int, size *goquery.Selection) {
		disabled := size.HasClass("disabled")
		avaialble = avaialble || !disabled
		currentSize := size.Find("span").First().Text()
		prod.Sizes = append(prod.Sizes, types.Size{
			Name:      currentSize,
			Available: !disabled,
		})
	})

	prod.Available = avaialble
	prod.Shop = ShopName
	prod.URL = url

	return prod, nil
}

func cleanupString(str string) string {
	return strings.Trim(str, " \n")
}
