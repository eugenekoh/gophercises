package linkparser

import (
	"container/list"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link is a data structure that holds link's href and text content in a HTML document
type Link struct {
	href string
	text string
}

// ParseLinks is a function that parses HTML to extract out all links.
// It uses a BFS algorithm and ignores nested links
func ParseLinks(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	links := make([]Link, 0)
	linkNodes := findLinkNodes(doc)

	for _, linkNode := range linkNodes {
		links = append(links, extractLinkInfo(linkNode))
	}

	return links, nil
}

func findLinkNodes(doc *html.Node) []*html.Node {

	linkNodes := make([]*html.Node, 0)

	// BFS using doubly linked list
	queue := list.New()
	queue.PushBack(doc)

	for queue.Len() != 0 {
		elem := queue.Front()
		node := elem.Value.(*html.Node)

		// iterate through all siblings
		for node != nil {

			switch {
			// find a tags
			case node.Type == html.ElementNode && node.Data == "a":
				linkNodes = append(linkNodes, node)

			default:
				// add first child to queue
				if node.FirstChild != nil {
					queue.PushBack(node.FirstChild)
				}
			}

			node = node.NextSibling
		}

		queue.Remove(elem)
	}

	return linkNodes
}

func extractLinkInfo(node *html.Node) Link {
	var link Link

	// extract href value, assumes href exists
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link.href = attr.Val
		}
	}

	// extract text , dfs
	link.text = getText(node)
	return link
}

func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	var ret string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		ret += getText(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
