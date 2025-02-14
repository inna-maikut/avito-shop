SHELL:=/bin/bash

oapi-codegen:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./api/oapi-codegen.yaml ./api/schema.yaml

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

lint-ci:
	golangci-lint run ./... --out-format=github-actions --timeout=5m

generate:
	go generate ./...

test:
	go test -cover ./...

tidy:
	go mod tidy

make_jwt_keys:
	openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem
	echo "JWT_SECRET=\"`sed -E 's/\$$/\\\n/g' ecprivatekey.pem`\"" >> .env
	rm ecprivatekey.pem