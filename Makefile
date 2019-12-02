MD5_COMMAND = $(shell openssl md5 serverless-blog-back | awk '{ print $$NF }')

all: load fmt test build package build_template

load:
	go get ./...

fmt:
	go fmt ./...

test:
	go test ./...

build:
	GOOS=linux go build

package:
	zip serverless-blog-back-${MD5_COMMAND}.zip serverless-blog-back

build_template:
	cat ./template.yml | sed 's/serverless-blog-back.zip/serverless-blog-back-${MD5_COMMAND}.zip/g' > deployment_template.yml

deploy: update_s3_bucket sam_deploy

update_s3_bucket: empty_s3_bucket copy_to_s3

empty_s3_bucket:
	aws s3 rm s3://${CODE_URI_BUCKET} --recursive

copy_to_s3:
	aws s3 cp ./serverless-blog-back*.zip s3://${CODE_URI_BUCKET}

sam_deploy:
	sam deploy --stack-name ${STACK_NAME} --template-file ./deployment_template.yml --capabilities CAPABILITY_IAM --parameter-overrides \
		CodeUriBucket=${CODE_URI_BUCKET}

clean:
	rm serverless-blog-back serverless-blog-back*.zip deployment_template.yml

run:
	sam local start-api

