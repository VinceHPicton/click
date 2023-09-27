# Click
The open source, non-profit dating app, focussed on finding people real relationships, not on profits.

### Key features:
- 50% male 50% female ratio with each region, if there is an excess of one gender, they have to wait in a queue until there is space
- 5 matches per account at a time - if you a new match after you have 5 you must remove some.

## Run locally
- Run /scripts/build.sh
- Run /scripts/up.sh

## Rebuild locally
- Stop the containers either with /scripts/down.sh or ctrl+C
- Run /scripts/clean.sh to reset state of containers
- Run /scripts/up.sh

* [Project layout](./docs/layout.md) -- this explains the rationale behind the layout

## Tools list
- atlas for DB migrations
- Docker
- Pgadmin
- sqlc