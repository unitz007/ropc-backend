FROM golang:1.17-alpine
LABEL authors="charles"

RUN mkdir /app

WORKDIR /app

COPY ./ropc .

CMD ["./ropc"]