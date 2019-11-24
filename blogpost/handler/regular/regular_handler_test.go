package regular

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/printezisn/serverless-blog-back/blogpost/service/mocks"

	"github.com/printezisn/serverless-blog-back/blogpost/model"

	globalModel "github.com/printezisn/serverless-blog-back/global/model"

	"github.com/aws/aws-lambda-go/events"
)

// TestHandleCreateWithInvalidInput tests that the POST "/posts" request returns the correct response when the input
// is invalid.
func TestHandleCreateWithInvalidInput(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	request := events.APIGatewayProxyRequest{Path: "/posts", HTTPMethod: "PUT", Body: "error"}
	response, _ := handler.Handle(request)

	if response.StatusCode != 400 {
		t.Errorf("The status code was expected to be 400, but it was %d.", response.StatusCode)
	}
}

// TestHandleCreateWithSuccess tests that the PUT "/posts" request returns the correct response when the operation is successful.
func TestHandleCreateWithSuccess(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	post := model.BlogPost{ID: "id"}
	postBytes, _ := json.Marshal(post)
	postJSON := string(postBytes)
	request := events.APIGatewayProxyRequest{Path: "/posts", HTTPMethod: "PUT", Body: postJSON}
	expectedResponse := globalModel.Response{Entity: "response", StatusCode: 200}
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	expectedResponseJSON := string(expectedResponseBytes)

	service.On("Create", post).Return(expectedResponse)

	actualResponse, _ := handler.Handle(request)

	if actualResponse.StatusCode != expectedResponse.StatusCode {
		t.Errorf("The status code was expected to be %d, but it was %d.", expectedResponse.StatusCode, actualResponse.StatusCode)
	}
	if actualResponse.Body != expectedResponseJSON {
		t.Error("The body was expected to be ", expectedResponseJSON, " but it was ", actualResponse.Body)
	}
}

// TestHandleUpdateWithInvalidInput tests that the POST "/posts" request returns the correct response when the input
// is invalid.
func TestHandleUpdateWithInvalidInput(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	request := events.APIGatewayProxyRequest{Path: "/posts", HTTPMethod: "POST", Body: "error"}
	response, _ := handler.Handle(request)

	if response.StatusCode != 400 {
		t.Errorf("The status code was expected to be 400, but it was %d.", response.StatusCode)
	}
}

// TestHandleUpdateWithSuccess tests that the POST "/posts" request returns the correct response when the operation is successful.
func TestHandleUpdateWithSuccess(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	post := model.BlogPost{ID: "id"}
	postBytes, _ := json.Marshal(post)
	postJSON := string(postBytes)
	request := events.APIGatewayProxyRequest{Path: "/posts", HTTPMethod: "POST", Body: postJSON}
	expectedResponse := globalModel.Response{Entity: "response", StatusCode: 200}
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	expectedResponseJSON := string(expectedResponseBytes)

	service.On("Update", post).Return(expectedResponse)

	actualResponse, _ := handler.Handle(request)

	if actualResponse.StatusCode != expectedResponse.StatusCode {
		t.Errorf("The status code was expected to be %d, but it was %d.", expectedResponse.StatusCode, actualResponse.StatusCode)
	}
	if actualResponse.Body != expectedResponseJSON {
		t.Error("The body was expected to be ", expectedResponseJSON, " but it was ", actualResponse.Body)
	}
}

// TestHandleDeleteWithSuccess tests that the DELETE "/posts" request returns the correct response when the operation is successful.
func TestHandleDeleteWithSuccess(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	request := events.APIGatewayProxyRequest{Path: "/posts/id", HTTPMethod: "DELETE", PathParameters: map[string]string{"id": "id"}}
	expectedResponse := globalModel.Response{Entity: "response", StatusCode: 200}
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	expectedResponseJSON := string(expectedResponseBytes)

	service.On("Delete", "id").Return(expectedResponse)

	actualResponse, _ := handler.Handle(request)

	if actualResponse.StatusCode != expectedResponse.StatusCode {
		t.Errorf("The status code was expected to be %d, but it was %d.", expectedResponse.StatusCode, actualResponse.StatusCode)
	}
	if actualResponse.Body != expectedResponseJSON {
		t.Error("The body was expected to be ", expectedResponseJSON, " but it was ", actualResponse.Body)
	}
}

