# -=-=-=-=-=-=- Compile Image -=-=-=-=-=-=-

FROM golang:1.16 AS stage-compile

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./cmd/govanity
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/govanity

# -=-=-=-=-=-=- Final Image -=-=-=-=-=-=-

FROM alpine:latest 

WORKDIR /root/
COPY --from=stage-compile /go/src/app/govanity ./
COPY --from=stage-compile /go/src/app/cmd/internal/templates/*.gotmpl ./cmd/internal/templates/

RUN apk --no-cache add ca-certificates

EXPOSE 8080

ENTRYPOINT [ "./govanity" ]  