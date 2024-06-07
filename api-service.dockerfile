FROM alpine:latest

RUN mkdir /app
RUN mkdir /csv
RUN mkdir /data

COPY csv/Player.csv /csv
COPY player_rest_server /app

ENV USE_MYSQL=1

CMD [ "/app/player_rest_server" ]
