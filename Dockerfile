FROM debian:stretch-slim

WORKDIR /

COPY active-defense-scheduler /usr/local/bin

CMD ["active-defense-scheduler"]