package repository 

import "gitlab.com/pragmaticreviews/golan-mux-api/entity"

type PostResposiotry interface{
	Save(post *entity.Post)(*entity.Post, error)
	FindAll() ([]entity.Post, error)
	FindByID(id string)(*entity.Post, error)
	Delete(post *entity.Post) error
}