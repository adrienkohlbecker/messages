package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/adrienkohlbecker/messages/model"
)

type Grouped struct {
	Group    string
	Messages model.Messages
}

var byGroup = make(map[string]*Grouped)
var tmpl *template.Template

func main() {

	msgs, err := model.Load("/Users/adrien/Desktop/messages-store/signal_g4.json")
	if err != nil {
		log.Fatal(err)
	}

	// bytes, err := ioutil.ReadFile("/Users/adrien/Desktop/messages-store/temp.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// var list [][]string
	// err = json.Unmarshal(bytes, &list)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// loc, err := time.LoadLocation("Europe/Paris")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// for _, item := range list {
	//
	// 	t, loopErr := time.Parse("Jan 2, 15:04:05", item[1])
	// 	if loopErr != nil {
	// 		log.Fatal(loopErr)
	// 	}
	// 	t = time.Date(2016, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	//
	// 	var att []*model.Attachment
	// 	if len(item) == 4 {
	// 		kind := "img"
	// 		if strings.HasSuffix(item[3], ".mp4") {
	// 			kind = "video"
	// 		}
	// 		att = append(att, &model.Attachment{
	// 			URL:  item[3],
	// 			Kind: kind,
	// 		})
	// 	}
	//
	// 	sender := item[0]
	// 	if sender == "Adrien" {
	// 		sender = ""
	// 	}
	//
	// 	msg := &model.Message{
	// 		ID:          "",
	// 		Sender:      sender,
	// 		Content:     item[2],
	// 		Timestamp:   t,
	// 		Sent:        item[0] == "Adrien",
	// 		Attachments: att,
	// 		Group:       "+33629448767~+33695218383",
	// 		Kind:        "signal",
	// 	}
	//
	// 	msgs = append(msgs, msg)
	//
	// }
	//
	// sort.Sort(msgs)
	//
	// msgs.Write("/Users/adrien/Desktop/messages-store/signal_g5_full.json")

	for _, msg := range msgs {

		_, ok := byGroup[msg.Group]
		if !ok {
			byGroup[msg.Group] = &Grouped{Group: msg.Group, Messages: make(model.Messages, 0)}
		}

		byGroup[msg.Group].Messages = append(byGroup[msg.Group].Messages, msg)

	}

	tmplStr := `
<style type="text/css">
	body {
	    font-family: Helvetica Neue;
	    font-size: 16px;
	}

	.group-title {
	    clear: both;
	    margin-top: 30px;
	}

	.message {
	    padding: 5px;
	    border: 1px solid #d6d6d6;
	    margin-top: 10px;
	    border-radius: 5px;
	}

	.message-body {
	    margin:0;
	}

	.message-time {
	    display: block;
	    font-size: 0.8em;
	    padding-top: 5px;
	}

	.message-sent {
	    float:right;
	    background-color: #e8e8e8;
	}

	.message-sent .message-time {
	    color: #8c8c8c;
	}

	.message-received {
	    float:left;
	    background-color: #4e97ff;
	    color: white;
	}

	.message-received .message-time {
	    color: #d8d8d8;
	}

	.message {
	    clear:both;
	}

	.message-sender {
	    font-size: 0.9em;
	    margin: 0;
	    margin-bottom: 5px;
	}

	.message-sent .message-sender {
			display: none;
	}

	.message-attachments img, .message-attachments video {
		max-width: 200px;
		display: block;
	}
</style>

{{range .}}
	<h1 class="group-title">{{.Group}}</h1>

	{{range .Messages}}
		<div class="message {{ if .Sent }}message-sent{{else}}message-received{{end}}">
			<h2 class="message-sender">
				{{.Sender}}
			</h2>
			<p class="message-body">
				{{.Content | formatContent }}
			</p>
			<div class="message-attachments">
				{{ range .Attachments }}
				  {{ if eq .Kind "img" }}
						<img src="/serve?path={{ .URL }}" />
				  {{ else if eq .Kind "png" }}
						<img src="/serve?path={{ .URL }}" />
				  {{ else if eq .Kind "gif" }}
						<img src="/serve?path={{ .URL }}" />
					{{ else if eq .Kind "video" }}
						<video controls src="/serve?path={{.URL}}" />
					{{ else if eq .Kind "audio" }}
						<audio controls src="/serve?path={{.URL}}" />
					{{ else if eq .Kind "vcard" }}
						<a href="/serve?path={{.URL}}">download vcard</a>
					{{ else }}
					  {{.URL}}
					{{ end }}
				{{ end }}
			</div>
			<small class="message-time">
				{{.Timestamp.Format "02/01/2006 15:04:05"}}
			</small>
		</div>
	{{ end }}

{{ end }}
`

	tmpl, err = template.New("msgs").Funcs(template.FuncMap{"hasPrefix": hasPrefix, "toURL": toURL, "formatContent": formatContent}).Parse(tmplStr)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", msgHandler)
	http.HandleFunc("/serve", serveHandler)

	log.Printf("Parsed %d messages", len(msgs))
	log.Printf("Listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))

}

func msgHandler(w http.ResponseWriter, r *http.Request) {

	err := tmpl.Execute(w, byGroup)
	if err != nil {
		panic(err)
	}

}

func serveHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Query().Get("path")

	if !strings.HasPrefix(path, "/Users/adrien/Dropbox/Applications/Messages/media") {
		panic("unauthorized")
	}

	if strings.HasSuffix(path, ".jpg") {
		w.Header().Add("Content-Type", "image/jpeg")
	} else if strings.HasSuffix(path, ".mp4") {
		w.Header().Add("Content-Type", "video/mp4")
	}

	file, err := os.Open(path)
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(500)
	}

	_, err = io.Copy(w, file)
	if err != nil {
		log.Printf("ERROR: %s", err)
		w.WriteHeader(500)
	}

}

func hasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func toURL(s string) template.URL {
	return template.URL(s)
}

func formatContent(s string) template.HTML {
	s = strings.Replace(s, "\n", "<br>", -1)
	s = strings.Replace(s, "\r", "<br>", -1)
	return template.HTML(s)
}
