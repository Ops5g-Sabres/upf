FROM golang:bullseye as build

RUN apt update && apt upgrade -qy
RUN apt install -y \
    build-essential \
    golang \
    make \
    ca-certificates \
    protobuf-compiler \
    vim \
    git

RUN mkdir -p /go
ENV GOPATH /go
ENV GOBIN /go/bin
ENV PATH "$PATH:$GOBIN"
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

COPY . /upf
WORKDIR /upf

RUN git submodule update --recursive --init
WORKDIR /upf/pfcpiface/click_pb/moa

RUN make protobuf
RUN mkdir -p /sabres
RUN mkdir -p /mgmt
RUN mkdir -p /sdcore

FROM build as export
COPY --from=build /upf/pfcpiface/click_pb/moa/pkg/sabres/*.pb.go /sabres
COPY --from=build /upf/pfcpiface/click_pb/moa/pkg/mgmt/*.pb.go /mgmt
COPY --from=build /upf/pfcpiface/click_pb/moa/pkg/sdcore/*.pb.go /sdcore
