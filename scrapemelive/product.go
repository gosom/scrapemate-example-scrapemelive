package scrapemelive

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Product struct {
	Title            string
	Price            string
	ShortDescription string
	Sku              string
	Categories       []string
	Tags             []string
}

func (o Product) CsvHeaders() []string {
	return []string{
		"title",
		"price",
		"short_description",
		"sku",
		"categories",
		"tags",
	}
}

func (o Product) CsvRow() []string {
	return []string{
		o.Title,
		o.Price,
		o.ShortDescription,
		o.Sku,
		strings.Join(o.Categories, ","),
		strings.Join(o.Tags, ","),
	}
}

func parseProduct(doc *goquery.Document) Product {
	return Product{
		Title:            parseTitle(doc),
		Price:            parsePrice(doc),
		ShortDescription: parseShortDescription(doc),
		Sku:              parseSku(doc),
		Categories:       parseCategories(doc),
		Tags:             parseTags(doc),
	}
}

func parseTitle(doc *goquery.Document) string {
	return doc.Find("h1.product_title").Text()
}

func parsePrice(doc *goquery.Document) string {
	return doc.Find("p.price").Text()
}

func parseShortDescription(doc *goquery.Document) string {
	return doc.Find("div.woocommerce-product-details__short-description>p").Text()
}

func parseSku(doc *goquery.Document) string {
	return doc.Find("span.sku").Text()
}

func parseCategories(doc *goquery.Document) []string {
	var categories []string
	doc.Find("div.product_meta > span.posted_in > a").Each(func(i int, s *goquery.Selection) {
		categories = append(categories, s.Text())
	})
	return categories
}

func parseTags(doc *goquery.Document) []string {
	var tags []string
	doc.Find("div.product_meta > span.tagged_as > a").Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})
	return tags
}
