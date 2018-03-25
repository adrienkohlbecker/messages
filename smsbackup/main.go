package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrienkohlbecker/messages/model"
	"github.com/vincent-petithory/dataurl"
)

type SMSes struct {
	XMLName xml.Name `xml:"smses"`
	Count   int      `xml:"count,attr"`
	SMSes   []SMS    `xml:"sms"`
	MMSes   []MMS    `xml:"mms"`
}

var Entities = map[string]string{
	"&#10;":            "\n",
	"&#13;":            "\r",
	"&#55356;&#57096;": "🌈",
	"&#55356;&#57158;": "🍆",
	"&#55356;&#57218;": "🎂",
	"&#55356;&#57224;": "🎈",
	"&#55356;&#57225;": "🎉",
	"&#55356;&#57226;": "🎊",
	"&#55356;&#57252;": "🎤",
	"&#55356;&#57270;": "🎶",
	"&#55356;&#57339;": "🏻🏻", // Emoji Modifier Fitzpatrick Type-1-2
	"&#55357;&#56338;": "🐒",
	"&#55357;&#56396;": "👌",
	"&#55357;&#56397;": "👍",
	"&#55357;&#56401;": "👑",
	"&#55357;&#56447;": "👿",
	"&#55357;&#56459;": "💋",
	"&#55357;&#56463;": "💏",
	"&#55357;&#56473;": "💙",
	"&#55357;&#56476;": "💜",
	"&#55357;&#56489;": "💩",
	"&#55357;&#56490;": "💪",
	"&#55357;&#56565;": "📵",
	"&#55357;&#56725;": "🖕🏻",
	"&#55357;&#56740;": "🖤",
	"&#55357;&#56832;": "😀",
	"&#55357;&#56833;": "😁",
	"&#55357;&#56834;": "😂",
	"&#55357;&#56835;": "😃",
	"&#55357;&#56836;": "😄",
	"&#55357;&#56837;": "😅",
	"&#55357;&#56838;": "😆",
	"&#55357;&#56839;": "😇",
	"&#55357;&#56841;": "😉",
	"&#55357;&#56842;": "😊",
	"&#55357;&#56843;": "😋",
	"&#55357;&#56845;": "😍",
	"&#55357;&#56846;": "😎",
	"&#55357;&#56847;": "😏",
	"&#55357;&#56848;": "😐",
	"&#55357;&#56849;": "😑",
	"&#55357;&#56850;": "😒",
	"&#55357;&#56851;": "😓",
	"&#55357;&#56853;": "😕",
	"&#55357;&#56855;": "😗",
	"&#55357;&#56856;": "😘",
	"&#55357;&#56858;": "😚",
	"&#55357;&#56859;": "😛",
	"&#55357;&#56860;": "😜",
	"&#55357;&#56862;": "😞",
	"&#55357;&#56864;": "😠",
	"&#55357;&#56866;": "😢",
	"&#55357;&#56873;": "😩",
	"&#55357;&#56875;": "😫",
	"&#55357;&#56876;": "😬",
	"&#55357;&#56877;": "😭",
	"&#55357;&#56879;": "😯",
	"&#55357;&#56880;": "😰",
	"&#55357;&#56881;": "😱",
	"&#55357;&#56883;": "😳",
	"&#55357;&#56885;": "😵",
	"&#55357;&#56887;": "😷",
	"&#55357;&#56891;": "😻",
	"&#55357;&#56898;": "🙂",
	"&#55357;&#56899;": "🙃",
	"&#55357;&#56900;": "🙄",
	"&#55357;&#56904;": "🙈",
	"&#55357;&#56906;": "🙊",
	"&#55357;&#56962;": "🚂",
	"&#55358;&#56595;": "🤓",
	"&#55358;&#56596;": "🤔",
	"❤&#65039;":        "❤️",
	"&#55356;&#57169;": "🍑",
	"&#55357;&#56840;": "😈",
	"&#55357;&#56857;": "😙",
	"☺&#65039;":        "☺️",
	"&#55357;&#56392;": "👈",
	"&#55358;&#56710;": "🦆",
	"&#55357;&#56869;": "😥",
	"&#55357;&#56357;": "🐥",
	"&#55357;&#57025;": "🛁",
}

var Addresses = map[string]string{
	"33785529239":          "+33785529239",
	"695959773":            "+33695959773",
	"33695218383":          "+33695218383",
	"33628255966":          "+33628255966",
	"(+33)615528008":       "+33615528008",
	"insert-address-token": "+33661779655",
}

type SMS struct {
	XMLName       xml.Name `xml:"sms"`
	Protocol      int      `xml:"protocol,attr"`
	Address       string   `xml:"address,attr"`
	Date          int64    `xml:"date,attr"`
	Type          int      `xml:"type,attr"`
	Subject       string   `xml:"subject,attr"`
	Body          string   `xml:"body,attr"`
	Toa           string   `xml:"toa,attr"`
	ScToa         string   `xml:"sc_toa,attr"`
	ServiceCenter string   `xml:"service_center,attr"`
	Read          int      `xml:"read,attr"`
	Status        int      `xml:"status,attr"`
	Locked        int      `xml:"locked,attr"`
}

