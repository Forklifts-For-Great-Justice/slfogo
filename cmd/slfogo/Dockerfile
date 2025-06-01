FROM	golang:latest	AS	build-stage

WORKDIR	/app
COPY	cmd/slfogo/go.mod	cmd/slfogo/go.sum	./
RUN	go mod download
COPY	cmd/slfogo/*.go	./
RUN	CGO_ENABLED=0 GOOS=linux go build -o /slfogo
FROM	alpine:latest
WORKDIR	/
COPY	--from=build-stage /slfogo	/slfogo
ENTRYPOINT	["/slfogo"]
