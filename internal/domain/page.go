package domain

import (
	"io"

	"golang.org/x/net/html"
)

type Extractor interface {
	ExtractPrice(body io.ReadCloser) float32
}

type PriceExtractor struct{}

func (ext PriceExtractor) ExtractPrice(body io.ReadCloser) float32 {
	doc, err := html.Parse(body)

	if err != nil {
		//TODO: implement error handling
		return 0.00
	}

	tag := "td"

	tagPs := doTraverse(doc, tag)
	tags := *tagPs
	t := tags[0]
	_ = t

	return 5.55
}

func doTraverse(doc *html.Node, tag string) *[]string {
	var data []string

	var traverse func(n *html.Node, data *[]string, tag string) *html.Node

	traverse = func(n *html.Node, data *[]string, tag string) *html.Node {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if isNodeWithPriceForBuying(c, tag) {
				*data = append(*data, c.Data)
			}

			res := traverse(c, data, tag)

			if res != nil {
				return res
			}
		}

		return nil
	}

	traverse(doc, &data, tag)

	return &data
}

func isNodeWithPriceForBuying(n *html.Node, tag string) bool {
	return n.Type == html.TextNode && n.Parent.Data == tag
}
