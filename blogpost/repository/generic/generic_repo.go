package generic

import "github.com/printezisn/serverless-blog-back/blogpost/model"

// Repo represents the repository layer for blog posts.
type Repo interface {
	Create(post model.BlogPost) (model.BlogPost, error)
	Update(revision int64, post model.BlogPost) (model.BlogPost, error)
	Get(id string) (model.BlogPost, bool, error)
	Delete(id string) (bool, error)
	GetAll(pageSize int64) ([]model.BlogPost, error)
	GetMore(lastID string, pageSize int64) ([]model.BlogPost, error)
}
