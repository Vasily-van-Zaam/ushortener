# go-musthave-shortener-tpl

Шаблон репозитория для практического трек "Веб-разработка на Go"

# Начало работы

1. Для старта сервера запустить `go run cmd/shortener/main.go`
2. Для сборки проекта в bin файл `go build cmd/shortener/main.go`
3. Генерация документации swagger `swag init -g cmd/shortener/main.go`

# Документация API

- Для проверки API запросов откройте в браузере http://localhost:8080/swagger-docs/

go tool pprof -http=":9090" -seconds=30 http://localhost:8888/debug/pprof/heap

go pprof -top -diff_base=profiles/base.pprof profiles/result.pprof
