SQLC_DATABASES = postgresql

build: services/storage/sqlc

# Mark this as a phony target so it always runs.
.PHONY: services/storage/sqlc
services/storage/sqlc: \
	./services/storage/sqlc/sqlc.json \
	$(patsubst %,./services/storage/sqlc/%/schema.sql,$(SQLC_DATABASES)) \
	$(patsubst %,./services/storage/sqlc/%/queries.sql,$(SQLC_DATABASES))

	@sqlc generate -f $<
