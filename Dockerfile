FROM golang:1.13 as build

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go get
RUN CGO_ENABLED=0 go build -a -o bmco

FROM golang:alpine

RUN mkdir /app
COPY --from=build /app/bmco /app
WORKDIR /app

ENTRYPOINT [ "./bmco" ]
CMD "./bmco"