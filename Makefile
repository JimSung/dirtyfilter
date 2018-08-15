all: g complie
g:
	@echo "generate protocols..."
	@protoc --go_out=plugins=grpc:. proto/*.proto
complie:
	@echo "complie binaries..."
	@go install 1>/dev/null
