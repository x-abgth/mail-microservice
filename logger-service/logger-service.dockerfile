FROM alpine:3.17.1

RUN mkdir /app

COPY loggerServiceApp /app

CMD ["/app/loggerServiceApp"]