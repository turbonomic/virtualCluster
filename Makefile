OUTPUT_DIR=./_output

product: clean
	env GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/containerChain.linux ./cmd

build: clean
	go build -o ${OUTPUT_DIR}/containerChain ./cmd

test: clean
	@go test -v -race ./pkg/...
	
.PHONY: clean
clean:
	@: if [ -f ${OUTPUT_DIR} ] then rm -rf ${OUTPUT_DIR} fi
