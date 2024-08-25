FROM ubuntu:latest

RUN apt update && apt install -y ca-certificates

COPY don-trade-golang ./
COPY ./html/index.html ./html/
COPY ./html/detailed.html ./html/
COPY ./html/logged_in.html ./html/
COPY ./html/admin.html ./html/
COPY ./html/user.html ./html/
COPY ./static/* ./static/
RUN mkdir /db

EXPOSE 8080

ENTRYPOINT ["./don-trade-golang", "-environment=fly"]
