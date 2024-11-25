FROM golang:1.23.1

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
ENV PATH="/go/bin:${PATH}"
WORKDIR /app

COPY ./RatesMicroservice/go.mod ./RatesMicroservice/go.sum ./ ./

RUN go mod download

COPY ./RatesMicroservice .

RUN go build -o rates

 EXPOSE 50051

CMD ["./rates"]