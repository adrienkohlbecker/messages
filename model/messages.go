package model

import (
	"encoding/json"
	"io/ioutil"
)

type Messages []*Message

func (msgs Messages) Len() int      { return len(msgs) }
func (msgs Messages) Swap(i, j int) { msgs[i], msgs[j] = msgs[j], msgs[i] }

func (msgs Messages) Less(i, j int) bool {
	if msgs[i].Group == msgs[j].Group {
		return msgs[i].Timestamp.Before(msgs[j].Timestamp)
	}
	return msgs[i].Group < msgs[j].Group
}

func Load(path string) (Messages, error) {

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var list Messages
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, err
	}

	return list, nil

}

func (msgs *Messages) Write(path string) error {

	msgsJSON, err := json.MarshalIndent(msgs, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, msgsJSON, 0644)
	if err != nil {
		return err
	}
	return nil

}
