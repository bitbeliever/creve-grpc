FROM golang:1.7 AS builder

COPY . .

RUN go build -o main

FROM busybox

COPY --from=build main /app

RUN ls
RUN ls /app

CMD ./app/main