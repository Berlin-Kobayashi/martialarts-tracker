# martialarts-tracker

## Setup

### Immutable Image
`docker build -t martialarts-tracker .`
`docker run martialarts-tracker`

### Dev Environment
`docker-compose up`

## Use Cases
-   add training unit

## Entities
-   training unit
    -   training series (e.g. "JKD I Dan")
    -   techniques
        -   type (e.g. "kick")
        -   name (e.g. "hook kick")
        -   description
    -   training methods
        -   type (e.g. "Counter")
        -   name (e.g. "Pak-Sao Game")
        -   description
        -   contains (e.g. "Pak Sao, Parry")
    -   exercise
        -   type (e.g. "Sparring")
        -   name (e.g. "Lead hand sparring")
        -   description
