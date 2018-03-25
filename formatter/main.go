package formatter

import (
	"html/template"
	"io"
	"sort"
	"strings"

	"github.com/adrienkohlbecker/messages/model"
)

const tmplStr = `
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
				  {{ else if eq .Kind "gifv" }}
				  		<video controls src="/serve?path={{.URL}}" autoplay loop />
					{{ else if eq .Kind "video" }}
						<video controls src="/serve?path={{.URL}}" />
					{{ else if eq .Kind "audio" }}
						<audio controls src="/serve?path={{.URL}}" />
					{{ else if eq .Kind "aac" }}
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

type Grouped struct {
	Group    string
	Messages model.Messages
}

func Format(msgs model.Messages, w io.Writer) error {

	tmpl, err := template.New("msgs").Funcs(template.FuncMap{"hasPrefix": hasPrefix, "toURL": toURL, "formatContent": formatContent}).Parse(tmplStr)
	if err != nil {
		return err
	}

	sort.Sort(msgs)

	var byGroup = make(map[string]*Grouped)

	for _, msg := range msgs {

		_, ok := byGroup[msg.Group]
		if !ok {
			byGroup[msg.Group] = &Grouped{Group: msg.Group, Messages: make(model.Messages, 0)}
		}

		byGroup[msg.Group].Messages = append(byGroup[msg.Group].Messages, msg)

	}

	err = tmpl.Execute(w, byGroup)
	if err != nil {
		return err
	}

	return nil

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
