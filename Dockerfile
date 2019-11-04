FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD bin/appbackend /

EXPOSE 8080

ENTRYPOINT ["/appbackend"]
