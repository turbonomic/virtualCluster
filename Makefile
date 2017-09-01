OUTPUT_DIR=./_output

bin=vCluster
product: clean
	env GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_DIR}/${bin}.linux ./cmd

build: clean
	go build -o ${OUTPUT_DIR}/${bin} ./cmd

test: clean
	@go test -v -race ./pkg/...
	
clean:
	@: if [ -f ${OUTPUT_DIR} ] then rm -rf ${OUTPUT_DIR} fi
