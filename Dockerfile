FROM golang:1.12.4 AS builder

# Build arguments
ARG binary_name=main
    # See ./sample-data/go-os-arch.csv for a table of OS & Architecture for your base image
ARG target_os=linux
ARG target_arch=amd64

# Build the server Binary
WORKDIR /app
#WORKDIR /go/src/${GIT_SERVER}/${GIT_ORG}/${GIT_REPO}

COPY go.mod .
COPY go.sum .

RUN go mod download

# Seems duplicative, and ideally not needed
COPY . .

RUN rm -rf /app/build
RUN CGO_ENABLED=0 GOOS=${target_os} GOARCH=${target_arch} go build -a -o /app/build/${binary_name} main.go

RUN ls /app

#-----------------------------------------------------------------------------------------------------------------------

FROM centos:7

LABEL author="Benjamin Smith"
COPY --from=builder ./app/build/main /usr/bin/main
RUN ["chmod", "-R", "+x", "/usr/bin/main"]

ENTRYPOINT ["tail", "-f", "/dev/null"]
