FROM scratch as runtime
#FROM golang:1.17.0-alpine
WORKDIR app
COPY /artifacts/k8s-patcher .
CMD ["./k8s-patcher"]
#ENTRYPOINT ["/bin/sh","-c","sleep infinity"]