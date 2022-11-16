FROM golang:1.18.6 AS builder

# Build arguments
ARG binary_name=UnrealVersionSelector
    # See ./sample-data/go-os-arch.csv for a table of OS & Architecture for your base image
ARG target_os=linux
ARG target_arch=amd64
ARG VERSION=dev
# Build the server Binary
WORKDIR /app
#WORKDIR /go/src/${GIT_SERVER}/${GIT_ORG}/${GIT_REPO}

COPY go.mod .
COPY go.sum .

RUN go mod download

# Seems duplicative, and ideally not needed
COPY . .

RUN rm -rf /app/build

RUN CGO_ENABLED=0 GOOS=${target_os} GOARCH=${target_arch} go build -a -ldflags " -X github.com/Benbentwo/UnrealGameVersionUpdater/pkg/version.Version=$(VERSION)" -o /app/build/${binary_name} main.go

#-----------------------------------------------------------------------------------------------------------------------

FROM centos:7

LABEL author="Benjamin Smith"
COPY --from=builder ./app/build/UnrealVersionSelector /usr/bin/UnrealVersionSelector
RUN ["chmod", "+x", "/usr/bin/UnrealVersionSelector"]

ENTRYPOINT ["/usr/bin/UnrealVersionSelector"]
