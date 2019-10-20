package page

import (
	"io"
	"sort"
)

type Blog struct {
	*Page
	Articles []Article
}

// ByDate is a custom type to sort the articles by their created date
type ByDate []Article

func (d ByDate) Len() int {
	return len(d)
}

func (d ByDate) Less(i, j int) bool {
	return d[i].Created.After(d[j].Created)
}

func (d ByDate) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (b *Blog) Render(w io.Writer) error {
	err := b.Template.Execute(w, b)
	return err
}

func (b *Blog) AddArticle(article *Article) {
	b.Articles = append(b.Articles, *article)
	sort.Sort(ByDate(b.Articles))
}

func (b *Blog) AddArticles(articles []Article) {
	for _, article := range articles {
		b.AddArticle(&article)
	}
}

func NewBlog() *Blog {
	pg := New()
	var articles []Article
	return &Blog{
		pg,
		articles,
	}
}
