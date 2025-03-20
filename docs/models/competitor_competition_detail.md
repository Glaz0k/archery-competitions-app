## Модель записи участника на соревнование

### Допустимые типы

```
<competitor_competition_detail>
{
  "competition_id": <number>,
  "competitor": <competitor_full>,
  "is_active": <bool>,
  "created_at": <YYYY-MM-DDThh:mm:ss±hh ISO 8601>
}
```

_Модели:_

- [`competitor_full`](competitor.md#full)

### Пример

```json
{
  "competition_id": 103,
  "competitor": {
    "id": 759265,
    "full_name": "Иванов Иван",
    "birth_date": "1999-08-24",
    "identity": "male",
    "bow": "classic",
    "rank": "first_class",
    "region": "Санкт-Петербург",
    "federation": null,
    "club": null
  },
  "is_active": true,
  "created_at": "2005-08-09T18:31:42+03"
}
```
