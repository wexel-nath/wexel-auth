#!/bin/bash
set -e

# TODO: wait for db to be healthy
sleep 20

# check for required vars
DB_HOST="${DB_HOST:?DB_HOST must be set}"
DB_NAME="${DB_NAME:?DB_NAME must be set}"
DB_PASS="${DB_PASS:?DB_PASS must be set}"
DB_PORT="${DB_PORT:?DB_PORT must be set}"
DB_USER="${DB_USER:?DB_USER must be set}"

connect() {
	PGPASSWORD="$DB_PASS" \
		 psql \
			-h "$DB_HOST" \
			-p "$DB_PORT" \
			-U "$DB_USER" \
			"$@"
}

run_db_query() {
	connect -tAc "$1"
}

run_query() {
	connect -d "$DB_NAME" -tAc "$1"
}

run_sql_file() {
	connect -d "$DB_NAME" -f "$1"
}

maybe_create_database() {
	echo "Creating database"

	if [[ $(run_db_query "SELECT 1 FROM pg_database WHERE datname = '$DB_NAME';") -eq 1 ]]; then
		echo "Database '$DB_NAME' exists"
	else
		run_db_query "CREATE DATABASE $DB_NAME;"
	fi
	echo
}

maybe_create_schema() {
	echo "Creating database schema"
	for file in /db/schema/*.sql; do
		if [[ "$file" =~ _(.+).sql ]]; then
			table_name=${BASH_REMATCH[1]}
		else
			continue
		fi

		echo "Running $file"
		if [[ $(run_query "SELECT 1 FROM pg_tables WHERE tablename = '$table_name';") -eq 1 ]]; then
			echo "Table '$table_name' exists"
		else
			run_sql_file "$file"
		fi
		echo
	done
}

maybe_run_updates() {
	echo "Running database migrations"
	for file in /db/updates/*.sql; do
		if [[ $(run_query "SELECT 1 FROM update WHERE update_id = '$file';") -eq 1 ]]; then
			continue
		fi

		echo "Running $file"
		run_sql_file "$file"

		echo "Saving update"
		run_query "INSERT INTO update (update_id) VALUES ('$file');"
		echo
	done
}

maybe_create_database
maybe_create_schema

# TODO: uncomment this once an update exists
#maybe_run_updates

echo "Database initialization complete."

exit 0
