## Запрос учётной записи

### Допустимые типы

```
<credentials>
{
  "login": <string>,
  "password": <string>
}
```

_Специальные значения:_

- `login`: `^[a-zA-Z0-9._-]{6,20}$`
- `password`: `^[a-zA-Z0-9._-]{6,20}$`

### Пример

```json
{
  "login": "ivanov.i.i",
  "password": "qwerty123"
}
```
