FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD bin/appbackend /
ADD db /db

EXPOSE 8080
EXPOSE 8081

ENTRYPOINT ["/appbackend"]
