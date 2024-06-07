# Intro2Gin
Minute REST-API built while learning the Gin Framework

## Table Of Contents
- [Running the project](#running-the-project)
  - [Requirements](#requirements)
  - [Instructions](#instructions)
- [API Routes](#api-routes)

## Running the project

#### Requirements
- Go
- PostgreSQL
- Make
- Git

#### Instructions
- Clone the project.
- Create a database.
- Use album.sql to create the tables and insert dummy data.
- Set the following environmental variables.
```bash
export DB_NAME=<your db name>
export DB_PASS=<your postgres password>
export DB_USER=<your postgres username>
```
- Cd into the project and run the following command
```bash
make run
```

- Make request to localhost:8080/albums

## API Routes

- GET /albums

  Returns a list of all albums

- POST /albums

  Adds a new album

  ### Body (json)
  ```json
  {
    "title": "Kolos",
    "artist": "Meshuggah",
    "price": 45.78
  }
  ```

- GET /albums/:id

  Returns an album with the specified id

  Example: localhost:8080/albums/3