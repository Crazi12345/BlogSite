package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Page struct {
	ID      int
	Title   string
	Content template.HTML
}
type Todo struct {
	ID    int
	Title string
	Done  bool
}

func mdToHTML(md []byte) template.HTML {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
     htmlString := string(markdown.Render(doc, renderer))

    // Replace "public/" with "static/" in the HTML
    htmlString = strings.ReplaceAll(htmlString, "public/", "static/")


    return template.HTML(htmlString)
}
func main() {

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("Blog/public"))))
	var pages []Page
	err := fs.WalkDir(os.DirFS("Blog/Posts"), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".md") {
			mds, err := fs.ReadFile(os.DirFS("Blog/Posts"), path)

			if err != nil {
				return err
			}
			htmlContent := mdToHTML(mds)

			// Create a new page
			page := Page{
				ID:      len(pages) + 1,
				Title:   strings.TrimSuffix(d.Name(), ".md"), // Extract title from filename
				Content: template.HTML(htmlContent),
			}
			pages = append(pages, page)

			// Create a unique route for each page
			http.HandleFunc("/"+page.Title, func(w http.ResponseWriter, r *http.Request) {
				tmpl, err := template.ParseFiles("test.html") // Use your template file
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				tmpl.Execute(w, page) // Pass the individual page to the template
			})

		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
