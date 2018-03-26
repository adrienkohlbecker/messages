package formatter

import (
	"html/template"
	"io"
	"sort"
	"strings"

	"github.com/adrienkohlbecker/messages/model"
)

const tmplStr = `
<html>
<head>
<head>
<meta charset="UTF-8">
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
</head>

<body>
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
						<img src="{{ .URL | toURL  }}" />
				  {{ else if eq .Kind "sticker" }}
						<img src="{{ .URL | toURL  }}" />
				  {{ else if eq .Kind "gifv" }}
				  		<video controls src="{{.URL | toURL }}" autoplay loop />
					{{ else if eq .Kind "video" }}
						<video controls src="{{.URL | toURL }}" />
					{{ else if eq .Kind "audio" }}
						<audio controls src="{{.URL | toURL }}" />
					{{ else if eq .Kind "document" }}
						<a href="{{.URL | toURL }}">download document</a>
					{{ else }}
					  {{.URL | toURL }}
					{{ end }}
				{{ end }}
			</div>
			<small class="message-time">
				{{.Timestamp.Format "02/01/2006 15:04:05"}}
			</small>
		</div>
	{{ end }}

{{ end }}
</body>
</html>
`

func Format(msgs model.Messages, w io.Writer) error {

	tmpl, err := template.New("msgs").Funcs(template.FuncMap{"hasPrefix": hasPrefix, "toURL": toURL, "formatContent": formatContent}).Parse(tmplStr)
	if err != nil {
		return err
	}

	sort.Sort(msgs)

	err = tmpl.Execute(w, msgs.ByGroup())
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
