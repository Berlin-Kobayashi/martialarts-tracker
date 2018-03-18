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
 - make referenced by recursive
 - add dummy tests
 - delete all references
 - read all
 - extract service package to own package
 - cover empty property reflection panic
 - optimize expanding by adding a temporary map of IDs
 - Test Mongo repository
 - Let Mongo repository keep the result type instead of overriding it with bson.M
 - improve error handling in service to always respond with proper status codes