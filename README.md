# GophKeeper

GophKeeper представляет собой клиент-серверную систему, позволяющую пользователю надёжно и безопасно хранить логины, пароли, бинарные данные и прочую приватную информацию.

## Типы хранимой информации

* пары логин/пароль;
* произвольные текстовые данные;
* произвольные бинарные данные;
* данные банковских карт.

## Запуск

Чтобы запустить сервер необходимо сгенерировать ключи и развернуть docker-compose:

```bash
make certs
docker-compose up
```

Далее можно запускать клиент и пользоваться программой

### Для нового пользователя

1. Пользователь получает клиент под необходимую ему платформу.
2. Пользователь проходит процедуру первичной регистрации.
3. Пользователь добавляет в клиент новые данные.
4. Клиент синхронизирует данные с сервером.

#### Для существующего пользователя

1. Пользователь получает клиент под необходимую ему платформу.
2. Пользователь проходит процедуру аутентификации.
3. Клиент синхронизирует данные с сервером.
4. Пользователь запрашивает данные.
5. Клиент отображает данные для пользователя.
