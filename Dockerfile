FROM golang:1.20.4-bullseye

RUN apt-get update -yq \
  && apt-get -yq install curl ca-certificates zip make \
  && curl -L https://deb.nodesource.com/setup_18.x | bash \
  && apt-get update -yq \
  && apt-get install -yq nodejs

RUN npm install -g pnpm

WORKDIR /app

COPY . .

RUN make install

RUN make build

CMD [ "make", "run" ]
