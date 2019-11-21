FROM bitnami/minideb:stretch
RUN mkdir /app
ADD ./bin/hobo /app/hobo
CMD ["/app/hobo"]
