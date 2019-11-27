all: load fmt test build package

load:
	go get ./...

fmt:
	go fmt ./...

test:
	go test ./...

build:
	GOOS=linux go build

package:
	zip serverless-blog-back.zip serverless-blog-back

deploy: copy_to_s3 sam_deploy
	
copy_to_s3:
	aws s3 cp ./serverless-blog-back.zip s3://${CODE_URI_BUCKET}

sam_deploy:
	sam deploy --stack-name ${STACK_NAME} --template-file ./template.yml --capabilities CAPABILITY_IAM --parameter-overrides \
		CodeUriBucket=${CODE_URI_BUCKET} CodeUriKey=serverless-blog-back.zip

clean:
	rm serverless-blog-back serverless-blog-back.zip

run:
	sam local start-api

