# CardStore service

[Главная](../README.md) / Makefile

<br>

## Запуск проекта

```
make run
```

---

## Создать бинарник

В корне проекта появится бинарник ./dist/cardstore

```
make build
```

---

## Поднять БД

Поднять postgres БД в docker контейнере.

**Логин / пароль:** `postgres` / `secret`

```
make db-up
```

---

## Отключить БД

```
make db-down
```

---

## Миграции

```
make migrate-up
```

## Добавление ключей

```
make add-keys
```

## Шифрование не зашифрованных карт и сессий, которые добавлялись ранее

```
make encrypt-unencrypted
```
