FROM golang:latest


RUN go get github.com/chetan-rns/backend

# Copy the source files
COPY go/src/github.com/chetan-rns/backend .
RUN go mod download


RUN go build -o backend ./cmd/



# Backend service runs on port 5000
EXPOSE 5000
USER 1000
ENTRYPOINT [ "./backend" ]  
