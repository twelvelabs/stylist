FROM alpine:3.18

COPY stylist /usr/local/bin/stylist

CMD ["/usr/local/bin/stylist"]
