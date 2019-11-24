package regular

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/printezisn/serverless-blog-back/blogpost/model"
	"github.com/printezisn/serverless-blog-back/blogpost/service/generic"
)

// Handler handles requests for blog posts.
type Handler struct {
	service generic.Service
}

// New creates and returns a new handler instance.
func New(service generic.Service) Handler {
	return Handler{service: service}
}

// Handle handles requests from the API Gateway.
func (handle *Handler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	path := strings.ToLower(request.Path)
	if strings.Index(path, "/posts") == 0 {
		if strings.ToLower(request.HTTPMethod) == "put" {
			return createBlogPost(handle.service, request)
		}
		if strings.ToLower(request.HTTPMethod) == "post" {
			return updateBlogPost(handle.service, request)
		}
		if strings.ToLower(request.HTTPMethod) == "delete" {
			return deleteBlogPost(handle.service, request)
		}
		if strings.ToLower(request.HTTPMethod) == "get" {
			if request.PathParameters["id"] != "" {
				return getBlogPost(handle.service, request)
			}
			if request.QueryStringParameters["lastID"] != "" && request.QueryStringParameters["lastCreationTimestamp"] != "" {
				return getMoreBlogPosts(handle.service, request)
			}

			return getAllBlogPosts(handle.service, request)
		}
	}

	return events.APIGatewayProxyResponse{
			Body: "The request is not supported",
			Headers: map[string]string{
				"Content-Type": "application/text",
			},
			StatusCode: 400},
		nil
}

func createBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var post model.BlogPost
	err := json.Unmarshal([]byte(request.Body), &post)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       "The input model is not valid.",
				StatusCode: 400,
			},
			nil
	}

	response := service.Create(post)
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body:       string(responseBytes),
			StatusCode: response.StatusCode,
		},
		nil
}

func updateBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var post model.BlogPost
	err := json.Unmarshal([]byte(request.Body), &post)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       "The input model is not valid.",
				StatusCode: 400,
			},
			nil
	}

	response := service.Update(post)
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body:       string(responseBytes),
			StatusCode: response.StatusCode,
		},
		nil
}

func deleteBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := service.Delete(request.PathParameters["id"])
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body:       string(responseBytes),
			StatusCode: response.StatusCode,
		},
		nil
}

func getBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := service.Get(request.PathParameters["id"])
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body:       string(responseBytes),
			StatusCode: response.StatusCode,
		},
		nil
}

func getAllBlogPosts(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := service.GetAll()
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body:       string(responseBytes),
			StatusCode: response.StatusCode,
		},
		nil
}

func getMoreBlogPosts(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lastID := request.QueryStringParameters["lastID"]
	lastCreationTimestamp, err := strconv.ParseInt(request.QueryStringParameters["lastCreationTimestamp"], 10, 64)

	if strings.TrimSpace(lastID) == "" || err != nil {
		return events.APIGatewayProxyResponse{
				Body: "The input is invalid.",
				Headers: map[string]string{
					"Content-Type": "application/text",
				},
				StatusCode: 400,
			},
			nil
	}

	response := service.GetMore(lastID, lastCreationTimestamp)
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body:       string(responseBytes),
			StatusCode: response.StatusCode,
		},
		nil
}
