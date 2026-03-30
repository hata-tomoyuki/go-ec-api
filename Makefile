.PHONY: setup test test-unit test-integration

## setup: git hooks を有効化する（初回クローン後に一度だけ実行）
setup:
	git config core.hooksPath .githooks
	@echo "✅ Git hooks enabled (.githooks/pre-push)"

## test: 全テストを実行（統合テスト含む・Docker 必要）
test:
	go test ./... -cover -timeout 120s

## test-unit: ユニットテストのみ実行（Docker 不要・高速）
test-unit:
	go test ./... -run 'Test[^I]' -timeout 60s -cover

## test-integration: 統合テストのみ実行（Docker 必要）
test-integration:
	go test ./... -run TestIntegration -v -timeout 120s
