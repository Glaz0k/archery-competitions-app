## Реализация JWT авторизации

### 1. Вид токена

1. Header:

```json
{
  "alg": "RS256",
  "typ": "JWT"
}
```

2. Payload:

```json
{
  "user_id": 934591243,
  "role": "organizer"
}
```

3. Signature

### 2. Допустимые роли

- `organizer` - организатор
- `competitor` - участник

### 3. Способ передачи:

Находится в заголовке `Authorization`:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

### 4. Обработка ошибок авторизации

- **Любой неаутентифицированный запрос**\
  `401 Unauthorized`\
  `No content`

- **Неподходящая роль**\
  `403 Forbidden`\
  `No content`
