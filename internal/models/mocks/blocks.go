package mocks

import (
	"blog_project/internal/models"
	"time"
)

var mockBlock = models.Block{
	ID:          1,
	Created:     time.Now(),
	ImageURL1:   "url1",
	ImageURL2:   "url2",
	ImageURL3:   "url3",
	Description: "sometext",
	Title:       "Title",
}

type BlockModel struct {
}

func (m *BlockModel) GetAll() ([]models.Block, error) {
	return []models.Block{mockBlock}, nil
}

func (m *BlockModel) Insert(title, description, imageURL1, imageURL2, imageURL3 string) error {
	return nil
}

func (m *BlockModel) Update(id int, title, description, imageURL1, imageURL2, imageURL3 string) error {
	if id == 1 {
		if title != "Title" {
			return models.ErrInvalidCredentials
		}
		return nil
	}
	return models.ErrNoRecord
}
