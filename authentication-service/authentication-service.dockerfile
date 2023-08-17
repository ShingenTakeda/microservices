#BUILD TINY DOCKER IMAGE
FROM alpine:latest

RUN mkdir /app

COPY authApp /app

CMD [ "/app/authApp" ]
