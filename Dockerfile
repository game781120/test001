FROM golang:1.21.5-alpine
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o /app/digital_visitor

ENTRYPOINT ["./digital_visitor"]