package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
)

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Chapter struct {
	Title   string   `json:"title"`
	Storys  []string `json:"story"`
	Options []Option `json:"options"`
}

type Story map[string]Chapter

func JSONStory(r io.Reader) (Story, error) {
	var story Story
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&story); err != nil {
		return nil, err
	}
	return story, nil
}

const defaultHandlerTmpl = `
<!DOCTYPE html>
<html lang="zh-CN">

<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Choose-Your-Own-Adventure</title>
</head>

<body>
  <h1>{{.Title}}</h1>
	{{range .Storys}}
		<p>{{.}}</p>
	{{end}}
  <ul>
		{{range .Options}}
    <li>
			<a href="/{{.Arc}}">{{.Text}}</a>
    </li>
		{{end}}
  </ul>
</body>

</html>`

type handler struct {
	s Story
}

func NewHandler(s Story) http.Handler {
	return handler{s}
}

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := tmpl.Execute(w, chapter)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "The Story Not Found", http.StatusNotFound)
}
