FROM golang:1.8
RUN mkdir -p /go/src/server
RUN mkdir -p /var/pool
COPY . /go/src/server/
ENV PORT=8080 
RUN go get github.com/cheikhshift/gos
RUN cd /go/src/server && gos deps && gos --export && go install
EXPOSE 8080
Expose 8700

CMD server
# healthcheck requires docker 1.12 and up.
# HEALTHCHECK --interval=20m --timeout=3s \
#  CMD curl -f http://localhost:8080/ || exit 1