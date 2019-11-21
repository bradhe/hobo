FROM minideb
RUN mkdir /app
ADD ./bin/location-search /app/location-search
CMD ["/app/location-search"]
