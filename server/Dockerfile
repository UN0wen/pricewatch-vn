from golang:alpine

WORKDIR /build
COPY . .
RUN go mod vendor
RUN go build -o server .

WORKDIR /bin

RUN cp /build/server .
# static react page
COPY ./build/ .
EXPOSE 3000

CMD ["/bin/server"]
