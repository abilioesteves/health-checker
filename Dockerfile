# BUILD
FROM abilioesteves/gowebbuilder:v0.7.0 as builder

ENV p $GOPATH/src/github.com/labbsr0x/health-checker

ADD ./ ${p}
WORKDIR ${p}
RUN go get -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /hc main.go

# PKG
FROM alpine

COPY --from=builder /hc /

CMD [ "/hc", "start" ]
