FROM	golang:latest	AS	build-stage

WORKDIR	/app
COPY	go.mod	go.sum	./
RUN	go mod download
COPY	cmd/slfogo/*.go	./
RUN	CGO_ENABLED=0 GOOS=linux go build -o /slfogo
ENTRYPOINT	["/slfogo"]
