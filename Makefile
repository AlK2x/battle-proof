include .env
export

.PHONY: init-db
init-db:
	@echo "-- building mysql db init scripts"
	bin/initdb.sh