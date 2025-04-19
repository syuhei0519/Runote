.PHONY: clean-all

# 🔥 Docker関連すべて削除（危険）
clean-all:
	docker-compose down --volumes --remove-orphans || true
	docker rm -f $$(docker ps -aq) || true
	docker volume prune -f || true
	docker network prune -f || true
	docker image prune -a -f || true
	docker system prune -a --volumes -f || true
	echo "✅ All Docker resources cleaned up."

run test:
		docker compose -f tests/docker-compose.integration.yaml up -d