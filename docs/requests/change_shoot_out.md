## Запрос изменения перестрелки

### Допустимые типы

```
<change_shoot_out>
{
  "score": <string>,
  "priority": <bool | null>
}
```

_Специальные значения:_

- `score`: `^(M|[1-9]|10|X)$`

### Пример

```json
{
  "score": 8,
  "priority": true
}
```
