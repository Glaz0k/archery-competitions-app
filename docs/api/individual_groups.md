# Контроллер для управления индивидуальными группами

## Управление группой

### `GET /individual_groups/{group_id}`

> Получить информацию о группе. Участникам доступны только те, в которых они участвуют

_Уровень доступа:_ `[admin, user]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`individual_group.json`](../models/individual_group.md)
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)

### `DELETE /individual_groups/{group_id}`

> [!WARNING]
> Безвозвратно удалить существующую группу (если такой группы нет - ничего не происходит). Является каскадным удалением и должно помимо самой группы удалить все сопутсвующие с ней ресурсы, такие как квалификации, финальная сетка, информация об участии и все результаты вместе с выстрелами

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно удалено**\
  `204 No content`

## Управление участниками группы

### `GET /individual_groups/{group_id}/competitors`

> Получить список участников в группе. Участник имеет доступ только если состоит в группе

_Уровень доступа:_ `[admin, user]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `200 Ok`\
  `[`[`competitor_group_detail.json`](../models/competitor_group_detail.md)`]`
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)

### `POST /individual_groups/{group_id}/competitors/sync`

> Обновить список участников, синхронизируя его вместе со списком зарегистрированных на соревнование. Применимо только к группам на этапе создания

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `200 Ok`\
  `[`[`competitor_group_detail.json`](../models/competitor_group_detail.md)`]`
- [**Невозможно обновить список группы не на этапе создания**](../policies/user_errors.md/#невозможно-выполнить-действие)
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)

## Управление квалицикацией (на уровне группы)

### `GET /individual_groups/{group_id}/qualification`

> Получить текущую таблицу квалификаций для группы. Участник имеет доступ только если принадлежит к группе

_Уровень доступа:_ `[admin, user]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`qualification_table.json`](../models/qualification_table.md)
- [**Группа или квалификация не найдена**](../policies/user_errors.md/#не-найдено)

### `POST /individual_groups/{group_id}/qualification/start`

> Начать квалификацию в группе с заданными параметрами (при применении к начатой - только получает таблицу).\
> Участники назначаются в группу из числа зарегистрированных на соревнование и имеющих активный статус на момент начала квалификации по соответсвующим параметрам.\
> Квалификацию можно начать, если количество зарегистрированных участников больше 3.\
> Для каждого участника создается заданное количество раундов и серий и выбирается тип серий в зависимости от вида лука:
>
> - 6-10, X – для дивизионов Классический и Блочный лук;
> - 1-10, X – для остальных дивизионов
>
> Для всех участников становится активной первая серия первого раунда

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Тело запроса:_

```json
{
  "distance": "18m",
  "round_count": 2,
  "ranges_count": 10,
  "range_size": 3
}
```

_Ответы:_

- **Успешно**\
  `201 Created | 200 Ok`\
  [`qualification_table.json`](../models/qualification_table.md)
- [**Неверные параметры квалификации**](../policies/user_errors.md/#неверные-параметры)
- [**Недостаточно участников для начала**](../policies/user_errors.md/#невозможно-выполнить-действие)
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)

### `POST /individual_groups/{group_id}/qualification/end`

> Закончить квалификации в группу, по полученным очкам распределяет участников по местам (при применении к завершенной - только получает таблицу). Необходимо, чтобы все участники завершили свои раунды.\
> Распределение по местам при одинаковом счёте необходимо уточнить (не уверен, по 9-кам, 10-кам в какой-то очередности)

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`qualification_table.json`](../models/qualification_table.md)
- [**Невозможно закончить квалификацию из-за незавершённых участниками раундов**](../policies/user_errors.md/#невозможно-выполнить-действие)
  ```json
  {
    "details": [{}]
  }
  ```
  - `details`: `[`[`competitor_shrinked.json`](../models/competitor.md/#shrinked)`]`
- [**Группа или квалификация не найдена**](../policies/user_errors.md/#не-найдено)

## Управление финальной сеткой (на уровне группы)

### `GET /individual_groups/{group_id}/final_grid`

> Получить финальную сетку группы. Участник имеет доступ только если принадлежит к группе

_Уровень доступа:_ `[admin, user]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `200 Ok`\
  [`final_grid.json`](../models/final_grid.md)
- [**Группа или сетка не найдена**](../policies/user_errors.md/#не-найдено)

### `POST /individual_groups/{group_id}/quarterfinal/start`

> Начать четвертьфинал в группе. Формирует первый уровень сетки в зависимости от занятых на этапе квалификации мест. Спарринги и места сверху вниз:
>
> - (1 - 8) | (5 - 4) | (3 - 6) | (7 - 2)
>
> При неполной квалификации (кол-во участников меньше 8) незанятые места остаются пустыми, а оставшийся участник автоматически становится победителем и проходит в следующий этап:
>
> - (<ins>1</ins> - N) | (5 - 4) | (<ins>3</ins> - N) | (N - <ins>2</ins>)
>
> Можно начать четвертьфинал только для группы с завершенной квалификацией (при применении к начатому четвертьфиналу и стадиям далее - только получает сетку).\
> Для участников, получивших автоматическую победу не создается серия, а место считается неактивным.

> [!IMPORTANT]
> Максимальное количество серий и тип счёта (общий или победные очки) зависит от типа группы:
>
> - Блочный лук - 5 серий по три выстрела в каждой серии;
> - Все остальные классы - 3-5 серий по три выстрела до 6 победных очков:
>   - выигрыш – 2 очка,
>   - ничья – по 1 очку каждому,
>   - проигрыш – 0 очков.
>
> В случае равного счёта по итогам 5 серий необходима перестрелка - 1 стрела, побеждает тот, чья стрела ближе к центру мишени (при одинаковом счёте, организатору необходимо выставить `priority`)

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `201 Created | 200 Ok`\
  [`final_grid.json`](../models/final_grid.md)
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)
- [**Невозможно начать четвертьфинал**](../policies/user_errors.md#невозможно-выполнить-действие)

### `POST /individual_groups/{group_id}/semifinal/start`

> Начать полуфинал в группе. Формирует второй уровень сетки в зависимости от исходов спаррингов 1-го уровня. Можно начать только если все спарринги окончены победой одного из участников (при применении к начатому полуфиналу и стадиям далее - только получает сетку)

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `201 Created | 200 Ok`\
  [`final_grid.json`](../models/final_grid.md)
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)
- [**Невозможно начать полуфинал**](../policies/user_errors.md#невозможно-выполнить-действие)

### `POST /individual_groups/{group_id}/final/start`

> Начать финал в группе. Формирует третий уровень сетки в зависимости от исходов спаррингов 2-го уровня. Можно начать только если все спарринги окончены победой одного из участников (при применении к начатому финалу - только получает сетку). В спарринг бронзы попадают проигравшие в полуфинале

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `201 Created | 200 Ok`\
  [`final_grid.json`](../models/final_grid.md)
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)
- [**Невозможно начать финал**](../policies/user_errors.md#невозможно-выполнить-действие)

### `POST /individual_groups/{group_id}/final/end`

> Закончить финал и перевести группу в состояние завершенной. Можно закончить только если все спарринги окончены победой одного из участников (при применении к законченому - только получает сетку).

_Уровень доступа:_ `[admin]`\
_Переменные пути:_

- `group_id` - id группы

_Ответы:_

- **Успешно**\
  `201 Created | 200 Ok`\
  [`final_grid.json`](../models/final_grid.md)
- [**Группа не найдена**](../policies/user_errors.md/#не-найдено)
- [**Невозможно закончить финал**](../policies/user_errors.md#невозможно-выполнить-действие)
