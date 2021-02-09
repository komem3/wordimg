FROM golang:1.15
WORKDIR /go/src/github.com/komem3/word_rand_img/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/server

FROM alpine:latest
ENV TZ=Asia/Tokyo
ENV FONT_PATH=./font.ttf
WORKDIR /root
COPY font.ttf .
COPY --from=0 /go/src/github.com/komem3/word_rand_img/app .
CMD ["./app"]
