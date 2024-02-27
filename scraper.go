package mscraper

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

type Song struct {
	Title  string `json:"title"`
	Artist string `json:"artist"`
}

type Genre string

const (
	KBALLAD = "0100"
	KDANCE  = "0200"
	KHIPHOP = "0300"
	KRB     = "0400"
	KINDY   = "0500"
	KROCK   = "0600"
	KTROT   = "0700"
	KBLUSE  = "0800"
	POP     = "0900"
	ROCK    = "1000"
	ELEC    = "1100"
	HIPHOP  = "1200"
	RB      = "1300"
	BLUSE   = "1400"
)

// Methods to be called outside as a library ----------------------
// V1
func GetNewestAlbumFromMelon() []Album {
	c := colly.NewCollector()
	return scrapeNewestAlubumMelon(c)
}

func GetNewestHipHopFromMelon() []Song {
	c := colly.NewCollector()
	return scrapeNewestHipHopSongsMelon(c)
}

// V2
func GetNewestSongsMelon(genreCode string) []Song {
	c := colly.NewCollector()
	return scrapeNewestSongsMelon(c, genreCode)
}

// Scrape Methods -------------------------------------------------
func scrapeNewestSongsMelon(c *colly.Collector, genreCode string) []Song {
	var songs []Song
	var currentSong Song

	c.OnHTML("div.wrap_song_info", func(h *colly.HTMLElement) {
		h.ForEach("div", func(_ int, div *colly.HTMLElement) {
			// Find the nested <a> tag within the <span> tag
			a := div.ChildText("span a")

			// Check if the <a> tag is the first or second one based on your HTML structure
			if div.Index == 0 {
				currentSong.Title = a
			} else if div.Index == 1 {
				currentSong.Artist = a

				// Add to list of songs
				songs = append(songs, currentSong)

				// Reset the current song for the next iteration
				currentSong = Song{}
			}
		})
	})

	c.Visit("https://www.melon.com/genre/song_list.htm?gnrCode=GN" + genreCode)
	return songs
}

func scrapeNewestAlubumMelon(c *colly.Collector) []Album {
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

func scrapeNewestHipHopSongsMelon(c *colly.Collector) []Song {

	var songs []Song
	var currentSong Song

	c.OnHTML("div.wrap_song_info", func(h *colly.HTMLElement) {
		h.ForEach("div", func(_ int, div *colly.HTMLElement) {
			// Find the nested <a> tag within the <span> tag
			a := div.ChildText("span a")

			// Check if the <a> tag is the first or second one based on your HTML structure
			if div.Index == 0 {
				currentSong.Title = a
			} else if div.Index == 1 {
				currentSong.Artist = a

				// Add to list of songs
				songs = append(songs, currentSong)

				// Reset the current song for the next iteration
				currentSong = Song{}
			}
		})
	})

	c.Visit("https://www.melon.com/genre/song_list.htm?gnrCode=GN0300&dtlGnrCode=")
	return songs
}

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

	os.WriteFile(fileName, content, 0644)
}
