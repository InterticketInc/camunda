CAMUNDA_VERSION=7.14.0

golang:
	@echo "--> Go Version"
	@go version

authors:
	@echo "--> Updating the AUTHORS"
	git log --format='%aN <%aE>' | sort -u > AUTHORS

test:
	go test -v ./...

changelog:
	git log $(shell git tag | tail -n1)..HEAD --no-merges --format=%B > changelog

docker-lint:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.27.0 golangci-lint run -v

update-openapi:
	curl -O https://app.camunda.com/nexus/repository/camunda-bpm/org/camunda/bpm/camunda-engine-rest-openapi/${CAMUNDA_VERSION}/camunda-engine-rest-openapi-${CAMUNDA_VERSION}.jar
	mv camunda-engine-rest-openapi-${CAMUNDA_VERSION}.jar camunda.zip
	mkdir -p camunda
	unzip camunda.zip openapi.json -d camunda
	rm camunda.zip
