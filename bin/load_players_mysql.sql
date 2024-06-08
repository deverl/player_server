USE rest_server;

DROP TABLE IF EXISTS players;

CREATE TABLE IF NOT EXISTS players (
    playerID     VARCHAR (60) PRIMARY KEY UNIQUE,
    birthYear    INTEGER,
    birthMonth   INTEGER,
    birthDay     INTEGER,
    birthCountry VARCHAR(60),
    birthState   VARCHAR(60),
    birthCity    VARCHAR(60),
    deathYear    INTEGER,
    deathMonth   INTEGER,
    deathDay     INTEGER,
    deathCountry VARCHAR(60),
    deathState   VARCHAR(60),
    deathCity    VARCHAR(60),
    nameFirst    VARCHAR(60),
    nameLast     VARCHAR(60),
    nameGiven    VARCHAR(60),
    weight       INTEGER,
    height       INTEGER,
    bats         VARCHAR(60),
    throws       VARCHAR(60),
    debut        DATE,
    finalGame    DATE,
    retroID      VARCHAR(60),
    bbrefID      VARCHAR(60)
);

CREATE TABLE IF NOT EXISTS config (
    id       INTEGER PRIMARY KEY UNIQUE,
    fileHash VARCHAR(255)
);

-- INSERT INTO config (id, fileHash) VALUES (1, "");:

LOAD DATA LOCAL INFILE '../api-service/csv/Player.csv' INTO TABLE players
FIELDS TERMINATED BY ',' 
 OPTIONALLY ENCLOSED BY '"'
 ESCAPED BY '"'
 LINES TERMINATED BY '\n'
IGNORE 1 LINES;
