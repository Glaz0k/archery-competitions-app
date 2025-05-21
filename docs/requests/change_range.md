## Запрос изменения серии

### Допустимые типы

```
<change_range>
{
  "range_ordinal": <number>,
  "shots": <[ <shot> ] | null>
}
```

_Модели:_

- [`shot`](../models/shot.md)

### Пример

```json
{
  "range_ordinal": 1,
  "shots": [
    {
      "shot_ordinal": 1,
      "score": "7"
    },
    {
      "shot_ordinal": 2,
      "score": "M"
    },
    {
      "shot_ordinal": 3,
      "score": "X"
    }
  ]
}
```
