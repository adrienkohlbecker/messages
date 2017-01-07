package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/adrienkohlbecker/messages/model"
	"github.com/k0kubun/pp"
)

var msgRegex = regexp.MustCompile("(\\d{1,2}/\\d{1,2}/\\d{1,2}, \\d{1,2}:\\d{1,2}) - ([\\w\\s]+): (.*)")
var attRegex = regexp.MustCompile("(.*) \\(file attached\\)(?:\\n(.*))?")
var loc *time.Location

func init() {

	l, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatal(err)
	}
	loc = l

}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: whatsapp PATH")
		os.Exit(1)
	}

	source := os.Args[1]

	if !FileExists(source) {
		log.Fatalf("Unable to find source at %s", source)
	}

	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var msgs model.Messages
	var previousTs time.Time
	var addedSec time.Duration

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "are now secured with end-to-end encryption") {
			continue
		}

		msg, loopErr := Parse(filepath.Base(source), line)
		if loopErr != nil {
			log.Fatal(loopErr)
		}

		if msg.Timestamp.Equal(previousTs) {
			msg.Timestamp = msg.Timestamp.Add(addedSec + time.Second)
			addedSec = addedSec + time.Second
		} else {
			addedSec = 0
			previousTs = msg.Timestamp
		}

		msgs = append(msgs, msg)
	}

	err = scanner.Err()
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(msgs)

	resultJSON, err := json.MarshalIndent(msgs, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join("/Users/adrien/Dropbox/Applications/Messages", filepath.Base(source)), resultJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

	pp.Println(filepath.Join("/Users/adrien/Dropbox/Applications/Messages", filepath.Base(source)))

}

func FileExists(path string) bool {

	_, err := os.Stat(path)
	if err != nil {
		return false
	}

	return true

}

func Parse(group string, line string) (*model.Message, error) {

	matches := msgRegex.FindStringSubmatch(line)

	pp.Println(matches)

	t, err := time.Parse("1/2/06, 15:04", matches[1])
	if err != nil {
		return nil, err
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)

	content := matches[3]
	content = strings.Replace(content, "\\n", "\n", -1)

	var attachments []*model.Attachment

	attMatches := attRegex.FindStringSubmatch(content)
	if len(attMatches) != 0 {
		attachments = append(attachments, &model.Attachment{Kind: "img", URL: "/Users/adrien/Dropbox/Applications/Messages/media/" + attMatches[1]})
		content = attMatches[2]
	}

	msg := &model.Message{
		ID:          "",
		Sender:      matches[2],
		Content:     content,
		Timestamp:   t,
		Sent:        false,
		Attachments: attachments,
		Group:       group,
		Kind:        "whatsapp",
	}

	if msg.Sender == "Adrien Kohlbecker" {
		msg.Sender = ""
		msg.Sent = true
	}

	return msg, nil

}
