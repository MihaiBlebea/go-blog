package main

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"go-blog/handler"
)

type Page struct {
	Title string
	Body  template.HTML
}

var templ *template.Template

func init() {
	templ = template.Must(template.ParseGlob("./templates/*.gohtml"))

}

func main() {

	hdl := handler.New(templ)
	router := httprouter.New()

	router.GET("/", hdl.GetHomepage)
	router.GET("/blog", hdl.GetBlog)
	router.GET("/blog/:slug", hdl.GetArticle)
	router.GET("/contact", hdl.GetContact)

	router.GET("/assets/:file", hdl.GetAsset)

	// Set the custom error page
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdl.Get404(w, r, http.StatusNotFound)
		return
	})

	http.ListenAndServe(":8081", router)
}
