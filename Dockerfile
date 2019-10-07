FROM golang
ADD . /go/src/github.com/Barrett17/market_prices
RUN go get -u github.com/gorilla/mux
RUN go install /go/src/github.com/Barrett17/market_prices
ENTRYPOINT /go/bin/market_prices
EXPOSE 8080
