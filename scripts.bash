#!/bin/bash

case $1 in

  postgres)
    docker run --rm --name postgres14 -v $HOME/db/psql:/var/lib/postgresql/data -p 5432:5432 -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -d postgres:14.4-alpine
    ;;

  createdb)
    docker exec postgres14 createdb --username=root --owner=root simple_bank
    ;;

  dropdb)
    docker exec postgres14 dropdb simple_bank
    ;;

  execdb)
    docker exec -it postgres14 psql -U root simple_bank
    ;;

  migrateup)
    migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up
    ;;

  migratedown)
    migrate -path db/migration/ -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
    ;;

  sqlc)
    sqlc generate
    ;;

  lint)
    golint ./...
    ;;

  test)
    go test -v -cover ./...
    ;;

  watch)
    find  -type f -name "*.go" | entr -r go run .
    ;;

  *)
    echo "Usage:
    postgres
    createdb
    dropdb
    execdb
    migrateup
    migratedown
    sqlc
    lint
    test
    watch"
    ;;
esac


