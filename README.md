# link shortener

[![CI](https://github.com/kun1ts4/link-shortener/actions/workflows/ci.yml/badge.svg)](https://github.com/kun1ts4/link-shortener/actions/workflows/ci.yml)

## Запуск

1. С хранилищем PostgreSQL:

```bash
make run-postgres
```

2. С хранилищем в памяти:

```bash
make run-memory
```

## API

- `POST /` — создание короткой ссылки. Тело запроса: JSON (`Content-Type: application/json`) с полем `url`. Ответ: `201` и JSON в теле с полями `short` и `original`.
- `GET /{short}.json` — получение информации о короткой ссылке в формате JSON (только JSON). Ответы:
  - 200: {"url":"<original>", "clicks":<number>} при успехе
  - 404: ссылка не найдена
  - 400: неверный формат запроса (например, если использовать `GET /{short}` без суффикса `.json`)

## Алгоритм сокращения

1. Вычисляется хэш оригинальной ссылки (SHA-256)
2. Хэш переводится в Hex-строку
3. Hex-строка делится на двузначные 16-ричные числа из них берутся первые N (длина ссылки в задании = 10)
4. числа переводятся в десятичные и конвертируются в символы алфавита (0-9, a-z, A-Z, _) по остатку от деления на длину алфавита
5. Получившаяся строка и есть короткая ссылка

## Конфигурация:

Настройки приложения (длина ссылок, алфавит, таймауты, другие параметры):

`config/config.yaml`

Секреты базы (пример, необходимо создать .env):

`.env.example`

## Примеры использования

Создать короткую ссылку:

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  --data '{"url":"https://example.com/very/long/url"}' \
  http://localhost:8080/
# -> в ответе будет JSON, например: {"short":"jonidmKLIk","original":"https://example.com/very/long/url"}
```

Получить данные короткой ссылки (JSON):

```bash
curl http://localhost:8080/jonidmKLIk.json
# -> {"url":"https://example.com/very/long/url","clicks":0}
```

Если вызвать `GET /jonidmKLIk` (без `.json`), сервис вернёт 400 и сообщение: "use .json format to get link"
