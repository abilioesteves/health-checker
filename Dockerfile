# BUILD
FROM abilioesteves/gowebbuilder:v0.7.0 as builder

ENV p $GOPATH/src/github.com/abilioesteves/health-checker

ADD ./ ${p}
WORKDIR ${p}
RUN go get -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /hc main.go

# PKG
FROM alpine

COPY --from=builder /hc /

ENV HC_PORT "37441"
ENV HC_LOG_LEVEL "info"
ENV HC_TARGET_HEALTH_URL ""
ENV HC_TARGET_NAME ""

CMD [ "/hc", "start" ]
