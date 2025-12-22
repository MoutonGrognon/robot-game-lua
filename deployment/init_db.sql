CREATE USER referee_user WITH PASSWORD :v1;
CREATE USER matchmaker_user WITH PASSWORD :v2;
CREATE USER rglua_user WITH PASSWORD :v3;
CREATE DATABASE rglua;

\c rglua;

CREATE OR REPLACE FUNCTION load_bot(path TEXT) RETURNS TEXT AS $$
  SELECT CAST(pg_read_file($1) AS TEXT);
$$ LANGUAGE sql SECURITY DEFINER;
ALTER FUNCTION load_bot(TEXT) OWNER TO postgres;

-- TODO: reduce name size
CREATE TABLE users (id UUID PRIMARY KEY, name VARCHAR (20));
CREATE TABLE bots (id UUID PRIMARY KEY, name VARCHAR (16), script TEXT, userId UUID REFERENCES users (id) , userName VARCHAR (20));
CREATE TABLE ranking (id UUID PRIMARY KEY, botId UUID REFERENCES bots (id), botName VARCHAR (16), elo INTEGER, winCount INTEGER, drawCount INTEGER, lossCount INTEGER);
CREATE TABLE matches (id UUID PRIMARY KEY, blueBotId UUID REFERENCES bots (id), redBotId UUID REFERENCES bots (id), blueBotName VARCHAR (16), redBotName VARCHAR (16), blueUserName VARCHAR (20), redUserName VARCHAR (20), date TIMESTAMP, compressedGame BYTEA, blueScore INTEGER, redScore INTEGER, ranked BOOLEAN);
-- Make pagination more efficient for matches
CREATE TABLE matches_count (count BIGINT);

CREATE OR REPLACE FUNCTION matches_count_update() RETURNS trigger AS $$
  BEGIN
    IF TG_OP = 'INSERT' THEN
      UPDATE matches_count SET count = count + 1;
      RETURN NEW;
    END IF;
    IF TG_OP = 'DELETE' THEN
      UPDATE matches_count SET count = count - 1;
      RETURN OLD;
    END IF;
    RETURN NULL;
  END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

CREATE CONSTRAINT TRIGGER matches_count_trigger 
  AFTER INSERT OR DELETE on matches
  DEFERRABLE INITIALLY DEFERRED
  FOR EACH ROW
    EXECUTE PROCEDURE matches_count_update();

GRANT SELECT ON TABLE bots TO referee_user;

GRANT SELECT ON TABLE bots TO matchmaker_user;
GRANT SELECT ON TABLE users TO matchmaker_user;
GRANT SELECT, INSERT, UPDATE ON TABLE ranking TO matchmaker_user;
GRANT SELECT, INSERT ON TABLE matches TO matchmaker_user;

-- TODO: add bouncer user
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE bots TO rglua_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE users TO rglua_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE ranking TO rglua_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE matches TO rglua_user;
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE matches_count TO rglua_user;

INSERT INTO matches_count(count) VALUES (0);
