package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "regexp"
)

var linkPattern = rgx(`<a href="(.+?)"`)

func rgx(pattern string) *regexp.Regexp {
  //returns regexp version of pattern
  return regexp.MustCompile(pattern)
}

type Callback func()

type Queue struct {
    elements []string
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

type Crawler struct {
    keyWords []string
    links Queue
    getNextLinks func([]byte)[][][]byte
    seen map[string]bool
}

func (c * Crawler) getPages() {
    for !c.links.isEmpty() {

        fmt.Println(len(c.links.elements))

        nextPage := c.next()
        nextLinks := c.getNextLinks(nextPage)

        for _, link := range nextLinks {
            c.links.push(string(link[1]))
        }
    }
}

func (c * Crawler) done(page string) (isDone bool) {
    return
}

func getNextLinks(page []byte) (nextLinks [][][]byte) {
    nextLinks = linkPattern.FindAllSubmatch(page, -1)
    return
}

func (c * Crawler) next() (nextPage []byte) {
    nextLink := c.links.pop()

    if (!c.seen[nextLink]) {
      c.seen[nextLink] = true
      resp, err := http.Get(nextLink)

      errHandler(err, func() {
        defer resp.Body.Close()
        body, err := ioutil.ReadAll(resp.Body)
        errHandler(err, func() {
          nextPage = body
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

func main() {
    start := "http://www.livestrong.com/"
    keyWords := []string{"health", "food", "diet", "superfood", "fat", "carbohydrates", "protein"}
    seen := make(map[string]bool)
    links := Queue{[]string{start}}
    crawler := Crawler{keyWords, links, getNextLinks, seen}

    crawler.getPages()

    //np := NutritionProfile {Map{initMap([]string{"fat", "carbohydrates", "protein"}), initMap}}

}
