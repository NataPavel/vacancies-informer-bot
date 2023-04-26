FROM golang
WORKDIR /vac_informer_tgbot
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-vac-info-bot
CMD ["/docker-vac-info-bot"]

