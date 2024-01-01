REGION ?= asia-northeast1

export PROJECT_ID ?= $(shell gcloud config get-value project)

setup:
	gcloud config set project $(PROJECT_ID)
	gcloud config set run/region $(REGION)

build:
	docker build -t otel-metrics-and-traces/app ./app
	docker build -t otel-metrics-and-traces/collector ./collector

auth:
	gcloud auth configure-docker $(REGION)-docker.pkg.dev

push: auth
	docker tag otel-metrics-and-traces/app $(REGION)-docker.pkg.dev/$(PROJECT_ID)/run-otel/app
	docker tag otel-metrics-and-traces/collector $(REGION)-docker.pkg.dev/$(PROJECT_ID)/run-otel/collector

gcloud-build-and-push:
	gcloud builds submit app --tag $(REGION)-docker.pkg.dev/$(PROJECT_ID)/run-otel/app:latest
	gcloud builds submit collector --tag $(REGION)-docker.pkg.dev/$(PROJECT_ID)/run-otel/collector:latest

deploy:
	@envsubst '$$PROJECT_ID' < service.yaml.template > service.yaml
	gcloud run services replace service.yaml --region $(REGION)