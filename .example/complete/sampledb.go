package main

import "fmt"

type Article struct {
	Author      string
	Description string
}

type Repository struct {
	articleMap map[int]Article
}

func New() Repository {
	return Repository{articleMap: make(map[int]Article)}
}

func (repo Repository) Create(article Article, id int) {
	repo.articleMap[id] = article
}
func (repo Repository) Delete(id int) {
	delete(repo.articleMap, id)
}
func (repo Repository) Get(id int) Article {
	return repo.articleMap[id]
}

func (article Article) String() string {
	return fmt.Sprintf("Author: %s \n Description: %s", article.Author, article.Description)
}
