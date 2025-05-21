## Модель таблицы квалификаций

### Допустимые типы

```
<qualification_table>
{
  "group_id": <number>,
  "distance": <string>,
  "round_count": <number>,
  "sections": [ <section> ]
}
```

_Модели:_

- [`section`](section.md)

### Пример

```json
{
  "group_id": 824,
  "distance": "18m",
  "round_count": 2,
  "sections": [
    {
      "id": 3052,
      "competitor": {
        "id": 759265,
        "full_name": "Иванов Иван"
      },
      "place": null,
      "rounds": [
        {
          "round_ordinal": 1,
          "is_ongoing": true,
          "total": 28
        },
        {
          "round_ordinal": 2,
          "is_ongoing": false,
          "total": 0
        }
      ],
      "total": 28,
      "10_s": 1,
      "9_s": 1,
      "rank_gained": null
    }
  ]
}
```
