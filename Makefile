BASE_STACK = docker compose -f docker-compose.yml
INTEGRATION_TEST_STACK = docker compose -f docker-compose-integration-test.yml

compose-up: ### Run docker compose
		$(BASE_STACK) up --build -d
.PHONY: compose-up

compose-down: ### Down docker compose
		$(BASE_STACK) down
.PHONY: compose-down

bench-test: ### Run a temporary test database for bench-test, then deletes it
		$(INTEGRATION_TEST_STACK) up --build -d test_db test_migrate_up
		TEST_PG_URL=postgres://testuser:testpass@localhost:5433/wb_tech_test?sslmode=disable \
			go test -bench=. ./internal/repo/cache -count=1 -benchmem
		$(INTEGRATION_TEST_STACK) down -v
.PHONY: bench_test

test: ### Run tests
		$(INTEGRATION_TEST_STACK) up --build -d test_db test_migrate_up
		TEST_PG_URL=postgres://testuser:testpass@localhost:5433/wb_tech_test?sslmode=disable \
			go test -v -race ./internal/...
		$(INTEGRATION_TEST_STACK) down -v
.PHONY: test

deps: ### deps tidy + verify
		go mod tidy && go mod verify
.PHONY: deps