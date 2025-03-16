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

<section>
{
  "id": <number>,
  "competitor": <competitor_shrinked>,
  "place": <number | null>,
  "rounds": [ <qualification_round_shrinked> ],
  "total": <number | null>,
  "10_s": <number | null>,
  "9_s": <number | null>,
  "rank_gained": <sports_rank | null>
}
```

_Модели:_

- [`competitor_shrinked`](competitor.md/#shrinked)
- [`qualification_round_shrinked`](qualification_round.md/#shrinked)

_Перечисления:_

- [`sports_rank`](../enums/sports_rank.md)

### Пример

```json
{
  "group_id": 824,
  "distance": "18m",
  "round_count": 2,
  "sections": [
    {
      "id": 3052,
      "competitor": {},
      "place": null,
      "rounds": [
        {
          "round_number": 1,
          "is_ongoing": true,
          "total": 28
        },
        {
          "round_number": 2,
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
