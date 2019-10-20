package handler

import (
	"errors"
	"go-blog/page"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type ArticleCache struct {
	FolderPath      string
	Cache           []page.Article
	CacheExpiration time.Time
	Duration        time.Duration
	Mutex           sync.Mutex
}

func NewArticleCache(folderPath string) *ArticleCache {
	return &ArticleCache{
		FolderPath: folderPath,
		Duration:   5 * time.Minute,
	}
}

func (ac *ArticleCache) Articles() ([]page.Article, error) {
	ac.Mutex.Lock()
	defer ac.Mutex.Unlock()

	if time.Now().Sub(ac.CacheExpiration) < 0 {
		return ac.Cache, nil
	}

	if ac.FolderPath == "" {
		return nil, errors.New("No folder path was supplied")
	}

	articles, err := loadArticlesFromFolder(ac.FolderPath)
	if err != nil {
		return nil, err
	}
	ac.Cache = articles
	ac.CacheExpiration = time.Now().Add(5 * time.Minute)

	return articles, nil
}

// loadArticlesFromFolder receives a path to the folder that holds the markdown articles
// and returns a collection of articles
func loadArticlesFromFolder(folderPath string) ([]page.Article, error) {
	files, err := ioutil.ReadDir("./content/articles")
	if err != nil {
		return nil, errors.New("Could not load the articles")
	}

	var articles []page.Article

	// Create a struct to hold the channel results
	type result struct {
		article page.Article
		err     error
	}

	// Create the channel
	resultChanel := make(chan result)

	for _, file := range files {
		// Create the go routine that has a self invoking function
		go func(fileName string) {
			content, err := ioutil.ReadFile("./content/articles/" + fileName)
			if err != nil {
				resultChanel <- result{
					err: err,
				}
			}

			slug := strings.Replace(fileName, "_", "-", -1)
			slug = strings.Replace(slug, ".md", "", -1)

			art := page.NewArticle(content)
			art.AddSlug(slug)

			resultChanel <- result{
				article: *art,
			}
		}(file.Name())
	}

	// Get all the results from the channel without blocking the go routines
	for i := 0; i < len(files); i++ {
		result := <-resultChanel
		if result.err != nil {
			continue
		}

		articles = append(articles, result.article)
	}

	return articles, nil
}
