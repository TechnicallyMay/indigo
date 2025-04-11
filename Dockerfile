# Build the Go app
FROM golang:1.22 AS build-stage

WORKDIR /indigo

COPY go.mod go.sum ./
COPY ./ ./
COPY ./tmpl /tmpl
COPY ./static /static

RUN echo "Listing file structure for build context."
RUN ls

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /indigo-run

FROM build-stage AS run-test-stage
RUN go test -v ./...

# Now make the actual container
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /indigo-run /indigo-run
COPY --from=build-stage /tmpl /tmpl
COPY --from=build-stage /static /static

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/indigo-run"]

