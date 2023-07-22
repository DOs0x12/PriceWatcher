package domain

import (
	"io"

	"golang.org/x/net/html"
)

type Extractor interface {
	ExtractRate(body io.ReadCloser) float32
}

type RateExtractor struct{}

func (ext RateExtractor) ExtractRate(body io.ReadCloser) float32 {
	doc, err := html.Parse(body)

	if err != nil {
		//TODO: implement error handling
		return 0.00
	}

	tag := "td"

	tagValue := doTraverse(doc, tag)
	tagData := *tagValue
	t := tagData[0]
	_ = t
	return 5.55
}

func doTraverse(doc *html.Node, tag string) *[]string {
	var data []string

	var traverse func(n *html.Node, data *[]string, tag string) *html.Node

	traverse = func(n *html.Node, data *[]string, tag string) *html.Node {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode && c.Parent.Data == tag {
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
