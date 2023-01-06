migratecreate:
	migrate -ext sql -dir db/migration -seq init_schema
migrateup:
	 migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down
createdb:
	docker exec -it postgres2 createdb simple_bank
dropdb:
	docker exec -it postgres2 dropdb simple_bank
accessdb:
	docker exec -it postgres2 psql  simple_bank
cc:
	go clean -testcache

.PHONY: accessdb migrateup migratedown migratecreate createdb dropdb cc