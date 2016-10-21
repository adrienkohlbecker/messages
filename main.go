package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/k0kubun/pp"
)

type MessageType int

const (
	MessageTypeUnknown MessageType = iota
	MessageTypeReceived
	MessageTypeSent
)

type Message struct {
	Contact  string
	Type     MessageType
	Date     time.Time
	Contents string
}

type SMSes struct {
	XMLName xml.Name `xml:"smses"`
	Count   int      `xml:"count,attr"`
	SMSes   []SMS    `xml:"sms"`
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
	"&#55357;&#56396;": "👌",
	"&#55357;&#56397;": "👍",
	"&#55357;&#56459;": "💋",
	"&#55357;&#56463;": "💏",
	"&#55357;&#56473;": "💙",
	"&#55357;&#56476;": "💜",
	"&#55357;&#56489;": "💩",
	"&#55357;&#56490;": "💪",
	"&#55357;&#56565;": "📵",
	"&#55357;&#56725;": "🖕🏻",
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
	"&#55357;&#56859;": "😛",
	"&#55357;&#56860;": "😜",
	"&#55357;&#56866;": "😢",
	"&#55357;&#56873;": "😩",
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
	"&#55357;&#56962;": "🚂",
}

var Addresses = map[string]string{
	"33785529239":    "+33785529239",
	"695959773":      "+33695959773",
	"33695218383":    "+33695218383",
	"33628255966":    "+33628255966",
	"(+33)615528008": "+33615528008",
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

type Messages map[string][]*Message

func (m Messages) Add(msg *Message) {

	_, ok := m[msg.Contact]
	if !ok {
		m[msg.Contact] = make([]*Message, 0)
	}
	m[msg.Contact] = append(m[msg.Contact], msg)

}

func main() {

	files := []string{"g4.xml", "signal.xml", "soshphone.xml", "telegram.xml", "whatsapp.xml", "sms-20161021040000.xml"}

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatal(err)
	}

	var msgs = Messages{}

	for _, file := range files {

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

			address := strings.Replace(sms.Address, " ", "", -1)

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

			msg := &Message{
				Contact:  address,
				Contents: sms.Body,
				Type:     MessageType(sms.Type),
				Date:     time.Unix(sms.Date/1000, sms.Date-(sms.Date/1000)*1000).In(loc),
			}
			msgs.Add(msg)
		}

	}

	pp.Println(msgs["+33622199552"])

}
