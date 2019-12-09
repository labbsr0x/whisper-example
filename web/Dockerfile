FROM abilioesteves/gowebbuilder:1.2.0 as builder

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /web 

## PKG
FROM alpine

RUN apk update \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates

COPY --from=builder /web /
COPY assets /assets

CMD [ "/web" ]