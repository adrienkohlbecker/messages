package main

import (
	"log"
	"sort"

	"github.com/adrienkohlbecker/messages/model"
)

func main() {

	msgs, err := model.Load("/Users/ak/Desktop/messages-store/signal.json")
	if err != nil {
		log.Fatal(err)
	}

	sort.Sort(msgs)

	err = msgs.Write("/Users/ak/Desktop/messages-store/signal.json")
	if err != nil {
		log.Fatal(err)
	}

}
