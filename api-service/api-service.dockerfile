FROM alpine:latest

RUN mkdir /app
RUN mkdir /csv

COPY player_server /app
COPY csv/Player.csv /csv

CMD [ "/app/player_server" ]
