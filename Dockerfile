# =================================================== #
# STAGE 1 - Build the project.                        #
# =================================================== #
FROM docker.io/golang:1.23.2-bookworm as builder

# Defaults for build
ARG GOARCH=amd64
ARG GOAMD64=v3

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN make build \
        binary_build_directory=. \
        binary_build_name=wex \
        GOOS=linux \
        GOARCH=${GOARCH} \
        GOAMD64=${GOAMD64}

# =================================================== #
# STAGE 2 - Build the production image.               #
# =================================================== #
FROM gcr.io/distroless/static:nonroot
LABEL authors="diego@diegoalves.info"

WORKDIR /api

COPY --from=builder /api/wex ./

ENTRYPOINT ["/api/wex"]
