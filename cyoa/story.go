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

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("").Parse(defaultHandlerTmpl))
}

type handler struct {
	s Story
	t *template.Template
}

// functional options
type handlerOption func(*handler)

// custom template
func WithTemplate(t *template.Template) handlerOption {
	return func(h *handler) {
		h.t = t
	}
}

func NewHandler(s Story, handlerOptions ...handlerOption) http.Handler {
	// default config
	h := handler{s, tmpl}

	// custom config
	for _, ho := range handlerOptions {
		ho(&h)
	}

	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "" || path == "/" {
		path = "/intro"
	}
	path = path[1:]

	if chapter, ok := h.s[path]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "The Story Not Found", http.StatusNotFound)
}
