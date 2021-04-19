FROM golang:1.16.2 as builder

WORKDIR $GOPATH/src/github.com/mainak90/k8sJobAdmissionWebHook
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/job-admission

FROM scratch
COPY --from=builder /go/bin/job-admission /go/bin/job-admission
ENTRYPOINT ["/go/bin/job-admission"]