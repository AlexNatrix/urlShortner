#Here 

FROM golang:latest AS build-stage

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux 

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o main ./cmd/service

#run tests in container
FROM build-stage AS run-test-stage
RUN go test -v ./...


#Deploy
FROM gcr.io/distroless/base-debian11 as build-release-stage

WORKDIR /

COPY --from=build-stage /build /build

USER nonroot:nonroot

EXPOSE 8082

ENTRYPOINT [ "./build/service.exe" ]