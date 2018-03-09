# martialarts-tracker

Tracking service for martial arts training.

## Setup

### Immutable Image
`docker build -t martialarts-tracker .`
`docker run martialarts-tracker`

### Dev Environment
`docker-compose up`

## Use Cases
-   add training unit

## Entities
See examples directory


## TODOs
 - make CRUD complete
 - read all
 - extract service package to own package
 - cover empty property reflection panic
 - optimize expanding by adding a temporary map of IDs
