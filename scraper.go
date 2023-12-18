package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/gocolly/colly"
)

type Album struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

// Able to retrieve new-release album from "melon.com" website
func main() {
	c := colly.NewCollector()
	albums := scrapeNewestAlubum(c)
	writeDataToJSON("music.json", albums)
}

// Scrape Methods -------------------------------------------------

func scrapeNewestAlubum(c *colly.Collector) []Album {
	var albums []Album

	c.OnHTML("div.info", func(h *colly.HTMLElement) {
		albumName := h.ChildText("a.album_name")
		artistName := h.ChildText("span.checkEllipsis a.artist_name")

		if albumName != "" || artistName != "" {

			albumInstance := Album{
				Name:   removeBetweenBrackets(albumName),
				Artist: removeBetweenBrackets(artistName),
			}
			albums = append(albums, albumInstance)
		}

	})
	c.Visit("https://www.melon.com/new/album/index.htm")
	return albums
}

// func scrapeNewestHipHopSongs(c *colly.Collector) {

// 	c.OnHTML("div.info", func(h *colly.HTMLElement) {
// 		albumName := h.ChildText("a.album_name")
// 		artistName := h.ChildText("span.checkEllipsis a.artist_name")

// 	})

// 	c.Visit("https://www.melon.com/new/album/index.htm")
// }

// Util Methods -----------------------------------------

func removeBetweenBrackets(input string) string {
	re := regexp.MustCompile(`\(.*?\)`)
	return re.ReplaceAllString(input, "")
}

func writeDataToJSON(fileName string, data any) {
	content, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error occured during json marshal process")
		return
	}

	os.WriteFile("music.json", content, 0644)
}
