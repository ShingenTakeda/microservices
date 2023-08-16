#BUILD TINY DOCKER IMAGE
FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

CMD [ "/app/brokerApp" ]
