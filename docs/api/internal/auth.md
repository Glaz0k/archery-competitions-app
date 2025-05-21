# Внутренний контроллер аутентификации

## Получение токена

### `POST /auth/sign_up`

> Зарегистрировать пользователя в системе. Выдаётся внутрення роль `user`.

_Тело запроса:_

[`credentials.json`](../requests/credentials.md)

_Ответы:_

- **Успешно**\
  `200 Ok`
  ```
  {
    "auth_data": <auth_data>,
    "token": <string>,  
  }
  ```
  - `auth_data`: [`auth_data.json`](../../models/auth_data.md)
- [**Неверные параметры регистрации**](../policies/user_errors.md/#неверные-параметры)
- [**Пользователь с таким `login` уже зарегистрирован**](../policies/user_errors.md/#ресурс-уже-существует)

### `POST /auth/sign_in`

> Войти с учётной записью пользователя.

_Тело запроса:_

[`credentials.json`](../requests/credentials.md)

_Ответы:_

- **Успешно**\
  `200 Ok`
  ```
  {
    "auth_data": <auth_data>,
    "token": <string>,  
  }
  ```
  - `auth_data`: [`auth_data.json`](../../models/auth_data.md)
- [**Неверные параметры входа**](../policies/user_errors.md/#неверные-параметры)

## Генерация учётных записей

### `POST /auth/generate/admin`

> Сгенерировать учётную запись админа с помощью пароля суперпользователя

_Тело запроса:_

```
{
  "superuser_password": <string>
}
```

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
- [**Неверный пароль**](../policies/user_errors.md/#неверные-параметры)

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
