go get -u github.com/pressly/goose/cmd/goose
name="$@"
MIGRATION_NAME="$@" env $(cat db.env) bash -c 'goose -dir=db/migrations postgres "host=${MIGRATE_DB_URL} user=${POSTGRES_USER} dbname=${POSTGRES_DB} sslmode=disable password=${POSTGRES_PASSWORD}" create "${MIGRATION_NAME}" sql'
