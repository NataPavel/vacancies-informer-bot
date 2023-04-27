FROM golang
WORKDIR /vac_informer_tgbot
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o docker-vac-info-bot cmd/main.go
CMD ["./docker-vac-info-bot"]

