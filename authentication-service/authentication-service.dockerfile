FROM alpine:3.17.1

RUN mkdir /app

COPY authApp /app

CMD ["/app/authApp"]