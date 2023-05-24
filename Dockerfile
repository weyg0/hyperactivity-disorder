FROM debian:stretch-slim

WORKDIR /

COPY active-defense-scheduler-framework /usr/local/bin

CMD ["active-defense-scheduler-framework"]