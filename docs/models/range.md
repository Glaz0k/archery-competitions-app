## Модель серии

### Допустимые типы

```
<range>
{
  "id": <number>,
  "range_number": <number>,
  "is_active": <bool>,
  "shots": <[ <shot> ] | null>,
  "range_score": <number | null>
}

<shot>
{
  "shot_number": <number>,
  "score": <string | null>
}
```

_Специальные значения:_

- `score`: `^(M|[1-9]|10|X)$`

### Пример

```json
{
  "id": 234098123,
  "range_number": 1,
  "is_active": false,
  "shots": [
    {
      "shot_number": 1,
      "score": "10"
    },
    {
      "shot_number": 2,
      "score": "9"
    },
    {
      "shot_number": 3,
      "score": "X"
    }
  ],
  "range_score": 29
}
```
