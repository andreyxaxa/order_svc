# Демонстрационный сервис, отображающий данные о заказе. 
Сервис получает данные заказов из очереди (Kafka), сохраняет их в базу данных (PostgreSQL) и кэширует в памяти для быстрого доступа.

[Быстрый старт](https://github.com/andreyxaxa/order_svc?tab=readme-ov-file#%D0%B7%D0%B0%D0%BF%D1%83%D1%81%D0%BA)

## Обзор

- Документация API - Swagger - http://localhost:8080/swagger
- Метрики - Prometheus metrics - http://localhost:8080/metrics
- Конфиг - [config/config.go](https://github.com/andreyxaxa/order_svc/blob/main/config/config.go). Читается из `.env` файла.
- Логгер - [pkg/logger/logger.go](https://github.com/andreyxaxa/order_svc/blob/main/pkg/logger/logger.go). Интерфейс позволяет подменить логгер.
- Graceful shutdown - [internal/app/app.go](https://github.com/andreyxaxa/order_svc/blob/main/internal/app/app.go).
- Кеширование данных (LRU-cache + TTL) - [internal/repo/cache/orders_cache.go](https://github.com/andreyxaxa/order_svc/blob/main/internal/repo/cache/orders_cache.go). Сама реализация кеша - [internal/repo/cache/lru/lru.go](https://github.com/andreyxaxa/order_svc/blob/main/internal/repo/cache/lru/lru.go).
- Восстановление кеша из БД при старте - [internal/app/app.go](https://github.com/andreyxaxa/order_svc/blob/main/internal/app/app.go). В `.env` укажем, сколько последних заказов по дате хотим восстановить, например 100 - `CACHE_PRELOAD_LIMIT=100`. Читаем из четырех таблиц, поэтому используем уровень изоляции repeatable read и транзакции.
- Удобная и гибкая конфигурация HTTP сервера - [pkg/httpserver/options.go](https://github.com/andreyxaxa/order_svc/blob/main/pkg/httpserver/options.go).
  Позволяет конфигурировать сервер в конструкторе таким образом:
  ```go
  httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port))
  ```
  Аналогичный подход с таким конфигурированием - [pkg/kafka/options.go](https://github.com/andreyxaxa/order_svc/blob/main/pkg/kafka/options.go), [pkg/postgres/options.go](https://github.com/andreyxaxa/order_svc/blob/main/pkg/postgres/options.go).
- В слое хэндлеров применяется версионирование - [internal/controller/http/v1](https://github.com/andreyxaxa/order_svc/tree/main/internal/controller/http/v1).
  Для версии v2 нужно будет просто добавить папку `http/v2` с таким же содержимым, в файле [internal/controller/http/router.go](https://github.com/andreyxaxa/order_svc/blob/main/internal/controller/http/router.go) добавить строку:
  ```go
  {
      v1.NewOrderRoutes(apiV1Group, o, l)
  }

  {
      v2.NewOrderRoutes(apiV1Group, o, l)
  }
  ```
- Можно получить данные заказа как в виде JSON - `v1/order/info?order_uid=...`, так и воспользоваться веб-интерфейсом - `v1/order/info/html`.

## Запуск

Клонируем репозиторий, выполняем:
```
make compose-up
```

## Тесты

Запустить тесты:
```
make test
```

Убедиться, что кеш ускоряет получение данных:
```
make bench-test
```
<img width="1036" height="165" alt="image" src="https://github.com/user-attachments/assets/ac2dc9da-d7f3-48d0-975f-0844bbd50bb7" />

## Прочие `make` команды
Зависимости:
```
make deps
```
docker compose down:
```
make compose-down
```

## API

#### JSON:
- http://localhost:8080/v1/order/info?order_uid=b563feb7b2b84b6test - вернет JSON со всей информацией о заказе с id `b563feb7b2b84b6test`;

#### HTML:
- http://localhost:8080/v1/order/info/html - форма для ввода id заказа, кнопка поиска;
- http://localhost:8080/v1/order/info/html?order_uid=b563feb7b2b84b6test - вернет HTML со всей информацией о заказе с id `b563feb7b2b84b6test`;
