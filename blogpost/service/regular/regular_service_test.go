package regular

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/stretchr/testify/mock"

	"github.com/printezisn/serverless-blog-back/blogpost/model"

	repoMocks "github.com/printezisn/serverless-blog-back/blogpost/repository/mocks"
)

// TestNew tests that the New method creates the service properly.
func TestNew(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)

	if service.repo != repo {
		t.Error("The repository is not set correctly.")
	}
}

// TestCreateWithValidationErrors tests that the Create method returns errors when the input is invalid.
func TestCreateWithValidationErrors(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{}

	response := service.Create(post)

	if response.StatusCode != 400 {
		t.Errorf("The status code was expected to be 400, but it was %d.", response.StatusCode)
	}
	if len(response.Errors) == 0 {
		t.Error("The response was expected to contain errors, but it didn't.")
	}
}

// TestCreateWithNonConditionalError tests that the Create method returns the correct response when an unexpected
// error occurs.
func TestCreateWithNonConditionalError(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}

	repo.On("Create", mock.MatchedBy(matchedByPost(post))).Return(post, errors.New("unexpected error"))

	response := service.Create(post)

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestCreateWithConditionalErrorAndConflicts tests that the Create method returns the correct response when the
// blog post is already stored with different values.
func TestCreateWithConditionalErrorAndConflicts(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}
	storedPost := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 2}

	err := awserr.New("ConditionalCheckFailedException", "error", errors.New("error"))
	requestFailure := awserr.NewRequestFailure(err, 400, "1")

	repo.On("Create", mock.MatchedBy(matchedByPost(post))).Return(post, requestFailure)
	repo.On("Get", post.ID).Return(storedPost, true, nil)

	response := service.Create(post)

	if response.StatusCode != 409 {
		t.Errorf("The status code was expected to be 409, but it was %d.", response.StatusCode)
	}
}

// TestCreateWithConditionalErrorAndNoConflicts tests that the Create method returns the correct response when the
// blog post is already stored with the same values.
func TestCreateWithConditionalErrorAndNoConflicts(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}
	storedPost := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}

	err := awserr.New("ConditionalCheckFailedException", "error", errors.New("error"))
	requestFailure := awserr.NewRequestFailure(err, 400, "1")

	repo.On("Create", mock.MatchedBy(matchedByPost(post))).Return(post, requestFailure)
	repo.On("Get", post.ID).Return(storedPost, true, nil)

	response := service.Create(post)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if response.Entity != storedPost {
		t.Error("The entity was expected to be ", storedPost, " but it was ", response.Entity)
	}
}

// TestCreateWithConditionalErrorAndFailure tests that the Create method returns the correct response when the
// blog post is already stored and an unexpected error occurs while retrieving it.
func TestCreateWithConditionalErrorAndFailure(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}

	err := awserr.New("ConditionalCheckFailedException", "error", errors.New("error"))
	requestFailure := awserr.NewRequestFailure(err, 400, "1")

	repo.On("Create", mock.MatchedBy(matchedByPost(post))).Return(post, requestFailure)
	repo.On("Get", post.ID).Return(post, false, errors.New("unexpected error"))

	response := service.Create(post)

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestCreateWithSuccess tests that the Create method returns the correct response when the operation is successful.
func TestCreateWithSuccess(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}

	repo.On("Create", mock.MatchedBy(matchedByPost(post))).Return(post, nil)

	response := service.Create(post)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if response.Entity != post {
		t.Error("The entity was expected to be ", post, " but it was ", response.Entity)
	}
}

// TestUpdateWithValidationErrors tests that the Update method returns errors when the input is invalid.
func TestUpdateWithValidationErrors(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{}

	response := service.Update(post)

	if response.StatusCode != 400 {
		t.Errorf("The status code was expected to be 400, but it was %d.", response.StatusCode)
	}
	if len(response.Errors) == 0 {
		t.Error("The response was expected to contain errors, but it didn't.")
	}
}

// TestUpdateWithNonConditionalError tests that the Update method returns the correct response when an unexpected
// error occurs.
func TestUpdateWithNonConditionalError(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}
	postUpdate := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 2}

	repo.On("Update", int64(1), mock.MatchedBy(matchedByPost(postUpdate))).Return(postUpdate, errors.New("unexpected error"))

	response := service.Update(post)

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestUpdateWithConditionalErrorAndConflicts tests that the Update method returns the correct response when the
// blog post is already stored with a different revision and different values.
func TestUpdateWithConditionalErrorAndConflicts(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}
	postUpdate := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 2}
	storedPost := model.BlogPost{ID: "id", Title: "title", Description: "descr2", Body: "body", Revision: 2}

	err := awserr.New("ConditionalCheckFailedException", "error", errors.New("error"))
	requestFailure := awserr.NewRequestFailure(err, 400, "1")

	repo.On("Update", int64(1), mock.MatchedBy(matchedByPost(postUpdate))).Return(postUpdate, requestFailure)
	repo.On("Get", post.ID).Return(storedPost, true, nil)

	response := service.Update(post)

	if response.StatusCode != 409 {
		t.Errorf("The status code was expected to be 409, but it was %d.", response.StatusCode)
	}
}

