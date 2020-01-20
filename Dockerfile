FROM golang:1.12.14-stretch
WORKDIR /app
COPY go.mod go.sum requirements.txt ./

RUN apt-get update && apt-get install -y \
    python3-pip

RUN pip3 install -r requirements.txt

# Download all dependencies.
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app

RUN go build -o validate .
# Expose port 5000 to the outside world


EXPOSE 5001
USER 1000
CMD ["./validate"]