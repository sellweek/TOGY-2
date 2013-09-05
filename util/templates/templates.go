package templates

import (
	"agility/util"
	"appengine/user"
	"html/template"
	"net/http"
	"os"
	"strings"
)

var (
	t         = "static" + string(os.PathSeparator) + "templates"
	templates *template.Template
)

//init injects few utility functions into templates we're using
func init() {
	templates = template.New("").Funcs(template.FuncMap{
		"equal": func(x, y int) bool {
			return x == y
		},
		"subtract": func(x, y int) int {
			return x - y
		},
		"add": func(x, y int) int {
			return x + y
		},
		"div": func(x, y int) float64 {
			return float64(x) / float64(y)
		},
		"mul": func(x, y int) int {
			return x * y
		}})
	// List of template files. When creating new template, add it here.
	templates = template.Must(parseFiles(templates, t))
	templates = template.Must(parseFiles(templates, t+string(os.PathSeparator)+"layout"))
}

//ParseFiles goes trough a folder (non-recursively), parsing and
//adding all HTML files into a template.
func parseFiles(t *template.Template, dir string) (temp *template.Template, err error) {
	f, err := os.Open(dir)
	if err != nil {
		return
	}

	fis, err := f.Readdir(0)
	if err != nil {
		return
	}

	filenames := make([]string, 0)

	for _, fi := range fis {
		if fi.IsDir() || getFileType(fi.Name()) != "html" {
			continue
		}

		filenames = append(filenames, dir+string(os.PathSeparator)+fi.Name())
	}

	temp, err = t.ParseFiles(filenames...)
	return
}

func getFileType(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}

	return parts[len(parts)-1]
}

//RenderLayout inserts template with given name into the layout and sets the title and pipeline.
//The template should be loaded inside templates variable
//If any arguments are provided after the context, they will be treated like links
//to JavaScript scripts to load in the header of the template.
func RenderLayout(c util.Context, tmpl string, title string, data interface{}, jsIncludes ...string) {
	RenderTemplate(c, "header.html", struct {
		Title      string
		JsIncludes []string
		Admin      bool
	}{title, jsIncludes, user.IsAdmin(c.Ac)})
	RenderTemplate(c, tmpl, data)
	RenderTemplate(c, "footer.html", nil)
}

//renderTemplate renders a single template
func RenderTemplate(c util.Context, tmpl string, data interface{}) {
	if err := templates.ExecuteTemplate(c.W, tmpl, data); err != nil {
		c.Ac.Errorf("Error 500. %v", err)
		http.Error(c.W, err.Error(), http.StatusInternalServerError)
	}
}
