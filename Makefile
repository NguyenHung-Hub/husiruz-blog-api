r:
	go run main.go
d:
	net start mongodb
redis:
	docker run --name blog_rdb -p 6377:6379 -d redis:6.2-alpine
rbash:
	docker exec -it blog_rdb /bin/ash
a:
	air -d
	
.PHONY: r d redis rbash a