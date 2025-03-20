## Запрос изменения серии

### Допустимые типы

```
<change_range>
{
  "range_number": <number>,
  "shots": <[ <shot> ] | null>
}
```

_Модели:_

- [`shot`](../models/shot.md)

### Пример

```json
{
  "range_number": 1,
  "shots": [
    {
      "shot_number": 1,
      "score": "7"
    },
    {
      "shot_number": 2,
      "score": "M"
    },
    {
      "shot_number": 3,
      "score": "X"
    }
  ]
}
```
