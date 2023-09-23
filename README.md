# CardStore service

Перед запуском создайте в корне проекта файл ".env". Пример можно взять из ".env-example"

## Документация по сервису

* [Makefile](./docs/makefile.md)
* [API карты](./docs/api-cards.md)
* [API сессии](./docs/api-sessions.md)

## Линтер

**Автоформатирование кода:**

```bash
gofmt -s -w .
```

**Запуск линтера:**

Руководство по установке __golangci-lint__:

https://golangci-lint.run/usage/install/

```bash
golangci-lint run ./...
```
