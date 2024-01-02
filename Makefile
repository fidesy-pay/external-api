# Constants

PROJECT_NAME=external-api
USER=fidesy-pay


PHONY: generate
generate:
	mkdir -p pkg/coingecko-api
	protoc --go_out=pkg/coingecko-api --go_opt=paths=import \
			--go-grpc_out=pkg/coingecko-api --go-grpc_opt=paths=import \
			--grpc-gateway_out=pkg/coingecko-api \
            --grpc-gateway_opt grpc_api_configuration=./api/${PROJECT_NAME}/coingecko-api.yaml \
            --grpc-gateway_opt allow_delete_body=true \
			api/${PROJECT_NAME}/coingecko-api.proto
	mv pkg/coingecko-api/github.com/${USER}/${PROJECT_NAME}/* pkg/coingecko-api
	rm -r pkg/coingecko-api/github.com

PHONY: clean
clean:
	 if docker inspect ${PROJECT_NAME} > /dev/null 2>&1; then docker rm -f ${PROJECT_NAME} && docker rmi -f ${PROJECT_NAME}; else echo "Container not found."; fi

PHONY: go-build
go-build:
	GOOS=linux GOARCH=amd64 go build -o ./main ./cmd/${PROJECT_NAME}
	mkdir -p bin
	mv main bin

PHONY: build
build:
	make go-build
	docker build --tag ${PROJECT_NAME} .

PHONY: run
run:
	make clean
	make build
	docker run --name ${PROJECT_NAME} --network=zoo -dp 7070:7070 -e GRPC_PORT=7070 -e PROXY_PORT=7071 -e SWAGGER_PORT=7072 -e METRICS_PORT=7073 -e APP_NAME=${PROJECT_NAME} -e ENV=local ${PROJECT_NAME}

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
      --openapiv2_opt grpc_api_configuration=./api/$(PROJECT_NAME)/$(PROJECT_NAME).yaml \
      --openapiv2_opt proto3_optional_nullable=true \
      --openapiv2_opt allow_delete_body=true \
      ./api/$(PROJECT_NAME)/$(PROJECT_NAME).proto

	mv api/$(PROJECT_NAME)/$(PROJECT_NAME).swagger.json ./swaggerui/swagger_temp.json
	jq '. + {"host": "$(APP_NAME).fidesy.xyz:$(PROXY_PORT)", "schemes": ["http"]}' ./swaggerui/swagger_temp.json > ./swaggerui/swagger.json
	rm ./swaggerui/swagger_temp.json
