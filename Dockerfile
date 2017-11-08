FROM golang:alpine
COPY proxy.go /app/
COPY run.sh /run.sh
EXPOSE 8080
CMD ["/run.sh"]