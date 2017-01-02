package main

import (
    "fmt"
    "net/http"
    "strings"
    "io/ioutil"
    "golang.org/x/net/html"
)

type Callback func()

type Queue struct {
    elements []string
}

func (q * Queue) isEmpty() (empty bool) {
    empty = len(q.elements) == 0
    return
}

func (q * Queue) pop() (nextElement string) {
    nextElement = q.elements[0]

    q.elements = q.elements[1:]

    return
}

func (q * Queue) push(newElement string) {
    q.elements = append(q.elements, newElement)
}


type Map struct {
    elements map[string]string
    initMap func([]string)map[string]string
}

func (m * Map) contains(field string) bool {
    var ok bool
    _, ok = m.elements[field]
    return ok
}

func initMap(fields []string) map[string]string {
    elements := make(map[string]string)
    for i := range fields {
      elements[fields[i]] = ""
    }

    return elements
}


type Crawler struct {
    keyWords []string
    links Queue
    seen map[string]bool
    text [][]string
}

func (c * Crawler) getPages() {
    for !c.links.isEmpty() {

        nextPage := c.next()

        data := getTextFromPage(nextPage)

        if (len(data["text"]) > 0) {
          fmt.Println(data["text"])
          c.text = append(c.text, data["text"])
        }

        nextLinks := data["link"]

        for _, l := range nextLinks {
          c.links.push(l)
        }
    }
}

func (c * Crawler) done(page string) (isDone bool) {
    return
}

func (c * Crawler) next() (nextPage string) {
    nextLink := c.links.pop()

    if (!c.seen[nextLink]) {
      c.seen[nextLink] = true
      resp, err := http.Get(nextLink)

      errHandler(err, func() {
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        errHandler(err, func() {
          nextPage = string(body)
        })
      })

    }

    return
}

type NutritionProfile struct {
    profile Map
}

func errHandler(err error, callback Callback) {
  if (err != nil) {
    fmt.Println("error")
    fmt.Println(err)
  } else {
    callback()
  }
}

func getTextFromPage(page string) (pageData map[string][]string) {
  doc, err := html.Parse(strings.NewReader(page))

  if (err == nil) {
    pageData = inspectParsedHTML(doc, make(map[string][]string))
  }

  return
}

func inspectParsedHTML(n *html.Node, pageData map[string][]string) map[string][]string {

    if n.Type == html.ElementNode && n.Data == "a" {
        for _, a := range n.Attr {
            if (a.Key == "href") {
              pageData["link"] = append(pageData["link"], a.Val)
            }
        }
    }

    if (n.Type == html.TextNode) {
      text := n.Data
      pageData["text"] = append(pageData["text"], text)
    }


    for c := n.FirstChild; c != nil; c = c.NextSibling {
        inspectParsedHTML(c, pageData)
    }

    return pageData
}

func main() {
    start := "http://www.livestrong.com/"
    keyWords := []string{"health", "food", "diet", "superfood", "fat", "carbohydrates", "protein"}
    seen := make(map[string]bool)
    links := Queue{[]string{start}}
    text := make([][]string, 0)
    crawler := Crawler{keyWords, links, seen, text}

    crawler.getPages()

    //np := NutritionProfile {Map{initMap([]string{"fat", "carbohydrates", "protein"}), initMap}}

}
