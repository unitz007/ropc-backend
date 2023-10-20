FROM golang:1.19
LABEL authors="charles"

RUN mkdir /app

WORKDIR /app

COPY ./ropc .

CMD ["./ropc"]