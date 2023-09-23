# CardStore service

[Главная](../README.md) / API по сессии

<br>

## Добавить сессию

**URL: `[POST] /v1/sessions`**

Тело запроса отправляется в формате JSON:

```
{
    "key": "session-key",
    "body": "{\"key1\":\"data1\"}"
}
```

---

## Получить данные сессии по ключу

**URL: `[GET] /v1/sessions/{sessionKey}`**

**Ответ:**

```json
{
    "id": 3,
    "key": "secret-key",
    "body": "{\"key1\":\"value1\"}",
    "expire_at": "2022-04-21T13:53:39.055583+05:00"
}
```
