# Контроллер аутентификации

## Регистрация

### `POST /auth/sign_up`

> Зарегистрировать пользователя в системе. Выдаётся внутрення роль `user`. После регистрации возвращается сессионная cookie и данные аутентификации

_Тело запроса:_

[`credentials.json`](../requests/credentials.md)

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`auth_data.json`](../models/auth_data.md)
- [**Неверные параметры регистрации**](../policies/user_errors.md/#неверные-параметры)
- [**Пользователь с таким `login` уже зарегистрирован**](../policies/user_errors.md/#ресурс-уже-существует)

## Управление сессией

### `POST /auth/sign_in`

> Войти с учётной записью пользователя. После регистрации возвращается сессионная cookie и данные аутентификации

_Тело запроса:_

[`credentials.json`](../requests/credentials.md)

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`auth_data.json`](../models/auth_data.md)
- [**Неверные параметры входа**](../policies/user_errors.md/#неверные-параметры)

### `POST /auth/logout`

> Выйти из учётной записи пользователя. После выхода удаляется связанный с cookie сессионный токен и для дальнейшей работы необходимо заново войти

_Ответы:_

- **Успешно**\
  `204 No content`

## Генерация учётных записей

### `POST /auth/generate/user`

> Сгенерировать учётную запись пользователя

_Уровень доступа:_ `[admin]`

_Ответы:_

- **Успешно**\
  `200 Ok`
  ```
  {
    "auth_data": <auth_data>,
    "credentials": <credentials>,  
  }
  ```
  - `auth_data`: [`auth_data.json`](../../models/auth_data.md)
  - `credentials`: [`credentials.json`](../../requests/credentials.md)
