go get -u github.com/pressly/goose/cmd/goose
env $(cat db.env) bash -c 'goose -dir=db/migrations postgres "host=${MIGRATE_DB_URL} user=${POSTGRES_USER} dbname=${POSTGRES_DB} sslmode=disable password=${POSTGRES_PASSWORD}" up'