// TestHandleGetWithSuccess tests that the GET "/posts/{id+}" request returns the correct response when the operation is successful.
func TestHandleGetWithSuccess(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	request := events.APIGatewayProxyRequest{Path: "/posts/id", HTTPMethod: "GET", PathParameters: map[string]string{"id": "id"}}
	expectedResponse := globalModel.Response{Entity: "response", StatusCode: 200}
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	expectedResponseJSON := string(expectedResponseBytes)

	service.On("Get", "id").Return(expectedResponse)

	actualResponse, _ := handler.Handle(request)

	if actualResponse.StatusCode != expectedResponse.StatusCode {
		t.Errorf("The status code was expected to be %d, but it was %d.", expectedResponse.StatusCode, actualResponse.StatusCode)
	}
	if actualResponse.Body != expectedResponseJSON {
		t.Error("The body was expected to be ", expectedResponseJSON, " but it was ", actualResponse.Body)
	}
}

// TestHandleGetAllWithSuccess tests that the GET "/posts" request returns the correct response when the operation is successful.
func TestHandleGetAllWithSuccess(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	request := events.APIGatewayProxyRequest{Path: "/posts", HTTPMethod: "GET"}
	expectedResponse := globalModel.Response{Entity: "response", StatusCode: 200}
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	expectedResponseJSON := string(expectedResponseBytes)

	service.On("GetAll").Return(expectedResponse)

	actualResponse, _ := handler.Handle(request)

	if actualResponse.StatusCode != expectedResponse.StatusCode {
		t.Errorf("The status code was expected to be %d, but it was %d.", expectedResponse.StatusCode, actualResponse.StatusCode)
	}
	if actualResponse.Body != expectedResponseJSON {
		t.Error("The body was expected to be ", expectedResponseJSON, " but it was ", actualResponse.Body)
	}
}

// TestHandleGetMoreWithSuccess tests that the GET "/posts?lastID=...&lastCreationTimestamp=..." request returns the correct
// response when the operation is successful.
func TestHandleGetMoreWithSuccess(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	lastID := "lastID"
	lastCreationTimestamp := int64(0)
	path := fmt.Sprintf("/posts?lastID=%s&lastCreationTimestamp=%d", lastID, lastCreationTimestamp)
	queryStringParameters := map[string]string{"lastID": lastID, "lastCreationTimestamp": strconv.FormatInt(lastCreationTimestamp, 10)}
	request := events.APIGatewayProxyRequest{Path: path, HTTPMethod: "GET", QueryStringParameters: queryStringParameters}
	expectedResponse := globalModel.Response{Entity: "response", StatusCode: 200}
	expectedResponseBytes, _ := json.Marshal(expectedResponse)
	expectedResponseJSON := string(expectedResponseBytes)

	service.On("GetMore", lastID, lastCreationTimestamp).Return(expectedResponse)

	actualResponse, _ := handler.Handle(request)

	if actualResponse.StatusCode != expectedResponse.StatusCode {
		t.Errorf("The status code was expected to be %d, but it was %d.", expectedResponse.StatusCode, actualResponse.StatusCode)
	}
	if actualResponse.Body != expectedResponseJSON {
		t.Error("The body was expected to be ", expectedResponseJSON, " but it was ", actualResponse.Body)
	}
}

// TestHandleGetMoreWithInvalidInput tests that the GET "/posts?lastID=...&lastCreationTimestamp=..." request returns the correct
// response when the input is invalid.
func TestHandleGetMoreWithInvalidInput(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	lastID := "lastID"
	lastCreationTimestamp := "invalid"
	path := fmt.Sprintf("/posts?lastID=%s&lastCreationTimestamp=%s", lastID, lastCreationTimestamp)
	queryStringParameters := map[string]string{"lastID": lastID, "lastCreationTimestamp": lastCreationTimestamp}
	request := events.APIGatewayProxyRequest{Path: path, HTTPMethod: "GET", QueryStringParameters: queryStringParameters}

	response, _ := handler.Handle(request)

	if response.StatusCode != 400 {
		t.Errorf("The status code was expected to be 400, but it was %d.", response.StatusCode)
	}
}

// TestInvalidRequest tests that the correct response is returned when the request is invalid.
func TestInvalidRequest(t *testing.T) {
	service := new(mocks.Service)
	handler := New(service)

	request := events.APIGatewayProxyRequest{Path: "/", HTTPMethod: "GET"}

	response, _ := handler.Handle(request)

	if response.StatusCode != 400 {
		t.Errorf("The status code was expected to be 400, but it was %d.", response.StatusCode)
	}
}
