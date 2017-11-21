FROM busybox
COPY meli /
ENTRYPOINT ["/meli"]