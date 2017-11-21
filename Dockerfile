FROM busybox

# the CA certs are downloaded from https://curl.haxx.se/docs/caextract.html
ADD testdata/ca-certificates.crt /etc/ssl/certs/

COPY meli /

ENTRYPOINT ["/meli"]