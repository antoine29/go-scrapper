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

func sendMessageToTgChannel(tgBotToken string, tgChannelName string, message string) error {
  encodedMsg := url.QueryEscape(message)
  url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", tgBotToken, tgChannelName, encodedMsg)
  err := HttpGet(url)
  if err != nil {
    return err
  }

  return nil
}

func HttpGet(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusAccepted) {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.New(fmt.Sprintf("Unsuccessful %d http response. Cannot read response body.\n", res.StatusCode))
		}

		return errors.New(fmt.Sprintf("Unsuccessful %d http response.\n%s", res.StatusCode, resBody))
	}

	return nil
}

