step by step:

1. install golang migrate via: go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest (go install, not go get)
2. in root folder type this: migrate create -ext sql -dir database/migration create_transactions_table
3. in newly created file, insert as the sql script in up, and drop table in down
4. run migration by writing this in terminal (match the ones in supabse): migrate -database "DB_CONN" -path database/migration up
5. once finish check in supabase for changes
6. Just in case you face error while doing migration/s, usually when deleting by down and then restart the up migration, this shows up: error: Dirty database version 20260207180240. Fix and force version. -> one solution is to manually search for transaction id in supabase, then set it false. finally repeat the up migration.