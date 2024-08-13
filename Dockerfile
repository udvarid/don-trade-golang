FROM ubuntu:latest

RUN apt update && apt install -y ca-certificates

COPY don-trade-golang ./
COPY ./html/*.html ./html/
RUN mkdir /db

EXPOSE 8080

ENTRYPOINT ["./don-trade-golang", "-environment=fly"]
