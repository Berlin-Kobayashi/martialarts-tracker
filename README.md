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
 - get referenced by
 - delete all references
 - read all
 - extract service package to own package
 - cover empty property reflection panic
 - optimize expanding by adding a temporary map of IDs
 - Test Mongo repository
 - improve error handling in service to always respond with proper status codes