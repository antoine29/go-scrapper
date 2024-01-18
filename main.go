package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

var targets []string = []string{
	"2010990022392",
	"20386975",
	"203.9",
}

var tgToken, tgTokenIsSet = os.LookupEnv("TG_TOKEN")
var tgChannelName, tgChannelNameIsSet = os.LookupEnv("TG_CHANNEL_NAME")

func main() {
  if !tgTokenIsSet || !tgChannelNameIsSet {
    slog.Error("TG env vars not set")
    os.Exit(1)
  }
  
  slog.Info(fmt.Sprintf("Scrapping started at %s", time.Now()))
  // sendMessageToTgChannel(
  //   tgToken, 
  //   tgChannelName, 
  //   fmt.Sprintf("Scrapping started at: %s", time.Now()),
  // )

	collector := colly.NewCollector()
  collector.SetRequestTimeout(5 * time.Minute)
	collector.OnHTML("p", func(element *colly.HTMLElement) {
		// fmt.Printf("%+v\n", e)
		text := strings.ToLower(element.Text)
		if targetMatchs(text) {
      slog.Info(fmt.Sprintf("Found: \n %s", text))
      sendMessageToTgChannel(tgToken, tgChannelName, text)
		} 
	})

  if err := collector.Visit("https://thor.organojudicial.gob.bo/"); err != nil {
    errorMessage := fmt.Sprintf("Error scrapping. \n %s", err.Error())
    slog.Error(errorMessage)
    os.Exit(1)
  }

  slog.Info("Scrapping finished")
}

func targetMatchs(text string) bool {
	for _, targuet := range targets {
		if strings.Contains(text, targuet) {
			return true
		}
	}

	return false
}

func sendMessageToTgChannel(tgBotToken string, tgChannelName string, message string) {
	encodedMsg := url.QueryEscape(message)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", tgBotToken, tgChannelName, encodedMsg)
	err := HttpGet(url)
	if err != nil {
    slog.Error(fmt.Sprintf("Error sending message to TG channel. \n %s", err.Error()))
	}
}

func HttpGet(url string) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}

	if !(res.StatusCode == http.StatusOK || res.StatusCode == http.StatusAccepted) {
		resBody, err := io.ReadAll(io.Reader(res.Body))
		if err != nil {
			return errors.New(fmt.Sprintf("Unsuccessful %d http response. Cannot read response body.\n", res.StatusCode))
		}

		return errors.New(fmt.Sprintf("Unsuccessful %d http response.\n%s", res.StatusCode, resBody))
	}

	return nil
}

