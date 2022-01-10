FROM golang:alpine

WORKDIR /go/src/devtipme
COPY . .
RUN go get

CMD ["go","run","main.go"]

#  docker build -t backend .
# docker run -p 5000:5000 --name backend -d backend