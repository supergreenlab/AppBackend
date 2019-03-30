FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD supergreenpromproxy /

EXPOSE 8080

ENTRYPOINT ["/supergreenpromproxy"]
