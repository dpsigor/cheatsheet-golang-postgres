#!/bin/bash

complete -W "postgres createdb dropdb execdb migrateup migratedown sqlc lint test watch" ./scripts.bash
