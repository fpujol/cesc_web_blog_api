FROM golang:1.19-alpine

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o build ./...

# ENV GIN_MODE=release
ENV HTTP_SERVER_ADDRESS=0.0.0.0:5000
ENV DB_SOURCE=postgresql://postgres:postgres@192.168.100.110:5432/cesc_blog?sslmode=disable

CMD ["build/cmd"]