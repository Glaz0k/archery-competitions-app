## Модель перестрелки

### Допустимые типы

```
<shoot_out>
{
  "id": <number>,
  "score": <string | null>,
  "priority": <bool | null>
}
```

_Специальные значения:_

- `score`: `^(M|[1-9]|10|X)$`

### Пример

Только созданная:
```json
{
  "id": 239871243,
  "score": null,
  "priority": null
}
```

Заполненная:
```json
{
  "id": 234098123,
  "score": "8",
  "priority": true
}
```
