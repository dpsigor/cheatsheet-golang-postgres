#!/bin/bash

complete -W "postgres createdb dropdb execdb migrateup migratedown sqlc lint test" ./scripts.bash
