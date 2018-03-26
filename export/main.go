package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/k0kubun/pp"

	"github.com/adrienkohlbecker/messages/formatter"
	"github.com/adrienkohlbecker/messages/model"
)

var msgs model.Messages
var _ = pp.Println

func main() {

	files := []string{
		"/Users/ak/Desktop/messages-store/signal.json",
		"/Users/ak/Desktop/messages-store/SignalPlaintextBackup.json",
		"/Users/ak/Desktop/messages-store/soshphone.json",
		"/Users/ak/Desktop/messages-store/g4.json",
		"/Users/ak/Desktop/messages-store/fb_Putzfrauplayboys.json",
		"/Users/ak/Desktop/messages-store/fb_Paris-ci les sorties.json",
		"/Users/ak/Desktop/messages-store/g5.json",
		"/Users/ak/Desktop/messages-store/telegram.json",
		"/Users/ak/Desktop/messages-store/whatsapp.json",
		"/Users/ak/Desktop/messages-store/facebook.json",
		"/Users/ak/Desktop/messages-store/sms-20170924174214.json",
	}

	for _, file := range files {
		m, err := model.Load(file)
		if err != nil {
			log.Fatal(err)
		}
		msgs = append(msgs, m...)
	}

	file, err := os.Open("/Users/ak/Desktop/messages-store/aliases.csv")
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(file)
	aliases := make(map[string]string)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		aliases[record[1]] = record[0]

	}

	for i, msg := range msgs {

		if val, ok := aliases[msg.Sender]; ok {
			msg.Sender = val
		}

		splitted := strings.Split(msg.Group, "~")
		for j, name := range splitted {
			if val, ok := aliases[name]; ok {
				splitted[j] = val
			}
		}

		sort.Strings(splitted)

		msg.Group = strings.Join(splitted, "~")
		msgs[i] = msg

	}

	err = os.RemoveAll("/Users/ak/export")
	if err != nil {
		panic(err)
	}

	byGroup := msgs.ByGroup()
	for _, grouped := range byGroup {
		err := exportGroup(grouped)
		if err != nil {
			panic(err)
		}
	}

}

func exportGroup(grp *model.Grouped) error {

	grpDir := filepath.Join("/Users/ak/export", grp.Group)
	jsonPath := filepath.Join(grpDir, fmt.Sprintf("%s.json", grp.Group))
	htmlPath := filepath.Join(grpDir, fmt.Sprintf("%s.html", grp.Group))

	err := os.MkdirAll(grpDir, 0755)
	if err != nil {
		return err
	}

	for _, msg := range grp.Messages {

		for _, att := range msg.Attachments {

			filename := filepath.Base(att.URL)
			newDir := filepath.Join("/Users/ak/export", "media", att.Kind)

			err := os.MkdirAll(newDir, 0755)
			if err != nil {
				return err
			}

			newPath := filepath.Join(newDir, filename)
			err = Copy(att.URL, newPath)
			if err != nil {
				return err
			}

			att.URL = newPath

		}

	}

	sort.Sort(grp.Messages)

	err = grp.Messages.Write(jsonPath)
	if err != nil {
		return err
	}

	file, err := os.Create(htmlPath)
	if err != nil {
		return err
	}

	for _, msg := range grp.Messages {
		for _, att := range msg.Attachments {
			att.URL = fmt.Sprintf("file://%s", att.URL)
		}
	}

	err = formatter.Format(grp.Messages, file)
	if err != nil {
		return err
	}

	file.Close()
	if err != nil {
		return err
	}

	return nil

}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
