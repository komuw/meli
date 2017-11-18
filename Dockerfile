FROM scratch
COPY meli /
ENTRYPOINT ["/meli -up"]