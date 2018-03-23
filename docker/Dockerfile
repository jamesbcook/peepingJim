FROM golang:1.10.0-alpine3.7

RUN apk --no-cache add git wget make tar nodejs nodejs-npm bash

RUN npm install -g phantomjs@2.1.1 --unsafe-perm

RUN wget -qO- "https://github.com/dustinblackman/phantomized/releases/download/2.1.1/dockerized-phantomjs.tar.gz" | tar xz -C /

RUN mkdir -p src/github.com/jamesbcook/

WORKDIR src/github.com/jamesbcook/

RUN git clone https://github.com/jamesbcook/peepingJim.git

WORKDIR peepingJim

RUN make docker

WORKDIR bin

CMD ["./peepingJim", "-h"]