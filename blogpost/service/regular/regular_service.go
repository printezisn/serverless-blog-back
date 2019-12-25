package regular

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/printezisn/serverless-blog-back/blogpost/model"
	postRepo "github.com/printezisn/serverless-blog-back/blogpost/repository/generic"
	gloBalModel "github.com/printezisn/serverless-blog-back/global/model"
)

// Service represents the regular service layer for blog posts.
type Service struct {
	repo     postRepo.Repo
	pageSize int64
}

// New creates a new instance of the regular service layer for blog posts.
func New(repo postRepo.Repo) Service {
	return Service{repo: repo, pageSize: 10}
}

// Create creates a new blog post.
func (service *Service) Create(post model.BlogPost) gloBalModel.Response {
	errs := post.Validate()
	if len(errs) > 0 {
		return gloBalModel.Response{Entity: post, Errors: errs, StatusCode: 400}
	}

	post.CreationTimestamp = time.Now().UTC().Unix()
	post.UpdateTimestamp = time.Now().UTC().Unix()

	newPost, err := service.repo.Create(post)

	if err != nil {
		log.Println("An error occurred while creating a new blog post: ", err)

		if requestFailure, ok := err.(awserr.RequestFailure); ok && requestFailure.Code() == "ConditionalCheckFailedException" {
			existingPost, found, err := service.repo.Get(newPost.ID)
			if err != nil {
				log.Println("An error occurred while fetching a blog post: ", err)
				return gloBalModel.Response{Entity: newPost, Errors: []string{}, StatusCode: 500}
			}

			newPost.CreationTimestamp = existingPost.CreationTimestamp
			newPost.UpdateTimestamp = existingPost.UpdateTimestamp

			if !found || existingPost != newPost {
				return gloBalModel.Response{Entity: existingPost, Errors: []string{}, StatusCode: 409}
			}

			return gloBalModel.Response{Entity: newPost, Errors: []string{}, StatusCode: 200}
		}

		return gloBalModel.Response{Entity: newPost, Errors: []string{}, StatusCode: 500}
	}

	return gloBalModel.Response{Entity: newPost, Errors: []string{}, StatusCode: 200}
}

// Update updates an existing blog post.
func (service *Service) Update(post model.BlogPost) gloBalModel.Response {
	errs := post.Validate()
	if len(errs) > 0 {
		return gloBalModel.Response{Entity: post, Errors: errs, StatusCode: 400}
	}

	post.UpdateTimestamp = time.Now().UTC().Unix()
	oldRevision := post.Revision
	post.Revision = oldRevision + 1

	updatedPost, err := service.repo.Update(oldRevision, post)

	if err != nil {
		log.Println("An error occurred while updating a blog post: ", err)

		if requestFailure, ok := err.(awserr.RequestFailure); ok && requestFailure.Code() == "ConditionalCheckFailedException" {
			existingPost, found, err := service.repo.Get(post.ID)
			if err != nil {
				log.Println("An error occurred while fetching a blog post: ", err)
				return gloBalModel.Response{Entity: post, Errors: []string{}, StatusCode: 500}
			}

			post.CreationTimestamp = existingPost.CreationTimestamp
			post.UpdateTimestamp = existingPost.UpdateTimestamp
			post.Revision = existingPost.Revision

			if !found || existingPost != post {
				return gloBalModel.Response{Entity: existingPost, Errors: []string{}, StatusCode: 409}
			}

			return gloBalModel.Response{Entity: existingPost, Errors: []string{}, StatusCode: 200}
		}

		return gloBalModel.Response{Entity: post, Errors: []string{}, StatusCode: 500}
	}

	return gloBalModel.Response{Entity: updatedPost, Errors: []string{}, StatusCode: 200}
}

// Delete deletes a blog post.
func (service *Service) Delete(id string) gloBalModel.Response {
	found, err := service.repo.Delete(id)

	if err != nil {
		log.Println("An error occurred while deleting a blog post: ", err)
		return gloBalModel.Response{Entity: id, Errors: []string{}, StatusCode: 500}
	}
	if !found {
		return gloBalModel.Response{Entity: id, Errors: []string{}, StatusCode: 404}
	}

	return gloBalModel.Response{Entity: id, Errors: []string{}, StatusCode: 200}
}

// Get fetches a blog post.
func (service *Service) Get(id string) gloBalModel.Response {
	post, found, err := service.repo.Get(id)

	if err != nil {
		log.Println("An error occurred while fetching a blog post: ", err)
		return gloBalModel.Response{Entity: post, Errors: []string{}, StatusCode: 500}
	}
	if !found {
		return gloBalModel.Response{Entity: post, Errors: []string{}, StatusCode: 404}
	}

	return gloBalModel.Response{Entity: post, Errors: []string{}, StatusCode: 200}
}

// GetAll fetches all blog posts (limit is 10).
func (service *Service) GetAll() gloBalModel.Response {
	posts, err := service.repo.GetAll(service.pageSize + 1)

	if err != nil {
		log.Println("An error occurred while fetching all blog posts: ", err)
		return gloBalModel.Response{Entity: posts, Errors: []string{}, StatusCode: 500}
	}

	hasMore := int64(len(posts)) > service.pageSize
	page := model.Page{Posts: posts}

	if hasMore {
		page.Posts = page.Posts[:service.pageSize]
		lastItem := page.Posts[service.pageSize-1]
		page.Cursor = lastItem.ID
	}

	return gloBalModel.Response{Entity: page, Errors: []string{}, StatusCode: 200}
}

// GetMore fetches more blog posts (limit is 10).
func (service *Service) GetMore(lastID string) gloBalModel.Response {
	posts, err := service.repo.GetMore(lastID, service.pageSize+1)

	if err != nil {
		log.Println("An error occurred while fetching more blog posts: ", err)
		return gloBalModel.Response{Entity: posts, Errors: []string{}, StatusCode: 500}
	}

	hasMore := int64(len(posts)) > service.pageSize
	page := model.Page{Posts: posts}

	if hasMore {
		page.Posts = page.Posts[:service.pageSize]
		lastItem := page.Posts[service.pageSize-1]
		page.Cursor = lastItem.ID
	}

	return gloBalModel.Response{Entity: page, Errors: []string{}, StatusCode: 200}
}
