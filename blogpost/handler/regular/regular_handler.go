package regular

import (
	"encoding/json"
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
			if request.QueryStringParameters["lastID"] != "" {
				return getMoreBlogPosts(handle.service, request)
			}

			return getAllBlogPosts(handle.service, request)
		}
		if strings.ToLower(request.HTTPMethod) == "options" {
			return events.APIGatewayProxyResponse{
					Body: "Success",
					Headers: map[string]string{
						"Content-Type":                 "application/text",
						"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
						"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
						"Access-Control-Allow-Origin":  "*",
					},
					StatusCode: 200},
				nil
		}
	}

	return events.APIGatewayProxyResponse{
			Body: "The request is not supported",
			Headers: map[string]string{
				"Content-Type":                 "application/text",
				"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Origin":  "*",
			},
			StatusCode: 400},
		nil
}

func createBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var post model.BlogPost
	err := json.Unmarshal([]byte(request.Body), &post)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body: "The input model is not valid.",
				Headers: map[string]string{
					"Content-Type":                 "application/text",
					"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
					"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
					"Access-Control-Allow-Origin":  "*",
				},
				StatusCode: 400,
			},
			nil
	}

	response := service.Create(post)
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body: string(responseBytes),
			Headers: map[string]string{
				"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Origin":  "*",
			},
			StatusCode: response.StatusCode,
		},
		nil
}

func updateBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var post model.BlogPost
	err := json.Unmarshal([]byte(request.Body), &post)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body: "The input model is not valid.",
				Headers: map[string]string{
					"Content-Type":                 "application/text",
					"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
					"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
					"Access-Control-Allow-Origin":  "*",
				},
				StatusCode: 400,
			},
			nil
	}

	response := service.Update(post)
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body: string(responseBytes),
			Headers: map[string]string{
				"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Origin":  "*",
			},
			StatusCode: response.StatusCode,
		},
		nil
}

func deleteBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := service.Delete(request.PathParameters["id"])
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body: string(responseBytes),
			Headers: map[string]string{
				"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Origin":  "*",
			},
			StatusCode: response.StatusCode,
		},
		nil
}

func getBlogPost(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := service.Get(request.PathParameters["id"])
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body: string(responseBytes),
			Headers: map[string]string{
				"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Origin":  "*",
			},
			StatusCode: response.StatusCode,
		},
		nil
}

func getAllBlogPosts(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	response := service.GetAll()
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body: string(responseBytes),
			Headers: map[string]string{
				"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Origin":  "*",
			},
			StatusCode: response.StatusCode,
		},
		nil
}

func getMoreBlogPosts(service generic.Service, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lastID := request.QueryStringParameters["lastID"]

	if strings.TrimSpace(lastID) == "" {
		return events.APIGatewayProxyResponse{
				Body: "The input is invalid.",
				Headers: map[string]string{
					"Content-Type":                 "application/text",
					"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
					"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
					"Access-Control-Allow-Origin":  "*",
				},
				StatusCode: 400,
			},
			nil
	}

	response := service.GetMore(lastID)
	responseBytes, _ := json.Marshal(response)

	return events.APIGatewayProxyResponse{
			Body: string(responseBytes),
			Headers: map[string]string{
				"Access-Control-Allow-Methods": "DELETE,GET,HEAD,OPTIONS,PATCH,POST,PUT",
				"Access-Control-Allow-Headers": "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token",
				"Access-Control-Allow-Origin":  "*",
			},
			StatusCode: response.StatusCode,
		},
		nil
}
