package dynamodb

import (
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/printezisn/serverless-blog-back/blogpost/model"
)

// Repo represents a repository for blog posts that uses DynamoDB.
type Repo struct {
	tableName string
	indexName string
	client    *dynamodb.DynamoDB
}

// New returns a new repository instance for blog posts that uses DynamoDB.
func New() Repo {
	tableName, ok := os.LookupEnv("DYNAMODB_TABLE_NAME")
	if !ok {
		tableName = "posts"
	}

	return Repo{tableName: tableName, indexName: "creationTimestamp-id-index", client: nil}
}

// createClient creates a new DynamoDB client.
func (repo *Repo) createClient() {
	if repo.client == nil {
		session := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
		repo.client = dynamodb.New(session)
	}
}

// Create creates a new blog post in the database.
func (repo *Repo) Create(post model.BlogPost) (model.BlogPost, error) {
	repo.createClient()

	item, _ := dynamodbattribute.MarshalMap(post)
	input := &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String(repo.tableName),
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}

	_, err := repo.client.PutItem(input)

	return post, err
}

// Update updates an existing blog post in the database.
func (repo *Repo) Update(revision int64, post model.BlogPost) (model.BlogPost, error) {
	repo.createClient()

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":title": {
				S: aws.String(post.Title),
			},
			":description": {
				S: aws.String(post.Description),
			},
			":body": {
				S: aws.String(post.Body),
			},
			":updateTimestamp": {
				N: aws.String(strconv.FormatInt(post.UpdateTimestamp, 10)),
			},
			":oldRevision": {
				N: aws.String(strconv.FormatInt(revision, 10)),
			},
			":newRevision": {
				N: aws.String(strconv.FormatInt(post.Revision, 10)),
			},
		},
		TableName: aws.String(repo.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(post.ID),
			},
		},
		ReturnValues:        aws.String("ALL_NEW"),
		ConditionExpression: aws.String("revision = :oldRevision"),
		UpdateExpression: aws.String("set title = :title, description = :description, body = :body, " +
			"updateTimestamp = :updateTimestamp, revision = :newRevision"),
	}

	response, err := repo.client.UpdateItem(input)
	if err != nil {
		return model.BlogPost{}, err
	}

	var updatedPost model.BlogPost
	err = dynamodbattribute.ConvertFromMap(response.Attributes, &updatedPost)

	return updatedPost, err
}

// Get searches and returns a blog post based on its id.
func (repo *Repo) Get(id string) (model.BlogPost, bool, error) {
	repo.createClient()

	queryInput := &dynamodb.QueryInput{
		TableName: aws.String(repo.tableName),
		KeyConditions: map[string]*dynamodb.Condition{
			"id": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(id),
					},
				},
			},
		},
	}

	response, err := repo.client.Query(queryInput)
	if err != nil {
		return model.BlogPost{}, false, err
	}

	var posts []model.BlogPost
	if err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &posts); err != nil {
		return model.BlogPost{}, false, err
	}

	if len(posts) == 0 {
		return model.BlogPost{}, false, err
	}

	return posts[0], true, nil
}

// GetAll loads all blog posts from the database.
func (repo *Repo) GetAll(pageSize int64) ([]model.BlogPost, error) {
	repo.createClient()

	scanInput := &dynamodb.ScanInput{
		AttributesToGet: []*string{aws.String("id"), aws.String("title"), aws.String("description"),
			aws.String("revision"), aws.String("creationTimestamp"), aws.String("updateTimestamp")},
		TableName: aws.String(repo.tableName),
		IndexName: aws.String(repo.indexName),
		Limit:     aws.Int64(pageSize),
	}

	response, err := repo.client.Scan(scanInput)
	if err != nil {
		return []model.BlogPost{}, err
	}

	var posts []model.BlogPost
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &posts)
	if posts == nil {
		posts = []model.BlogPost{}
	}

	return posts, err
}

// GetMore loads more blog posts from the database.
func (repo *Repo) GetMore(lastID string, lastCreationTimestamp int64, pageSize int64) ([]model.BlogPost, error) {
	repo.createClient()

	scanInput := &dynamodb.ScanInput{
		AttributesToGet: []*string{aws.String("id"), aws.String("title"), aws.String("description"),
			aws.String("revision"), aws.String("creationTimestamp"), aws.String("updateTimestamp")},
		ExclusiveStartKey: map[string]*dynamodb.AttributeValue{
			"creationTimestamp": {
				N: aws.String(strconv.FormatInt(lastCreationTimestamp, 10)),
			},
			"id": {
				S: aws.String(lastID),
			},
		},
		TableName: aws.String(repo.tableName),
		IndexName: aws.String(repo.indexName),
		Limit:     aws.Int64(pageSize),
	}

	response, err := repo.client.Scan(scanInput)
	if err != nil {
		return []model.BlogPost{}, err
	}

	var posts []model.BlogPost
	err = dynamodbattribute.UnmarshalListOfMaps(response.Items, &posts)
	if posts == nil {
		posts = []model.BlogPost{}
	}

	return posts, err
}

// Delete deletes a blog post from the database.
func (repo *Repo) Delete(id string) (bool, error) {
	repo.createClient()

	deleteInput := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
		ReturnValues: aws.String("ALL_OLD"),
		TableName:    aws.String(repo.tableName),
	}

	response, err := repo.client.DeleteItem(deleteInput)
	if err != nil {
		return false, err
	}

	return len(response.Attributes) > 0, err
}
