# Модель раунда квалификации

## full

### Допустимые типы

```
<qualification_round_full>
{
  "section_id": <number>,
  "round_number": <number>,
  "is_ongoing": <bool>
  "range_group": <range_group>
}

<range_group>
{
  "id": <number>,
  "ranges_max_count": <number>,
  "range_size": <number>,
  "ranges": <[ <range> ]>,
  "total_score": <number>
}

<range>
{
  "id": <number>,
  "range_number": <number>,
  "is_ongoing": <bool>,
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
  "section_id": 12034987,
  "round_number": 1,
  "is_ongoing": true,
  "range_group": {
    "id": 1297461324,
    "ranges_max_count": 10,
    "range_size": 3,
    "ranges": [
      {
        "id": 234098123,
        "range_number": 1,
        "is_ongoing": false,
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
      },
      {
        "id": 23987614,
        "range_number": 2,
        "is_ongoing": true,
        "shots": [
          {
            "shot_number": 1,
            "score": "8"
          },
          {
            "shot_number": 2,
            "score": null
          },
          {
            "shot_number": 3,
            "score": null
          }
        ],
        "range_score": 8
      },
      {
        "id": 43563568,
        "range_number": 3,
        "is_ongoing": false,
        "shots": null,
        "range_score": null
      }
    ],
    "total_score": 37
  }
}
```

## shrinked

### Допустимые типы

```
<qualification_round_shrinked>
{
  "round_number": <number>,
  "is_ongoing": <bool>,
  "total_score": <number | null>
}
```

### Пример

```json
{
  "round_number": 2,
  "is_ongoing": true,
  "total_score": 10
}
```
