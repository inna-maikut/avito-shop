SHELL:=/bin/bash

oapi-codegen:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=./api/oapi-codegen.yaml ./api/schema.yaml

run-local:
	go run ./cmd/server/main.go

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run --fix ./...

lint-ci:
	golangci-lint run ./... --out-format=github-actions --timeout=5m

generate:
	go generate ./...

test-cover-no-integration:
	go test -cover ./...

test-cover:
	go test --tags integration -cover ./...

test-api:
	go test --tags integration ./test/integration

test-repository:
	go test --tags integration ./internal/repository

test-total-cover-no-integration:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out && rm cover.out

test-total-cover:
	go test --tags integration ./... -coverprofile cover.out && go tool cover -func cover.out && rm cover.out

tidy:
	go mod tidy

make_jwt_keys:
	openssl ecparam -name prime256v1 -genkey -noout -out ecprivatekey.pem
	echo "JWT_SECRET=\"`sed -E 's/\$$/\\\n/g' ecprivatekey.pem`\"" >> .env
	rm ecprivatekey.pem

load-generate-targets:
	go run test/load/generate_targets.go

RATE=100
load-test-2-98:
	cat test/load/targets_2-98 | \
		vegeta attack -rate=${RATE}/s -lazy -format=http -duration=60s > test/load/results.2-98.${RATE}rps.bin
	cat test/load/results.2-98.${RATE}rps.bin | vegeta report -type=hdrplot > test/load/results.2-98.${RATE}rps.txt
	cat test/load/results.2-98.${RATE}rps.bin | vegeta report >> test/load/results.2-98.${RATE}rps.txt

load-test-10-90:
	cat test/load/targets_10-90 | \
		vegeta attack -rate=${RATE}/s -lazy -format=http -duration=60s > test/load/results.10-90.${RATE}rps.bin
	cat test/load/results.10-90.${RATE}rps.bin | vegeta report -type=hdrplot > test/load/results.10-90.${RATE}rps.txt
	cat test/load/results.10-90.${RATE}rps.bin | vegeta report >> test/load/results.10-90.${RATE}rps.txt

load-test-50-50:
	cat test/load/targets_50-50 | \
		vegeta attack -rate=${RATE}/s -lazy -format=http -duration=60s > test/load/results.50-50.${RATE}rps.bin
	cat test/load/results.50-50.${RATE}rps.bin | vegeta report -type=hdrplot > test/load/results.50-50.${RATE}rps.txt
	cat test/load/results.50-50.${RATE}rps.bin | vegeta report >> test/load/results.50-50.${RATE}rps.txt

load-test-plot:
	vegeta plot -title 2-98.${RATE}rps test/load/results.2-98.${RATE}rps.bin > test/load/results.plot.2-98.${RATE}rps.html
	vegeta plot -title 10-90.${RATE}rps test/load/results.10-90.${RATE}rps.bin > test/load/results.plot.10-90.${RATE}rps.html
	vegeta plot -title 50-50.${RATE}rps test/load/results.50-50.${RATE}rps.bin > test/load/results.plot.50-50.${RATE}rps.html
#	vegeta plot test/load/results.2-98.${RATE}rps.bin \
#		test/load/results.10-90.${RATE}rps.bin \
#        test/load/results.50-50.${RATE}rps.bin > test/load/results.plot.${RATE}rps.html