
FROM	golang:latest	AS	build-stage
WORKDIR	/app
COPY	cmd/testbin/go.mod	./
RUN	go mod download
COPY	cmd/testbin/*.go	./
RUN	CGO_ENABLED=0 GOOS=linux go build -o /testbin
FROM	scratch
WORKDIR	/
COPY	--from=build-stage /testbin	/testbin
ENTRYPOINT	["/testbin"]
