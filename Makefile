# Makefile para Code Explainer com Docker

APP_NAME=code-explainer
DOCKER_USER=mvcbotelho
DOCKER_TAG=latest
IMAGE=$(DOCKER_USER)/$(APP_NAME):$(DOCKER_TAG)

.PHONY: run build tidy docker docker-run docker-push clean

run:
	go run main.go

build:
	go build -o $(APP_NAME)

tidy:
	go mod tidy

docker:
	docker build -t $(APP_NAME) .

docker-run:
	docker run --rm -it --env-file .env $(APP_NAME)

docker-tag:
	docker tag $(APP_NAME) $(IMAGE)

docker-push: docker docker-tag
	docker push $(IMAGE)

clean:
	rm -f $(APP_NAME)
