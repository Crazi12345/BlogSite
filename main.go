package main

import (
	"html/template"
	"log"
	"net/http"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
    "io/fs"
	"os"
    "fmt"
)

type Page struct {
    ID int
    Content template.HTML
}
type Todo struct {
	ID    int
	Title string
	Done  bool
}

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
func main() {

    mds,_  := fs.ReadFile(os.DirFS("/home/tired_atlas/Documents/Main Brain"), "Laravel.md")
    md := []byte(mds)
	html := mdToHTML(md)

    fmt.Println(string(html))
    pages := Page{
        ID: 1,
       Content: template.HTML(string(html)),
    }
	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("test.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
        w.Header().Set("Content-Type", "text/html; charset=utf-8")

        // Directly write HTML output

		tmpl.Execute(w, pages)
	})
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
