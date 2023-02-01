# because we are already building the app from Makefile
FROM alpine:3.17.1

RUN mkdir /app

COPY mailerApp /app
COPY templates /templates

CMD ["/app/mailerApp"]