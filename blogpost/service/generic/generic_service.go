package generic

import (
	"github.com/printezisn/serverless-blog-back/blogpost/model"
	gloBalModel "github.com/printezisn/serverless-blog-back/global/model"
)

// Service represents the service layer for blog posts.
type Service interface {
	Create(post model.BlogPost) gloBalModel.Response
	Update(post model.BlogPost) gloBalModel.Response
	Delete(id string) gloBalModel.Response
	Get(id string) gloBalModel.Response
	GetAll() gloBalModel.Response
	GetMore(lastID string, lastCreationTimestamp int64) gloBalModel.Response
}