// TestUpdateWithConditionalErrorAndNoConflicts tests that the Update method returns the correct response when the
// blog post is already stored with a different revision but same values.
func TestUpdateWithConditionalErrorAndNoConflicts(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}
	postUpdate := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 2}
	storedPost := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 2}

	err := awserr.New("ConditionalCheckFailedException", "error", errors.New("error"))
	requestFailure := awserr.NewRequestFailure(err, 400, "1")

	repo.On("Update", int64(1), mock.MatchedBy(matchedByPost(postUpdate))).Return(postUpdate, requestFailure)
	repo.On("Get", post.ID).Return(storedPost, true, nil)

	response := service.Update(post)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if response.Entity != storedPost {
		t.Error("The entity was expected to be ", storedPost, " but it was ", response.Entity)
	}
}

// TestUpdateWithConditionalErrorAndFailure tests that the Update method returns the correct response when the
// blog post is stored with a different revision and an unexpected error occurs while retrieving it.
func TestUpdateWithConditionalErrorAndFailure(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}
	postUpdate := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 2}

	err := awserr.New("ConditionalCheckFailedException", "error", errors.New("error"))
	requestFailure := awserr.NewRequestFailure(err, 400, "1")

	repo.On("Update", int64(1), mock.MatchedBy(matchedByPost(postUpdate))).Return(postUpdate, requestFailure)
	repo.On("Get", post.ID).Return(post, false, errors.New("unexpected error"))

	response := service.Update(post)

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestUpdateWithSuccess tests that the Update method returns the correct response when the operation is successful.
func TestUpdateWithSuccess(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}
	postUpdate := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 2}

	repo.On("Update", int64(1), mock.MatchedBy(matchedByPost(postUpdate))).Return(postUpdate, nil)

	response := service.Update(post)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if response.Entity != postUpdate {
		t.Error("The entity was expected to be ", post, " but it was ", response.Entity)
	}
}

