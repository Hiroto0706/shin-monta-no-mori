DATABASE_URL=postgresql://postgres:password@localhost:5432/shin-monta-no-mori?sslmode=disable
TEST_DATABASE_URL=postgresql://postgres:password@localhost:5432/shin-monta-no-mori-test?sslmode=disable

.PHONY: createdb
createdb:
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec -it shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test

.PHONY: dropdb
dropdb:
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres  shin-monta-no-mori
	docker exec -it shin-monta-no-mori-db dropdb --username=postgres shin-monta-no-mori-test

.PHONY: new_migration
new_migration:
	migrate create -ext sql -dir server/internal/db/migration -seq $(name)

.PHONY: migrateup
migrateup:
	migrate -path server/internal/db/migration -database "$(DATABASE_URL)" -verbose up
	migrate -path server/internal/db/migration -database "$(TEST_DATABASE_URL)" -verbose up

.PHONY: migrateup1
migrateup1:
	migrate -path server/internal/db/migration -database "$(DATABASE_URL)" -verbose up 1
	migrate -path server/internal/db/migration -database "$(TEST_DATABASE_URL)" -verbose up 1

.PHONY: migratedown
migratedown:
	migrate -path server/internal/db/migration -database "$(DATABASE_URL)" -verbose down
	migrate -path server/internal/db/migration -database "$(TEST_DATABASE_URL)" -verbose down

.PHONY: migratedown1
migratedown1:
	migrate -path server/internal/db/migration -database "$(DATABASE_URL)" -verbose down 1
	migrate -path server/internal/db/migration -database "$(TEST_DATABASE_URL)" -verbose down 1

.PHONY: dc-up
dc-up:
	docker compose up -d

.PHONY: dc-down
dc-down:
	docker compose down

.PHONY: serve
serve:
	cd ./server && air -c .air.toml

.PHONY: front
front:
	cd ./client && npm run dev

.PHONY: sqlc
sqlc:
	cd server/ && sqlc generate

.PHONY: test
test:
	# テスト実行環境の構築
	docker exec shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test
	migrate -path server/internal/db/migration -database "$(TEST_DATABASE_URL)" -verbose up
	mkdir -p ./server/coverage

	# 各サブディレクトリのテストを実行し、個別のカバレッジファイルを生成
	cd server && go test ./api/admin/... -coverprofile=./coverage/api_admin.out
	cd server && go test ./api/user/... -coverprofile=./coverage/api_user.out
	cd server && go test ./api/middleware/... -coverprofile=./coverage/api_middleware.out
	cd server && go test ./pkg/... -coverprofile=./coverage/pkg.out
	cd server && go test ./internal/db/... -coverprofile=./coverage/db.out
	cd server && go test ./internal/domains/... -coverprofile=./coverage/domains.out

	# カバレッジファイルの結合
	echo "mode: set" > ./server/coverage/coverage.out
	tail -n +2 ./server/coverage/api_admin.out >> ./server/coverage/coverage.out
	tail -n +2 ./server/coverage/api_user.out >> ./server/coverage/coverage.out
	tail -n +2 ./server/coverage/api_middleware.out >> ./server/coverage/coverage.out
	tail -n +2 ./server/coverage/pkg.out >> ./server/coverage/coverage.out
	tail -n +2 ./server/coverage/db.out >> ./server/coverage/coverage.out
	tail -n +2 ./server/coverage/domains.out >> ./server/coverage/coverage.out

	# テスト結果の集計・出力
	cd server && go tool cover -func=./coverage/coverage.out > ./coverage/report.txt
	cd server && go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html
	./tools/aggregate_coverage.sh ./server/coverage/report.txt

# テストが途中で失敗したなどの理由でテスト環境が汚れてしまった時に使う
.PHONY: test-reset
test-reset:
	docker exec shin-monta-no-mori-db dropdb --username=postgres --if-exists shin-monta-no-mori-test
	docker exec shin-monta-no-mori-db createdb --username=postgres --owner=postgres shin-monta-no-mori-test
	migrate -path server/internal/db/migration -database "$(TEST_DATABASE_URL)" -verbose up

.PHONY: update-cors-setting
update-cors-setting:
	cd client/ && gsutil cors set cors-config.json gs://shin-monta-no-mori

.PHONY: deploy-server
deploy-server:
	cd server/ && flyctl deploy