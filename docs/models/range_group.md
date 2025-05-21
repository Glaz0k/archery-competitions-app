## Модель группы серий

### Допустимые типы

```
<range_group>
{
  "id": <number>,
  "ranges_max_count": <number>,
  "range_size": <number>,
  "type": <range_type>
  "ranges": <[ <range> ]>,
  "total_score": <number | null>
}
```

_Модели:_

- [`range`](range.md)

_Перечисления:_

- [`range_type`](../enums/range_type.md)

### Пример

```json
{
  "id": 1297461324,
  "ranges_max_count": 10,
  "range_size": 3,
  "ranges": [
    {
      "id": 234098123,
      "range_ordinal": 1,
      "is_active": false,
      "shots": [
        {
          "shot_ordinal": 1,
          "score": "10"
        },
        {
          "shot_ordinal": 2,
          "score": "9"
        },
        {
          "shot_ordinal": 3,
          "score": "X"
        }
      ],
      "range_score": 29
    },
    {
      "id": 23987614,
      "range_ordinal": 2,
      "is_active": true,
      "shots": [
        {
          "shot_ordinal": 1,
          "score": "8"
        },
        {
          "shot_ordinal": 2,
          "score": null
        },
        {
          "shot_ordinal": 3,
          "score": null
        }
      ],
      "range_score": 8
    },
    {
      "id": 43563568,
      "range_ordinal": 3,
      "is_active": false,
      "shots": null,
      "range_score": null
    }
  ],
  "total_score": 37
}
```