// TestDeleteWithError tests that the Delete method returns the correct response when an unexpected error occurs.
func TestDeleteWithError(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	id := "id"

	repo.On("Delete", id).Return(false, errors.New("unexpected error"))

	response := service.Delete(id)

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestDeleteWithNotFound tests that the Delete method returns the correct response when the blog post is not found.
func TestDeleteWithNotFound(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	id := "id"

	repo.On("Delete", id).Return(false, nil)

	response := service.Delete(id)

	if response.StatusCode != 404 {
		t.Errorf("The status code was expected to be 404, but it was %d.", response.StatusCode)
	}
}

// TestDeleteWithFound tests that the Delete method returns the correct response when the blog post is found.
func TestDeleteWithFound(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	id := "id"

	repo.On("Delete", id).Return(true, nil)

	response := service.Delete(id)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
}

// TestGetWithError tests that the Get method returns the correct response when an unexpected error occurs.
func TestGetWithError(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	id := "id"

	repo.On("Get", id).Return(model.BlogPost{}, false, errors.New("unexpected error"))

	response := service.Get(id)

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestGetWithNotFound tests that the Get method returns the correct response when the blog post is not found.
func TestGetWithNotFound(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}

	repo.On("Get", post.ID).Return(post, false, nil)

	response := service.Get(post.ID)

	if response.StatusCode != 404 {
		t.Errorf("The status code was expected to be 404, but it was %d.", response.StatusCode)
	}
}

// TestGetWithFound tests that the Get method returns the correct response when the blog post is found.
func TestGetWithFound(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	post := model.BlogPost{ID: "id", Title: "title", Description: "descr", Body: "body", Revision: 1}

	repo.On("Get", post.ID).Return(post, true, nil)

	response := service.Get(post.ID)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if response.Entity != post {
		t.Error("The entity was expected to be ", post, " but it was ", response.Entity)
	}
}

// TestGetAllWithError tests that the GetAll method returns the correct response when there is an unexpected error.
func TestGetAllWithError(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)

	repo.On("GetAll", service.pageSize+1).Return([]model.BlogPost{}, errors.New("unexpected error"))

	response := service.GetAll()

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestGetAllWithFewItems tests that the GetAll method returns the correct response when there are only a few items.
func TestGetAllWithFewItems(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	posts := []model.BlogPost{
		model.BlogPost{ID: "id1", Title: "title1", Description: "descr1", Body: "body1", Revision: 1},
		model.BlogPost{ID: "id2", Title: "title2", Description: "descr2", Body: "body2", Revision: 1},
	}
	cursor := model.Cursor{}

	repo.On("GetAll", service.pageSize+1).Return(posts, nil)

	response := service.GetAll()

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if _, ok := response.Entity.(model.Page); !ok {
		t.Error("The entity was expected to be of type Page, but it was ", reflect.TypeOf(response.Entity))
	}

	page, _ := response.Entity.(model.Page)
	if !compareSlices(page.Posts, posts) {
		t.Error("The posts were expected to be ", posts, " but they were ", page.Posts)
	}
	if page.Cursor != cursor {
		t.Error("The cursor was expected to be ", cursor, " but it was ", page.Cursor)
	}
}

// TestGetAllWithMoreItems tests that the GetAll method returns the correct response when there are more items.
func TestGetAllWithMoreItems(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	service.pageSize = 1
	posts := []model.BlogPost{
		model.BlogPost{ID: "id1", Title: "title1", Description: "descr1", Body: "body1", Revision: 1, CreationTimestamp: 2},
		model.BlogPost{ID: "id2", Title: "title2", Description: "descr2", Body: "body2", Revision: 1, CreationTimestamp: 1},
	}
	cursor := model.Cursor{CreationTimestamp: 2, ID: "id1"}
	pageSlice := posts[:1]

	repo.On("GetAll", service.pageSize+1).Return(posts, nil)

	response := service.GetAll()

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if _, ok := response.Entity.(model.Page); !ok {
		t.Error("The entity was expected to be of type Page, but it was ", reflect.TypeOf(response.Entity))
	}

	page, _ := response.Entity.(model.Page)
	if !compareSlices(page.Posts, pageSlice) {
		t.Error("The posts were expected to be ", pageSlice, " but they were ", page.Posts)
	}
	if page.Cursor != cursor {
		t.Error("The cursor was expected to be ", cursor, " but it was ", page.Cursor)
	}
}

// TestGetMoreWithError tests that the GetMore method returns the correct response when there is an unexpected error.
func TestGetMoreWithError(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	lastID := "id"
	lastCreationTimestamp := int64(1)

	repo.On("GetMore", lastID, lastCreationTimestamp, service.pageSize+1).Return([]model.BlogPost{}, errors.New("unexpected error"))

	response := service.GetMore(lastID, lastCreationTimestamp)

	if response.StatusCode != 500 {
		t.Errorf("The status code was expected to be 500, but it was %d.", response.StatusCode)
	}
}

// TestGetMoreWithFewItems tests that the GetMore method returns the correct response when there are only a few items.
func TestGetMoreWithFewItems(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	posts := []model.BlogPost{
		model.BlogPost{ID: "id1", Title: "title1", Description: "descr1", Body: "body1", Revision: 1},
		model.BlogPost{ID: "id2", Title: "title2", Description: "descr2", Body: "body2", Revision: 1},
	}
	cursor := model.Cursor{}
	lastID := "id"
	lastCreationTimestamp := int64(1)

	repo.On("GetMore", lastID, lastCreationTimestamp, service.pageSize+1).Return(posts, nil)

	response := service.GetMore(lastID, lastCreationTimestamp)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if _, ok := response.Entity.(model.Page); !ok {
		t.Error("The entity was expected to be of type Page, but it was ", reflect.TypeOf(response.Entity))
	}

	page, _ := response.Entity.(model.Page)
	if !compareSlices(page.Posts, posts) {
		t.Error("The posts were expected to be ", posts, " but they were ", page.Posts)
	}
	if page.Cursor != cursor {
		t.Error("The cursor was expected to be ", cursor, " but it was ", page.Cursor)
	}
}

// TestGetMoreWithMoreItems tests that the GetMore method returns the correct response when there are more items.
func TestGetMoreWithMoreItems(t *testing.T) {
	repo := new(repoMocks.Repo)
	service := New(repo)
	service.pageSize = 1
	posts := []model.BlogPost{
		model.BlogPost{ID: "id1", Title: "title1", Description: "descr1", Body: "body1", Revision: 1, CreationTimestamp: 2},
		model.BlogPost{ID: "id2", Title: "title2", Description: "descr2", Body: "body2", Revision: 1, CreationTimestamp: 1},
	}
	cursor := model.Cursor{CreationTimestamp: 2, ID: "id1"}
	pageSlice := posts[:1]
	lastID := "id"
	lastCreationTimestamp := int64(1)

	repo.On("GetMore", lastID, lastCreationTimestamp, service.pageSize+1).Return(posts, nil)

	response := service.GetMore(lastID, lastCreationTimestamp)

	if response.StatusCode != 200 {
		t.Errorf("The status code was expected to be 200, but it was %d.", response.StatusCode)
	}
	if _, ok := response.Entity.(model.Page); !ok {
		t.Error("The entity was expected to be of type Page, but it was ", reflect.TypeOf(response.Entity))
	}

	page, _ := response.Entity.(model.Page)
	if !compareSlices(page.Posts, pageSlice) {
		t.Error("The posts were expected to be ", pageSlice, " but they were ", page.Posts)
	}
	if page.Cursor != cursor {
		t.Error("The cursor was expected to be ", cursor, " but it was ", page.Cursor)
	}
}

func matchedByPost(expectedPost model.BlogPost) func(model.BlogPost) bool {
	return func(actualPost model.BlogPost) bool {
		return actualPost.ID == expectedPost.ID && actualPost.Title == expectedPost.Title &&
			actualPost.Description == expectedPost.Description &&
			actualPost.Body == expectedPost.Body && actualPost.Revision == expectedPost.Revision
	}
}

func compareSlices(arr1, arr2 []model.BlogPost) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}
