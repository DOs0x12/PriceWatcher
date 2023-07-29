package domain

import (
	"errors"
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
		priceErr := "cannot parse the body to an HTML document: %w"

		return 0.00, fmt.Errorf(priceErr, err)
	}

	tag := "td"
	re := regexp.MustCompile(regstring)

	data := doTraverse(doc, tag, re)

	if data == "" {
		traverceErr := fmt.Sprintf("the document does not have a price value with the tag: %v", tag)

		return 0.00, errors.New(traverceErr)
	}

	return getPrice(data)
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

func getPrice(data string) (float32, error) {
	trailingSymCnt := 17
	dataLen := len(data)

	if (dataLen == 0) || (dataLen <= trailingSymCnt) {
		lenErr := fmt.Sprintf("the length of the data string is not valid: %v", dataLen)

		return 0.00, errors.New(lenErr)
	}

	data = data[trailingSymCnt : dataLen-1]
	price, err := strconv.ParseFloat(data, 32)
	if err != nil {
		parseErr := fmt.Sprintf("cannot parse the string data: %v", data)

		return 0.00, errors.New(parseErr)
	}

	return float32(price), nil
}
