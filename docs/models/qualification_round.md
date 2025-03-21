# Модель раунда квалификации

## full

### Допустимые типы

```
<qualification_round_full>
{
  "section_id": <number>,
  "round_ordinal": <number>,
  "is_active": <bool>,
  "range_group": <range_group>
}
```

_Модели:_

- [`range_group`](range_group.md)

### Пример

```json
{
  "section_id": 12034987,
  "round_ordinal": 1,
  "is_active": true,
  "range_group": {
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
}
```

## shrinked

### Допустимые типы

```
<qualification_round_shrinked>
{
  "round_ordinal": <number>,
  "is_active": <bool>,
  "total_score": <number | null>
}
```

### Пример

```json
{
  "round_ordinal": 2,
  "is_active": true,
  "total_score": 10
}
```
