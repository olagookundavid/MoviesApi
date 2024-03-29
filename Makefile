# ==================================================================================== # 
# HELPERS 
# ==================================================================================== #

## help: print this help message
.PHONY: help 
help: 
	@echo 'Usage:' 
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm 
confirm: 
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== # 
# DEVELOPMENT 
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo 'Starting server...'
	go run ./cmd/api

## db/psql: connect to the database using psql and docker
.PHONY: db/psql
db/psql:
	@echo 'connecting to db'
	docker exec -it post-db bash
	psql postgres://greenlight:greenlight@localhost/greenlight?sslmode=disable
 
## db/migrate/up: apply all up database migrations
.PHONY: db/migrate/up
db/migrate/up:
	@echo 'Running up migrations...' 
	goose postgres postgres://greenlight:greenlight@localhost/greenlight up
	
# ==================================================================================== # 
# QUALITY CONTROL 
# ==================================================================================== # 
## audit: tidy dependencies and format, vet and test all code 

.PHONY: audit 
audit: vendor
	@echo 'Formatting code...' 
	go fmt ./... 
	@echo 'Vetting code...' 
	go vet ./... 
	staticcheck ./... 
	@echo 'Running tests...' 
	go test -race -vet=off ./...

## vendor: tidy and vendor dependencies 
.PHONY: vendor 
vendor: 
	@echo 'Tidying and verifying module dependencies...' 
	go mod tidy 
	go mod verify 
	@echo 'Vendoring dependencies...' 
	go mod vendor

# ==================================================================================== # 
# BUILD 
# ==================================================================================== # 

## build/api: build the cmd/api application 
.PHONY: build/api 
build/api: 
	@echo 'Building cmd/api...' 
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api