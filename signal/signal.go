package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/adrienkohlbecker/messages/model"
	"github.com/k0kubun/pp"
)

type SrcMessage struct {
	ID          string              `json:"id"`
	Sender      string              `json:"sender"`
	Content     string              `json:"content"`
	Timestamp   string              `json:"timestamp"`
	Sent        bool                `json:"sent"`
	Attachments []*model.Attachment `json:"attachments"`
	Group       string              `json:"group"`
	Kind        string              `json:"kind"`
}

var HTML = map[string]string{
	"&amp;": "&",
	"<br>":  "\n",
	"<img src=\"/images/emoji/apple/1f44d.png\" class=\"emoji\" title=\":+1:\">":               "👍",
	"<img src=\"/images/emoji/apple/1f600.png\" class=\"emoji\" title=\":grinning:\">":         "😀",
	"<img src=\"/images/emoji/apple/1f602.png\" class=\"emoji\" title=\":joy:\">":              "😂",
	"<img src=\"/images/emoji/apple/1f603.png\" class=\"emoji\" title=\":smiley:\">":           "😃",
	"<img src=\"/images/emoji/apple/1f609.png\" class=\"emoji\" title=\":wink:\">":             "😉",
	"<img src=\"/images/emoji/apple/1f618.png\" class=\"emoji\" title=\":kissing_heart:\">":    "😘",
	"<img src=\"/images/emoji/apple/1f61b.png\" class=\"emoji\" title=\":stuck_out_tongue:\">": "😛",
	"<img src=\"/images/emoji/apple/1f61e.png\" class=\"emoji\" title=\":disappointed:\">":     "😞",
	"<img src=\"/images/emoji/apple/1f620.png\" class=\"emoji\" title=\":angry:\">":            "😠",
	"<img src=\"/images/emoji/apple/1f62c.png\" class=\"emoji\" title=\":grimacing:\">":        "😬",
	"<img src=\"/images/emoji/apple/1f633.png\" class=\"emoji\" title=\":flushed:\">":          "😳",
	"<img src=\"/images/emoji/apple/1f60d.png\" class=\"emoji\" title=\":heart_eyes:\">":       "😍",
	"<img src=\"/images/emoji/apple/1f613.png\" class=\"emoji\" title=\":sweat:\">":            "😓",
	"<img src=\"/images/emoji/apple/1f628.png\" class=\"emoji\" title=\":fearful:\">":          "😨",
	"<img src=\"/images/emoji/apple/1f911.png\" class=\"emoji\" title=\":money_mouth_face:\">": "🤑",
	"<img src=\"/images/emoji/apple/1f494.png\" class=\"emoji\" title=\":broken_heart:\">":     "💔",
	"<img src=\"/images/emoji/apple/1f60a.png\" class=\"emoji\" title=\":blush:\">":            "😊",
}
var LinkRegex = regexp.MustCompile("<a href=\"([^\"]+)\" target=\"_blank\">[^<]+</a>")

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: signal PATH")
		os.Exit(1)
	}

	source := os.Args[1]

	if !FileExists(source) {
		log.Fatalf("Unable to find source at %s", source)
	}

	list, err := Load(source)
	if err != nil {
		log.Fatal(err)
	}

	for _, msg := range list {

		for k, v := range HTML {
			msg.Content = strings.Replace(msg.Content, k, v, -1)
		}

		msg.Content = LinkRegex.ReplaceAllString(msg.Content, "$1")

		if strings.Contains(msg.Content, "<img") {
			log.Fatalf("Unmatched emoji in: %s", msg.Content)
		}
		if strings.Contains(msg.Content, "<") {
			log.Fatalf("Unmatched HTML in: %s", msg.Content)
		}

		for _, attachment := range msg.Attachments {

			name := filepath.Base(attachment.URL)
			var filename string

			switch attachment.Kind {
			case "img":
				filename = fmt.Sprintf("%s.jpeg", name)
			case "video":
				filename = fmt.Sprintf("%s.mp4", name)
			default:
				log.Fatal("unknown type")
			}

			path := filepath.Join("/Users/adrien/Downloads", filename)
			newPath := filepath.Join("/Users/adrien/Dropbox/Applications/Messages/media", filename)
			if !FileExists(path) && !FileExists(newPath) {
				log.Fatalf("File does not exist: %s (%s)", attachment.URL, path)
			}

			if FileExists(path) {
				err = os.Rename(path, newPath)
				if err != nil {
					log.Fatal(err)
				}
			}

			attachment.URL = newPath

		}

	}

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatal(err)
	}

	var result model.Messages

	for _, msg := range list {

		c, errFor := strconv.ParseInt(msg.Timestamp, 10, 64)
		if errFor != nil {
			log.Fatal(errFor)
		}

		s := c / 1000
		ms := (c - s*1000) * 1000000

		t := time.Unix(s, ms).In(loc)

		result = append(result, &model.Message{
			ID:          msg.ID,
			Sender:      msg.Sender,
			Content:     msg.Content,
			Timestamp:   t,
			Sent:        msg.Sent,
			Attachments: msg.Attachments,
			Group:       msg.Group,
			Kind:        msg.Kind,
		})

	}

	sort.Sort(result)

	resultJSON, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join("/Users/adrien/Dropbox/Applications/Messages", filepath.Base(source)), resultJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println(filepath.Join("/Users/adrien/Dropbox/Applications/Messages", filepath.Base(source)))

}

func Load(path string) ([]*SrcMessage, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list []*SrcMessage
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}

	return list, nil

}

func FileExists(path string) bool {

	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true

}
