package page

import (
	"html/template"
	"io"

	"github.com/gomarkdown/markdown"
)

type IPage interface {
	Render(w io.Writer) error
}

type Page struct {
	Template *template.Template
	Title    string
	Body     template.HTML
	Css      template.CSS
}

func (p *Page) AddTemplate(t *template.Template) {
	p.Template = t
}

func (p *Page) AddTitle(title string) {
	p.Title = title
}

func (p *Page) AddBody(body string) {
	p.Body = template.HTML(body)
}

func (p *Page) AddMarkdownBody(body []byte) {
	htmlContent := markdown.ToHTML(body, nil, nil)
	p.Body = template.HTML(htmlContent)
}

func (p *Page) AddCss(css string) {
	p.Css = template.CSS(css)
}

func (p *Page) Render(w io.Writer) error {
	err := p.Template.Execute(w, p)
	return err
}

func New() *Page {
	return &Page{}
}
