## Модель индивидуальной группы

### Допустимые типы

```
<individual_group>
{
  "id": <number>,
  "competition_id": <number>,
  "bow": <bow_class>,
  "identity": <gender | null>,
  "state": <group_state>
}
```

_Перечисления:_

- [`bow_class`](../enums/bow_class.md)
- [`gender`](../enums/gender.md)
- [`group_state`](../enums/group_state.md)

_Пояснения:_

- null в поле `identity` означает "Объединённый"

### Пример

```json
{
  "id": 812,
  "competition_id": 103,
  "bow": "classic",
  "identity": "male",
  "state": "created"
}
```
