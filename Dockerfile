FROM bitnami/minideb:stretch
RUN mkdir /app
ADD ./bin/location-search /app/location-search
CMD ["/app/location-search"]
