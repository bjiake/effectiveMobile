FROM golang:1.22.2

RUN mkdir /apps
WORKDIR /apps
COPY . /apps
RUN go build /apps/cmd/app
CMD ["./app"]