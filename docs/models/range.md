## Модель серии

### Допустимые типы

```
<range>
{
  "id": <number>,
  "range_ordinal": <number>,
  "is_active": <bool>,
  "shots": <[ <shot> ] | null>,
  "range_score": <number | null>
}
```

_Модели:_

- [`shot`](shot.md)

### Пример

```json
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
}
```