type MMS struct {
	XMLName      xml.Name `xml:"mms"`
	MsgBox       int      `xml:"msg_box,attr"`
	Date         int64    `xml:"date,attr"`
	AddressTilde string   `xml:"address,attr"`
	Parts        MMSParts `xml:"parts"`
	Addrs        MMSAddrs `xml:"addrs"`
}

type MMSParts struct {
	XMLName xml.Name  `xml:"parts"`
	Parts   []MMSPart `xml:"part"`
}

type MMSPart struct {
	XMLName xml.Name `xml:"part"`
	CT      string   `xml:"ct,attr"`
	FN      string   `xml:"fn,attr"`
	Text    string   `xml:"text,attr"`
	Data    string   `xml:"data,attr"`
}

type MMSAddrs struct {
	XMLName xml.Name  `xml:"addrs"`
	Addrs   []MMSAddr `xml:"addr"`
}

type MMSAddr struct {
	XMLName xml.Name `xml:"addr"`
	Address string   `xml:"address,attr"`
	Type    int      `xml:"type,attr"`
}

func main() {

	file := "/Users/ak/Dropbox/SignalPlaintextBackup.xml"

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatal(err)
	}

	var msgs model.Messages

	f, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range Entities {
		f = bytes.Replace(f, []byte(k), []byte(v), -1)
	}

	if bytes.Contains(f, []byte("&#")) {
		i := bytes.Index(f, []byte("&#"))
		log.Fatalf("XML contains unknown entities %s", f[i-32:i+32])
	}

	v := SMSes{}
	err = xml.Unmarshal(f, &v)
	if err != nil {
		log.Fatal(err)
	}

	for _, sms := range v.SMSes {

		address := parseAddr(sms.Address)

		var sender string
		if sms.Type != 2 {
			sender = address
		}

		msgs = append(msgs, &model.Message{
			ID:          "",
			Sender:      sender,
			Content:     strings.TrimSpace(sms.Body),
			Timestamp:   time.Unix(sms.Date/1000, (sms.Date-(sms.Date/1000)*1000)*1000000).In(loc),
			Sent:        (sms.Type == 2),
			Attachments: make([]*model.Attachment, 0),
			Group:       address,
			Kind:        "sms",
		})
	}

	for _, mms := range v.MMSes {

		var sender string
		for _, addr := range mms.Addrs.Addrs {
			if addr.Type == 137 {
				sender = parseAddr(addr.Address)
			}
		}

		var content string
		for _, part := range mms.Parts.Parts {
			if part.CT == "text/plain" {
				content = part.Text
			}
		}

		var attachments []*model.Attachment
		for _, part := range mms.Parts.Parts {
			if part.CT != "text/plain" && part.CT != "application/smil" {

				url, loopErr := dataurl.DecodeString("data:" + part.CT + ";base64," + part.Data)
				if loopErr != nil {
					log.Fatal(loopErr)
				}

				var extension string
				var kind string

				switch part.CT {
				case "image/jpeg":
					extension = ".jpg"
					kind = "img"
				case "image/png":
					extension = ".png"
					kind = "img"
				case "text/x-vCard":
					extension = ".vcard"
					kind = "vcard"
				case "video/3gpp":
					extension = ".3gp"
					kind = "video"
				default:
					log.Fatalf("unknown extension for %s", part.CT)
				}

				path := part.FN
				if strings.HasPrefix(part.FN, "part-") || part.FN == "null" {
					path = fmt.Sprintf("%x", md5.Sum(url.Data)) + extension
				}

				log.Printf("Wrote %s", path)
				path = filepath.Join("/Users/ak/Desktop/messages-store/mediatmp", path)

				loopErr = ioutil.WriteFile(path, url.Data, os.FileMode(0644))
				if loopErr != nil {
					log.Fatal(loopErr)
				}

				attachments = append(attachments, &model.Attachment{
					Kind: kind,
					URL:  path,
				})

			}
		}

		groupElt := strings.Split(mms.AddressTilde, "~")
		for i := range groupElt {
			groupElt[i] = parseAddr(groupElt[i])
		}
		sort.Strings(groupElt)
		group := strings.Join(groupElt, "~")

		msgs = append(msgs, &model.Message{
			ID:          "",
			Sender:      sender,
			Content:     strings.TrimSpace(content),
			Timestamp:   time.Unix(mms.Date/1000, (mms.Date-(mms.Date/1000)*1000)*1000000).In(loc),
			Sent:        (mms.MsgBox == 2),
			Attachments: attachments,
			Group:       group,
			Kind:        "mms",
		})

	}

	sort.Sort(msgs)

	msgsJSON, err := json.MarshalIndent(msgs, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filepath.Join("/Users/ak/Desktop/messages-store", strings.Replace(filepath.Base(file), filepath.Ext(file), ".json", 1)), msgsJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//pp.Println(v.MMSes)

}

func parseAddr(s string) string {
	address := strings.Replace(s, " ", "", -1)

	new, ok := Addresses[address]
	if ok {
		address = new
	} else {
		if !strings.HasPrefix(address, "+") {
			if strings.HasPrefix(address, "0") {
				address = "+33" + address[1:len(address)]
			} else {
				log.Printf("[WARN]: Unknown number format: %s", address)
			}
		}
	}

	return address
}
