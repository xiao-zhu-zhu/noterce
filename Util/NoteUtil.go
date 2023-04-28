package Util

import (
	"github.com/antchfx/htmlquery"
	"github.com/parnurzeal/gorequest"
	"net/url"
)

func WriteNote(noteid string, note string) {
	agent := gorequest.New()
	agent.Post("https://note.ms/"+noteid).AppendHeader("Referer", "https://note.ms/"+noteid).Send("&t=" + url.QueryEscape(note)).End()
}

func GetNote(noteid string) string {
	agent := gorequest.New()
	resp, _, _ := agent.Get("https://note.ms/" + noteid).End()

	if resp == nil {
		return ""
	}

	parse, err := htmlquery.Parse(resp.Body)
	if err != nil {
		return "解析失败"
	}
	one := htmlquery.FindOne(parse, "/html/body/div[1]/div[1]/div/div/textarea")

	text := htmlquery.InnerText(one)
	return text
}
