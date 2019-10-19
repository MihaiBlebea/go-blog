package page

import "io"

type Blog struct {
	*Page
	Articles []Article
}

func (b *Blog) Render(w io.Writer) error {
	err := b.Template.Execute(w, b)
	return err
}

func (b *Blog) AddArticle(article *Article) {
	b.Articles = append(b.Articles, *article)
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
