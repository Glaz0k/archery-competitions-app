# `Обновления`

## [update] CreateCompetition

#### method: POST

- Исправлен эндпоинт
- Добавлена проверка на существоание

#### endpoint: /api/cups/{cup_id}/competitions

Тело запроса:

```json
{
  "stage": "I",
  "start_date": "2023-10-01T00:00:00Z",
  "end_date": "2023-10-10T00:00:00Z",
  "is_ended": false
}
```

## [new] Edit Competition

- Возможность изменения даты соревнования

#### endpoint: /api/competitions/{competition_id}

#### method: PUT

Тело запроса:

```json
{
  "start_date": "2023-10-01T00:00:00Z",
  "end_date": "2023-10-10T00:00:00Z"
}
```

Ответ:

```json
{
  "id": 123,
  "cup_id": 456,
  "stage": "Quarterfinals",
  "start_date": "2023-10-15T00:00:00Z",
  "end_date": "2023-10-25T00:00:00Z",
  "is_ended": false
}
```

## [new] GetCup

- Получение информации о кубке

#### endpoint: /api/сups/{cup_id}

#### method: GET

Ответ:

```json
{
  "id": 2,
  "title": "Test Cup",
  "address": "Test Address",
  "season": "2023/2024"
}
```

## [new] CreateCup

- Создание кубка

#### endpoint: /api/cups

#### method: POST

Тело запроса:

```json
{
  "title": "Test Cup",
  "address": "Test Address",
  "season": "2023/2024"
}
```

## [update] CreateIndividualGroup

- Исправлен эндпоинт

#### endpoint: /api/cups/{cup_id}/individual-groups

#### method: POST

Тело запроса:

```json
{
  "bow": "classic",
  "identity": "male",
  "state": "created"
}
```

## `TODO:`

### - Рефакторинг Qualification, Ranges, Shots