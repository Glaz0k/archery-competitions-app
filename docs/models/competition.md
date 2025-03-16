## Модель соревнования

### Допустимые типы

```
<competition>
{
  "id": <number>,
  "stage": <competition_stage>,
  "start_date": <date ISO 8601 | null>,
  "end_date": <date ISO 8601 | null>,
  "is_ended": <bool>
}
```

_Перечисления:_

- [`competition_stage`](../enums/competition_stage.md)

### Пример

```json
{
  "id": 103,
  "stage": "III",
  "start_date": "2025-01-18",
  "end_date": "2025-01-19",
  "is_ended": false
}
```
