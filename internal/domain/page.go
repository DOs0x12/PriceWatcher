package domain

import (
	"fmt"
	"io"
	"regexp"
	"strconv"

	"golang.org/x/net/html"
)

const regstring = `(^ покупка: [0-9]{4,5}\.[0-9][0-9])`

type Extractor interface {
	ExtractPrice(body io.ReadCloser) (float32, error)
}

type PriceExtractor struct{}

func (ext PriceExtractor) ExtractPrice(body io.ReadCloser) (float32, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return 0.00, fmt.Errorf("cannot parse the body to an HTML document: %w", err)
	}

	tag := "td"
	re := regexp.MustCompile(regstring)

	data := doTraverse(doc, tag, re)

	if data == "" {
		return 0.00, fmt.Errorf("the document does not have a price value with the tag: %v", tag)
	}

	return getPrice(data), nil
}

func doTraverse(doc *html.Node, tag string, re *regexp.Regexp) string {
	var traverse func(n *html.Node, tag string, re *regexp.Regexp) string

	traverse = func(n *html.Node, tag string, re *regexp.Regexp) string {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if isNodeWithPriceForBuying(c, tag, re) {
				return c.Data
			}

			data := traverse(c, tag, re)

			if data != "" {
				return data
			}
		}

		return ""
	}

	return traverse(doc, tag, re)
}

func isNodeWithPriceForBuying(n *html.Node, tag string, re *regexp.Regexp) bool {
	return n.Type == html.TextNode && n.Parent.Data == tag && re.MatchString(n.Data)
}

func getPrice(data string) float32 {
	re := regexp.MustCompile(`([0-9]{4,5}\.[0-9][0-9])`)
	match := re.FindStringSubmatch(data)[0]
	price, _ := strconv.ParseFloat(match, 32)

	return float32(price)
}
