## Модель места в спарринге

### Допустимые типы

```
<sparring_place>
{
  "id": <number>,
  "competitor": <competitor_shrinked>,
  "range_group": <range_group>,
  "is_active": <bool>,
  "shoot_out": <shoot_out | null>,
  "sparring_score": <number>
}
```

_Модели:_

- [`competitor_shrinked`](competitor.md#shrinked)
- [`range_group`](range_group.md)
- [`shoot_out`](shoot_out.md)

### Пример

```json
{
  "id": 3458761,
  "competitor": {
    "id": 759265,
    "full_name": "Иванов Иван"
  },
  "range_group": {
    "id": 1297461324,
    "ranges_max_count": 5,
    "range_size": 3,
    "ranges": [
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
      },
      {
        "id": 23987614,
        "range_number": 2,
        "is_active": true,
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
        "is_active": false,
        "shots": null,
        "range_score": null
      }
    ],
    "total_score": 37
  },
  "is_active": true,
  "shoot_out": null,
  "sparring_score": 37
}
```
