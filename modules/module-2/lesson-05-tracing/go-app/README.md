# Demo Go Application with OpenTelemetry Tracing

Демонстрационное приложение на Go с распределенной трассировкой.

## Описание

Приложение симулирует обработку заказов с несколькими этапами:
- Валидация заказа
- Проверка наличия на складе
- Обработка платежа
- Отправка заказа

Каждый этап создает span с метриками и атрибутами.

## Запуск

### Локально

```bash
go mod download
go run main.go
```

### Docker

```bash
docker build -t demo-go-app .
docker run -p 8080:8080 demo-go-app
```

## Тестирование

```bash
curl -X POST http://localhost:8080/order
```

## Endpoints

- `POST /order` - создать новый заказ
- `GET /health` - проверка здоровья приложения
