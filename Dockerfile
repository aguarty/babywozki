FROM golang:1.14 AS stage
ARG PACKAGE_NAME
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -ldflags="-w -s -X main.version=${VERSION} -X main.commitHash=${COMMIT}" -o /babywozki -a -installsuffix cgo cmd/*.go

FROM alpine:3.9
COPY /static /static
COPY --from=stage /babywozki /opt/application
ENTRYPOINT [ "/opt/application" ]
