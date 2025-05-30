
FROM	golang:latest	AS	build-stage
WORKDIR	/app
COPY	go.mod	go.sum	./
RUN	go mod download
COPY	cmd/testbin/*.go	./
RUN	CGO_ENABLED=0 GOOS=linux go build -o /testbin
ENTRYPOINT	["/testbin"]
