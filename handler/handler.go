package handler

import (
	"fmt"
	"go-blog/page"

	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	templates *template.Template
}

func (h *Handler) GetHomepage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pg := page.New()

	pg.AddTemplate(h.templates.Lookup("home.gohtml"))
	pg.AddBody("<p>Hello from the homepage</p>")
	pg.AddTitle("Homepage")

	err := pg.Render(w)
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) GetBlog(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	blg := page.NewBlog()

	files, err := ioutil.ReadDir("./content")

	for _, file := range files {
		content, err := ioutil.ReadFile("./content/" + file.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}

		slug := strings.Replace(file.Name(), "_", "-", -1)
		slug = strings.Replace(slug, ".md", "", -1)

		art := page.NewArticle(content)
		art.AddSlug(slug)

		blg.AddArticle(art)
	}

	blg.AddTemplate(h.templates.Lookup("blog.gohtml"))
	blg.AddTitle("Blog")

	err = blg.Render(w)
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("slug")

	fileName := strings.Replace(slug, "-", "_", -1)
	content, err := ioutil.ReadFile("./content/" + fileName + ".md")
	if err != nil {
		h.Get404(w, r, http.StatusNotFound)
		return
	}

	art := page.NewArticle(content)
	art.AddTemplate(h.templates.Lookup("article.gohtml"))
	art.AddSlug(slug)

	art.Render(w)
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) GetContact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	pg := page.New()

	pg.AddTemplate(h.templates.Lookup("contact.gohtml"))
	pg.AddTitle("Contact")

	err := pg.Render(w)
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) Get404(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		pg := page.New()

		pg.AddTemplate(h.templates.Lookup("404.gohtml"))
		pg.AddTitle("404 | Please try again")

		err := pg.Render(w)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	file := ps.ByName("file")

	content, err := ioutil.ReadFile("./assets/" + file)
	if err != nil {
		log.Panic(err)
	}
	w.Write(content)
}

func New(templates *template.Template) *Handler {
	return &Handler{templates}
}
