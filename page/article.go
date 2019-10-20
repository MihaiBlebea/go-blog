package page

import (
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

type Article struct {
	*Page
	Title           string
	Category        string
	Created         time.Time
	Author          string
	Slug            string
	FeatureImage    string
	RelatedArticles []Article
}

func (ar *Article) Render(w io.Writer) error {
	err := ar.Template.Execute(w, ar)
	return err
}

func (ar *Article) CreatedAt() string {
	return ar.Created.Format("2006 Jan 02")
}

func (ar *Article) AddSlug(slug string) {
	ar.Slug = slug
}

func (ar *Article) AddRelated(article *Article) {
	ar.RelatedArticles = append(ar.RelatedArticles, *article)
}

func NewArticle(content []byte) *Article {

	meta, err := extractMeta(string(content))
	if err != nil {
		log.Panic(err)
	}

	title, err := getValueForKey(meta, "title")
	if err != nil {
		log.Panic(err)
	}

	created, err := getValueForKey(meta, "date")
	if err != nil {
		log.Panic(err)
	}

	category, err := getValueForKey(meta, "category")
	if err != nil {
		log.Panic(err)
	}

	author, err := getValueForKey(meta, "author")
	if err != nil {
		log.Panic(err)
	}

	image, err := getValueForKey(meta, "feature-image")
	if err != nil {
		log.Panic(err)
	}

	articleBody := removeMeta(string(content), meta)

	pg := New()
	pg.AddTitle(title)
	pg.AddMarkdownBody([]byte(articleBody))

	createdDate, err := time.Parse("2006-Jan-02", created)
	if err != nil {
		log.Panic(err)
	}

	return &Article{
		Page:         pg,
		Title:        title,
		Category:     category,
		Created:      createdDate,
		Author:       author,
		FeatureImage: image,
	}
}

func extractMeta(content string) (string, error) {
	re := regexp.MustCompile(`(?ms)---.*---`)
	found := re.FindString(string(content))
	if found == "" {
		return "", errors.New("No meta found in the content")
	}

	found = strings.Replace(found, "---", "", -1)
	found = strings.TrimSpace(found)

	return found, nil
}

func removeMeta(content, meta string) string {
	noMetaContent := strings.Replace(content, meta, "", -1)
	noMetaContent = strings.Replace(noMetaContent, "---", "", -1)
	return strings.TrimSpace(noMetaContent)
}

func getValueForKey(content, key string) (string, error) {
	re := regexp.MustCompile(key + ": .*")
	found := re.FindString(content)
	if found == "" {
		return "", errors.New("No key found")
	}

	found = strings.Replace(found, key+":", "", -1)
	found = strings.TrimSpace(found)

	return found, nil
}
