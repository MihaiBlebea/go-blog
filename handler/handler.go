package handler

import (
	"fmt"
	"go-blog/page"
	"time"

	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	templates    *template.Template
	articleCache *ArticleCache
}

func (h *Handler) GetHomepage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	home := page.NewHome()
	home.AddTemplate(h.templates.Lookup("home.gohtml"))

	err := home.Render(w)
	if err != nil {
		log.Panic(err)
	}
}

func (h *Handler) GetBlog(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	startTime := time.Now()

	blg := page.NewBlog()

	articles, err := h.articleCache.Articles()
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
	fmt.Println("/blog : ", completeTime)
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	startTime := time.Now()

	slug := ps.ByName("slug")

	// Get all the articles
	articles, err := h.articleCache.Articles()
	if err != nil {
		log.Panic(err)
	}

	var art page.Article
	found := false
	for _, article := range articles {
		if article.Slug == slug {
			art = article
			found = true
		}
	}

	if found == false {
		h.Get404(w, r, http.StatusNotFound)
		return
	}

	art.AddTemplate(h.templates.Lookup("article.gohtml"))

	for _, relatedArticle := range articles[:3] {
		art.AddRelated(&relatedArticle)
	}

	art.Render(w)
	if err != nil {
		log.Panic(err)
	}

	completeTime := time.Now().Sub(startTime)
	fmt.Println("/article : ", completeTime)
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
	articleCache := NewArticleCache("./content/articles")

	go func() {
		ticker := time.NewTicker(time.Minute * 3)

		for {
			tempCache := NewArticleCache("./content/articles")
			tempCache.Articles()

			articleCache.Mutex.Lock()
			articleCache.Cache = tempCache.Cache
			articleCache.CacheExpiration = tempCache.CacheExpiration
			articleCache.Mutex.Unlock()

			<-ticker.C
		}
	}()

	return &Handler{
		templates,
		articleCache,
	}
}
