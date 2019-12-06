FROM golang:latest


# Copy the source files
COPY ./ ./

# Backend service runs on port 5000
EXPOSE 5000

ENTRYPOINT [ "./backend" ]  