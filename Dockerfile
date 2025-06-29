# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# [START cloudrun_helloworld_dockerfile]
# [START run_helloworld_dockerfile]

# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:alpine3.21 AS builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.mod go.sum ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Install swag for API documentation generation.
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go

# Build the binary.
RUN mkdir build
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/server ./cmd/api

# Use the official Debian slim image for a lean production container.
FROM alpine:latest
RUN apk update && apk add --no-cache ca-certificates tzdata && update-ca-certificates
ENV TZ=Asia/Jakarta

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/build/server /app/server
COPY .env .env

RUN chmod +x /app/server
EXPOSE 3001

# Run the web service on container startup.
CMD ["/app/server"]