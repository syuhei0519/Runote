generate:
	npx prisma generate

# Prisma マイグレーション名を手動で渡せる
migrate:
	@if [ -z "$(name)" ]; then \
		read -p "マイグレーション名を入力してください: " name; \
		npx -y prisma migrate dev --name $$name && make generate; \
	else \
		npx -y prisma migrate dev --name $(name) && make generate; \
	fi

test-init:
	dotenv -e ./test.env -- printenv DATABASE_URL && dotenv -e ./test.env -- npx prisma db push --skip-generate

test: test-init
	npm test