# Модель участника

## full

### Допустимые типы

```
<competitor_full>
{
  "id": <number>,
  "full_name": <string>,
  "birth_date": <date ISO 8601 | null>,
  "identity": <gender | null>,
  "bow": <bow_class | null>,
  "rank": <sports_rank | null>,
  "region": <string | null>,
  "federation": <string | null>,
  "club": <string | null>
}
```

_Перечисления:_

- [`gender`](../enums/gender.md)
- [`bow_class`](../enums/bow_class.md)
- [`sports_rank`](../enums/sports_rank.md)

### Пример

```json
{
  "id": 759265,
  "full_name": "Иванов Иван",
  "birth_date": "1999-08-24",
  "identity": "male",
  "bow": "classic",
  "rank": "first_class",
  "region": "Санкт-Петербург",
  "federation": null,
  "club": null
}
```

## shrinked

### Допустимые типы

```
<competitor_shrinked>
{
  "id": <number>,
  "full_name": <string>
}
```

### Пример

```json
{
  "id": 759265,
  "full_name": "Иванов Иван"
}
```
