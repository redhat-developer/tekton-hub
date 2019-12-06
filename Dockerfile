FROM golang:latest

COPY . .

EXPOSE 5000
USER 1000

ENTRYPOINT [ "./backend" ]  