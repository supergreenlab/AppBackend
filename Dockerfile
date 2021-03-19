FROM debian

RUN apt-get update && \
    apt-get install -y libmagickwand-dev
COPY assets /usr/local/share/appbackend

COPY ca-certificates.crt /etc/ssl/certs/
COPY bin/appbackend /
COPY db /db

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/appbackend"]
