FROM golang:1.23

# Set the working directory to /app
WORKDIR /app
# Copy the server directory to the container
COPY . ./
# Download the dependencies
RUN go mod download
# Build the server
RUN go build -o /main
# expose the port
EXPOSE 8080

ENV SERVER_ADDRESS_ONE=server2
ENV SERVER_ADDRESS_TWO=server3

CMD [ "/main" ]