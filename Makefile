# Constants

PROJECT_NAME=external-api
USER=fidesy-pay

APP_NAME=${PROJECT_NAME}-stage

PROTOS := coingecko-api

PHONY: generate
generate:
	$(foreach project,$(PROTOS), \
        mkdir -p pkg/$(project); \
        protoc --go_out=pkg/$(project) --go_opt=paths=import \
            --go-grpc_out=pkg/$(project) --go-grpc_opt=paths=import \
            --grpc-gateway_out=pkg/$(project) \
            --grpc-gateway_opt grpc_api_configuration=./api/$(project)/$(project).yaml \
            --grpc-gateway_opt allow_delete_body=true \
            api/$(project)/$(project).proto; \
        mv pkg/$(project)/github.com/${USER}/external-api/* pkg/$(project); \
        rm -r pkg/$(project)/github.com; \
    )


PHONY: clean
clean:
	 if docker inspect ${APP_NAME} > /dev/null 2>&1; then docker rm -f ${APP_NAME} && docker rmi -f ${APP_NAME}; else echo "Container not found."; fi

PHONY: go-build
go-build:
	GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/${PROJECT_NAME}
	mkdir -p bin
	mv main bin

PHONY: build
build:
	make go-build
	docker build --tag ${APP_NAME} .

PHONY: run
run:
	make clean
	make build
	docker run --name ${APP_NAME} --network=zoo -dp 7070:7070 -e GRPC_PORT=7070 -e PROXY_PORT=7071 -e SWAGGER_PORT=7072 -e METRICS_PORT=7073 -e APP_NAME=${APP_NAME} -e ENV=local ${APP_NAME}

PHONY: migrate-up
migrate-up:
	#docker exec -it practice psql -U postgres -c "create database crypto_service"
	goose -dir ./migrations postgres "postgres://postgres:postgres@localhost/crypto_service?sslmode=disable" up

PHONY: migrate-down
migrate-down:
	goose -dir ./migrations postgres "postgresql://user:pass@host:port/db?sslmode=disable" down

PHONY: generate-swagger
generate-swagger:
	protoc -I . --openapiv2_out ./ \
		--experimental_allow_proto3_optional=true \
		 --openapiv2_opt grpc_api_configuration=./api/coingecko-api/coingecko-api.yaml \
		--openapiv2_opt proto3_optional_nullable=true \
		--openapiv2_opt allow_delete_body=true \
		./api/coingecko-api/coingecko-api.proto

	mv api/coingecko-api/coingecko-api.swagger.json ./swaggerui/swagger_temp.json
	jq '. + {"host": "$(SERVER_HOST):$(PROXY_PORT)", "schemes": ["http"]}' ./swaggerui/swagger_temp.json > ./swaggerui/swagger.json
	rm ./swaggerui/swagger_temp.json

