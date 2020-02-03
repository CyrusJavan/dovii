# dovii
Namesake [parachromis dovii](https://en.wikipedia.org/wiki/Parachromis_dovii "Wikipedia")

A key-value data store used for learning more about the internals of a database. The goal is to incrementally improve the robustness and speed of the database, ideally reaching a point of an eventually consistent data store like Cassandra or Redis.

## How to run

- Need to have a local Docker daemon running and docker-compose installed
- Run with 3 dovii instances `docker-compose up --scale dovii=3`

## API

### GET
- `curl --header "Host: dovii.local" http://localhost/<KEY>`

### SET
- `curl --header "Host: dovii.local" -X POST http://localhost/<KEY>/<VALUE>`

## Changelog

- String keys and values stored in memory
- String keys and values stored in json on disk
- Bitcask support
- Deploy multiple instances with docker-compose, load balanced behind traefik
  - Currently does not return correct values for GETs, this is expected, now we do RAFT to get values consistent across cluster
