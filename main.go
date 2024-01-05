package main

import (
  "fmt"
  "strings"
	"github.com/gocolly/colly/v2"
)

var targets []string = []string{
  "2010990022392",
  "20386975",
  "203.9",
}

func main(){
	c := colly.NewCollector()

	c.OnHTML("p", func(element *colly.HTMLElement) {
		// fmt.Printf("%+v\n", e)
    text := strings.ToLower(element.Text)
    if targetMatchs(text) {
      sendTgMessage(text)
    }
	})

	c.Visit("https://thor.organojudicial.gob.bo/")
}

func targetMatchs(text string) bool {
  for _, targuet := range targets {
    if strings.Contains(text, targuet) {
      return true
    }
  }

  return false
}

func sendTgMessage(text string) {

}

