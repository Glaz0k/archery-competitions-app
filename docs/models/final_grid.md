## Модель финальной сетки

### Допустимые типы

```
<final_grid>
{
  "group_id": <number>,
  "quarterfinal": <quarterfinal>,
  "semifinal": <semifinal | null>,
  "final": <final | null>
}

<quarterfinal>
{
  "sparring_1": <sparring>,
  "sparring_2": <sparring>,
  "sparring_3": <sparring>,
  "sparring_4": <sparring>
}

<semifinal>
{
  "sparring_5": <sparring>,
  "sparring_6": <sparring>
}

<final>
{
  "sparring_gold": <sparring>,
  "sparring_bronze": <sparring>
}
```

_Модели:_

- [`sparring`](sparring.md)

### Пример

```json
{
  "group_id": 324098,
  "quarterfinal": {
    "sparring_1": {
      "id": 22340981,
      "top_place": {
        "id": 3458761,
        "competitor": {
          "id": 759265,
          "full_name": "Иванов Иван"
        },
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
        },
        "is_active": true
      },
      "bot_place": null,
      "state": "top_win"
    },
    "sparring_2": null,
    "sparring_3": null,
    "sparring_4": null
  },
  "semifinal": null,
  "final": null
}
```
