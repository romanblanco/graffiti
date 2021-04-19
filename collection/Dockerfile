FROM golang:alpine
WORKDIR /app
ADD . /app/
RUN apk update && \
    apk add --update git && \
    apk add --update openssh
RUN go get github.com/google/open-location-code/go
RUN go get github.com/op/go-logging
RUN go get github.com/romanblanco/go-ipfs-api@e9a1904e1ee945c2be8373afdd7cec026e6e31a0
RUN go get github.com/rwcarlsen/goexif/exif
RUN go build -o graffiti
ENTRYPOINT ["./graffiti"]
