FROM golang:1.19


# Set destination for COPY
WORKDIR /go/src/ropc-backend

## Download Go modules
#COPY go.mod go.sum ./
#RUN go mod download
#
#
#COPY *.go ./
#
## Build
#RUN go build -o ./ropc-backend .

COPY . .

RUN go get -d ./...
RUN go build -o ./ropc-backend -buildvcs=false

# Run
CMD ["./ropc-backend"]