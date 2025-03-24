# Внутренний контроллер аутентификации

## Получение токена

### `POST /auth/registration`

> Зарегистрировать пользователя в системе. Выдаётся внутрення роль `user`.

_Тело запроса:_

[`credentials.json`](../requests/credentials.md)

_Ответы:_

- **Успешно**\
  `200 Ok`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "type": "Bearer"
  }
  ```
- [**Неверные параметры регистрации**](../policies/user_errors.md/#неверные-параметры)
- [**Пользователь с таким `login` уже зарегистрирован**](../policies/user_errors.md/#ресурс-уже-существует)

### `POST /auth/login`

> Войти с учётной записью пользователя.

_Тело запроса:_

[`credentials.json`](../requests/credentials.md)

_Ответы:_

- **Успешно**\
  `200 Ok`
  ```json
  {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "type": "Bearer"
  }
  ```
- [**Неверные параметры входа**](../policies/user_errors.md/#неверные-параметры)
