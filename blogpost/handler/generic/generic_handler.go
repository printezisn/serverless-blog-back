package generic

import "github.com/aws/aws-lambda-go/events"

// Handler handles requests for blog posts
type Handler interface {
	Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}
