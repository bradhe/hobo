FROM bitnami/minideb:stretch
RUN apt update && apt install -y ca-certificates
RUN mkdir /app
ADD ./bin/hobo /app/hobo
CMD ["/app/hobo"]
