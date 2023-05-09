FROM golang

WORKDIR /event-horizon
COPY src /event-horizon
COPY go.mod /event-horizon
COPY Makefile /event-horizon
RUN go mod tidy
#CMD ["make", "run"]