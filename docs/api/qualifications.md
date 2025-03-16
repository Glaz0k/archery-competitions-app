# Контроллер для управления квалификациями

## Управление квалификацией

### `GET /qualifications/{group_id}`

> Получает текущую таблицу квалификаций для группы. Участник имеет доступ только если принадлежит к группе

_Уровень доступа:_ `[admin, user]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`qualification_table.json`](../models/qualification_table.md)
- [**Квалификация не найдена**](../policies/user_errors.md/#не-найдено)

### `GET /qualifications/{group_id}/rounds/{round_number}?competitor_id=`

> Получает информацию о раунде квалификации для отдельного участника. Участник имеет доступ, если принадлежит к группе

_Уровень доступа:_ `[admin, user]`\
_Переменные пути:_

- `group_id` - id группы
- `round_number` - номер раунда

_Параметры запроса:_

- `competitor_id [required]` - id участника

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`qualification_round_full.json`](../models/qualification_round.md/#full)
- [**Квалификация, раунд или участник не найдены**](../policies/user_errors.md/#не-найдено)
