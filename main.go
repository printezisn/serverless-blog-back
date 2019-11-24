package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	regularHandler "github.com/printezisn/serverless-blog-back/blogpost/handler/regular"
	"github.com/printezisn/serverless-blog-back/blogpost/repository/dynamodb"
	regularService "github.com/printezisn/serverless-blog-back/blogpost/service/regular"
)

func main() {
	repo := dynamodb.New()
	service := regularService.New(&repo)
	handler := regularHandler.New(&service)

	lambda.Start(handler.Handle)
}
