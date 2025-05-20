## 1. Регистрация участника на соревнование

### `POST /competitions/{competition_id}/register`

_Уровень доступа:_ `[organizer]`\
_Тело запроса:_

```json
{
  "competitor_id": 724560469214
}
```

_Ответы:_

- **Успешно зарегистрирован или уже был зарегистрирован**\
  `201 Created / 200 Ok`
  ```json
  {
    "id": 9275023853,
    "competition_id": 812031285712
    "competitor_id": 724560469214
    "created_at": "2025-03-10 22:12:44.126597+03"
  }
  ```
- **Неподходящий уровень доступа**\
  `403 Forbidden`
  ```json
  {
    "error": "Insufficient access level"
  }
  ```
- **Соревнования с таким id не существует**\
  `404 Not found`
  ```json
  {
    "error": "Competition not found"
  }
  ```
- **Участника с таким id не существует**\
  `404 Not found`
  ```json
  {
    "error": "Competitor not found"
  }
  ```
- **Внутренняя ошибка сервера**\
  `500 Interanal Server Error`\
  `No content`
