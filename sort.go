package main

import (
	"log"
	"sort"

	"github.com/adrienkohlbecker/messages/model"
)

func main() {

	msgs, err := model.Load("/Users/adrien/Desktop/messages-store/whatsapp.json")
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(msgs)

	err = msgs.Write("/Users/adrien/Desktop/messages-store/whatsapp.json")
	if err != nil {
		log.Fatal(err)
	}

}
