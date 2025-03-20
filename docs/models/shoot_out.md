## Модель перестрелки

### Допустимые типы

```
<shoot_out>
{
  "id": <number>,
  "score": <string>,
  "priority": <bool | null>
}
```

_Специальные значения:_

- `score`: `^(M|[1-9]|10|X)$`

### Пример

```json
{
  "id": 234098123,
  "score": "8",
  "priority": true
}
```
