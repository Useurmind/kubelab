FROM golang as build

WORKDIR /app

RUN go get github.com/magefile/mage

# cache dependencies
COPY go.mod go.sum ./
RUN go mod download

ADD . .
RUN mage -compile magic
RUN mage Projects:Build

FROM ubuntu as run
WORKDIR /app

COPY --from=build /app/services/projects/bin/ /app/

CMD [ "/app/projects" ]