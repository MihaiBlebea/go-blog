package handler

import (
	"errors"
	"fmt"
	"go-blog/page"
	"io/ioutil"
	"strings"
)

// loadArticlesFromFolder receives a path to the folder that holds the markdown articles and returns a collection of articles
func loadArticlesFromFolder(folderPath string) ([]page.Article, error) {
	files, err := ioutil.ReadDir("./content")
	if err != nil {
		return nil, errors.New("Could not load the articles")
	}

	var articles []page.Article
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

		articles = append(articles, *art)
	}

	return articles, nil
}
