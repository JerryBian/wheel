FROM golang:buster AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN VER=$(cat VERSION) && \
    GITHASH=$(git rev-parse --short HEAD) && \
    BUILDTIME=`date "+%Y-%m-%dT%H:%M:%S"` && \
    go build -v -ldflags "-X 'main.AppVer=$VER' -X 'main.BuildTime=$BUILDTIME' -X 'main.GitHash=$GITHASH'" -o /wheel

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /wheel /wheel

EXPOSE 9000

ENTRYPOINT ["/wheel"]