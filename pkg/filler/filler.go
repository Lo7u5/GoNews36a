package filler

import (
	"GoNews36a/pkg/dbase/postgresql"
	"encoding/xml"
	strip "github.com/grokify/html-strip-tags-go"
	"io"
	"net/http"
	"strings"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Chanel  Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	Link        string `xml:"link"`
}

// Parse получает данные из rss, разкодирует xml и возвращает массив новостей
func Parse(url string) ([]postgresql.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var f Feed
	err = xml.Unmarshal(b, &f)
	if err != nil {
		return nil, err
	}
	var data []postgresql.Post
	for _, item := range f.Chanel.Items {
		var p postgresql.Post
		p.Title = item.Title
		p.Content = item.Description
		p.Content = strip.StripTags(p.Content)
		p.Link = item.Link
		item.PubDate = strings.ReplaceAll(item.PubDate, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubDate)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}
		data = append(data, p)
	}
	return data, nil
}
