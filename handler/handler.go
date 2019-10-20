package handler

import (
	"fmt"
	"go-blog/page"
	"time"

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
	startTime := time.Now()

	blg := page.NewBlog()

	articles, err := loadArticlesFromFolder("./content")
	if err != nil {
		log.Panic(err)
	}

	blg.AddArticles(articles)
	blg.AddTemplate(h.templates.Lookup("blog.gohtml"))
	blg.AddTitle("Blog")

	err = blg.Render(w)
	if err != nil {
		log.Panic(err)
	}

	completeTime := time.Now().Sub(startTime)
	fmt.Println(completeTime)
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

	articles, err := loadArticlesFromFolder("./content")
	if err != nil {
		log.Panic(err)
	}

	for _, relatedArticle := range articles[:3] {
		art.AddRelated(&relatedArticle)
	}

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
