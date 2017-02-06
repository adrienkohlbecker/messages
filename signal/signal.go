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
	"<img src=\"/images/emoji/apple/1f37e.png\" class=\"emoji\" title=\":champagne:\">":                    "🍾",
	"<img src=\"/images/emoji/apple/1f440.png\" class=\"emoji\" title=\":eyes:\">":                         "👀",
	"<img src=\"/images/emoji/apple/1f448.png\" class=\"emoji\" title=\":point_left:\">":                   "👈",
	"<img src=\"/images/emoji/apple/1f449.png\" class=\"emoji\" title=\":point_right:\">":                  "👉",
	"<img src=\"/images/emoji/apple/1f44c.png\" class=\"emoji\" title=\":ok_hand:\">":                      "👌",
	"<img src=\"/images/emoji/apple/1f44d.png\" class=\"emoji\" title=\":+1:\">":                           "👍",
	"<img src=\"/images/emoji/apple/1f44f.png\" class=\"emoji\" title=\":clap:\">":                         "👏",
	"<img src=\"/images/emoji/apple/1f481-1f3fb.png\" class=\"emoji\" title=\"information_desk_person\">":  "💁🏻",
	"<img src=\"/images/emoji/apple/1f494.png\" class=\"emoji\" title=\":broken_heart:\">":                 "💔",
	"<img src=\"/images/emoji/apple/1f49a.png\" class=\"emoji\" title=\":green_heart:\">":                  "💚",
	"<img src=\"/images/emoji/apple/1f4ac.png\" class=\"emoji\" title=\":speech_balloon:\">":               "💬",
	"<img src=\"/images/emoji/apple/1f600.png\" class=\"emoji\" title=\":grinning:\">":                     "😀",
	"<img src=\"/images/emoji/apple/1f601.png\" class=\"emoji\" title=\":grin:\">":                         "😁",
	"<img src=\"/images/emoji/apple/1f602.png\" class=\"emoji\" title=\":joy:\">":                          "😂",
	"<img src=\"/images/emoji/apple/1f603.png\" class=\"emoji\" title=\":smiley:\">":                       "😃",
	"<img src=\"/images/emoji/apple/1f605.png\" class=\"emoji\" title=\":sweat_smile:\">":                  "😅",
	"<img src=\"/images/emoji/apple/1f609.png\" class=\"emoji\" title=\":wink:\">":                         "😉",
	"<img src=\"/images/emoji/apple/1f60a.png\" class=\"emoji\" title=\":blush:\">":                        "😊",
	"<img src=\"/images/emoji/apple/1f60b.png\" class=\"emoji\" title=\":yum:\">":                          "😋",
	"<img src=\"/images/emoji/apple/1f60d.png\" class=\"emoji\" title=\":heart_eyes:\">":                   "😍",
	"<img src=\"/images/emoji/apple/1f60e.png\" class=\"emoji\" title=\":sunglasses:\">":                   "😎",
	"<img src=\"/images/emoji/apple/1f60f.png\" class=\"emoji\" title=\":smirk:\">":                        "😏",
	"<img src=\"/images/emoji/apple/1f613.png\" class=\"emoji\" title=\":sweat:\">":                        "😓",
	"<img src=\"/images/emoji/apple/1f615.png\" class=\"emoji\" title=\":confused:\">":                     "😕",
	"<img src=\"/images/emoji/apple/1f617.png\" class=\"emoji\" title=\":kissing:\">":                      "😗",
	"<img src=\"/images/emoji/apple/1f618.png\" class=\"emoji\" title=\":kissing_heart:\">":                "😘",
	"<img src=\"/images/emoji/apple/1f61b.png\" class=\"emoji\" title=\":stuck_out_tongue:\">":             "😛",
	"<img src=\"/images/emoji/apple/1f61c.png\" class=\"emoji\" title=\":stuck_out_tongue_winking_eye:\">": "😜",
	"<img src=\"/images/emoji/apple/1f61e.png\" class=\"emoji\" title=\":disappointed:\">":                 "😞",
	"<img src=\"/images/emoji/apple/1f620.png\" class=\"emoji\" title=\":angry:\">":                        "😠",
	"<img src=\"/images/emoji/apple/1f621.png\" class=\"emoji\" title=\":rage:\">":                         "😡",
	"<img src=\"/images/emoji/apple/1f622.png\" class=\"emoji\" title=\":cry:\">":                          "😢",
	"<img src=\"/images/emoji/apple/1f625.png\" class=\"emoji\" title=\":disappointed_relieved:\">":        "😥",
	"<img src=\"/images/emoji/apple/1f628.png\" class=\"emoji\" title=\":fearful:\">":                      "😨",
	"<img src=\"/images/emoji/apple/1f629.png\" class=\"emoji\" title=\":weary:\">":                        "😩",
	"<img src=\"/images/emoji/apple/1f62c.png\" class=\"emoji\" title=\":grimacing:\">":                    "😬",
	"<img src=\"/images/emoji/apple/1f62d.png\" class=\"emoji\" title=\":sob:\">":                          "😭",
	"<img src=\"/images/emoji/apple/1f631.png\" class=\"emoji\" title=\":scream:\">":                       "😱",
	"<img src=\"/images/emoji/apple/1f633.png\" class=\"emoji\" title=\":flushed:\">":                      "😳",
	"<img src=\"/images/emoji/apple/1f635.png\" class=\"emoji\" title=\":dizzy_face:\">":                   "😵",
	"<img src=\"/images/emoji/apple/1f636.png\" class=\"emoji\" title=\":no_mouth:\">":                     "😶",
	"<img src=\"/images/emoji/apple/1f642.png\" class=\"emoji\" title=\":slightly_smiling_face:\">":        "🙂",
	"<img src=\"/images/emoji/apple/1f643.png\" class=\"emoji\" title=\":upside_down_face:\">":             "🙃",
	"<img src=\"/images/emoji/apple/1f64a.png\" class=\"emoji\" title=\":speak_no_evil:\">":                "🙊",
	"<img src=\"/images/emoji/apple/1f64c-1f3fd.png\" class=\"emoji\" title=\"raised_hands\">":             "🙌🏽",
	"<img src=\"/images/emoji/apple/1f911.png\" class=\"emoji\" title=\":money_mouth_face:\">":             "🤑",
	"<img src=\"/images/emoji/apple/1f913.png\" class=\"emoji\" title=\":nerd_face:\">":                    "🤓",
	"<img src=\"/images/emoji/apple/1f914.png\" class=\"emoji\" title=\":thinking_face:\">":                "🤔",
	"<img src=\"/images/emoji/apple/1f917.png\" class=\"emoji\" title=\":hugging_face:\">":                 "🤗",
	"<img src=\"/images/emoji/apple/1f918.png\" class=\"emoji\" title=\":the_horns:\">":                    "🤘",
	"<img src=\"/images/emoji/apple/2122.png\" class=\"emoji\" title=\":tm:\">":                            "™",
	"<img src=\"/images/emoji/apple/2639.png\" class=\"emoji\" title=\":white_frowning_face:\">":           "☹️",
	"<img src=\"/images/emoji/apple/26c4.png\" class=\"emoji\" title=\":snowman_without_snow:\">":          "⛄",
	"<img src=\"/images/emoji/apple/2764.png\" class=\"emoji\" title=\":heart:\">":                         "❤️",
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
			case "audio":
				filename = fmt.Sprintf("%s.mp3", name)
			case "gif":
				filename = fmt.Sprintf("%s.gif", name)
			case "png":
				filename = fmt.Sprintf("%s.png", name)
			case "aac":
				filename = fmt.Sprintf("%s.aac", name)
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
