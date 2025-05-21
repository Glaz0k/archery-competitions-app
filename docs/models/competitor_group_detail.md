## Модель записи участника в группе

### Допустимые типы

```
<competitor_competition_detail>
{
  "group_id": <number>,
  "competitor": <competitor_full>
}
```

_Модели:_

- [`competitor_full`](competitor.md#full)

### Пример

```json
{
  "group_id": 103,
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
  }
}
```
