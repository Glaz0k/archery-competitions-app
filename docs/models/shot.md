## Модель выстрела

### Допустимые типы

```
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
  "shot_number": 1,
  "score": "10"
}
```
