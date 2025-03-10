# Приложение для организации соревнований по стрельбе из лука
## 1. Организация репозитория
### 1.1. Структура
```
archery-competitions-app/
│
├── backend/
│   ├── api-gateway/
│   │   ├── nginx.conf
│   │   └── Dockerfile
│   ├── redis-session-storage/
│   │   └── Dockerfile
│   ├── auth-server/
│   │   ├── src/
│   │   └── Dockerfile
│   ├── postgres-auth/
│   │   └── Dockerfile
│   ├── app-server/
│   │   ├── src/
│   │   └── Dockerfile
│   ├── postgres-app/
│   │   └── Dockerfile
│   └── docker-compose.yml
│
├── frontend/
│   ├── web-app/
│   │   ├── src/
│   │   └── Dockerfile
│   ├── mobile-app/
│   └── docker-compose.yml
│
├── docs/api/
│
├── README.md
└── .gitignore
```
* backend/:
  * api-gateway/: Содержит конфигурацию Nginx и Dockerfile для создания образа API Gateway.
  * redis-session-storage/: Конфигурация Docker для Redis.
  * auth-server/: Java проект для сервера аутентификации.
  * postgres-auth/: Конфигурация Docker для базы данных Postgres, используемой сервером аутентификации.
  * app-server/: Go проект для монолитного сервиса, реализующего бизнес-логику.
  * postgres-app/: Конфигурация Docker для базы данных Postgres, используемой приложением.
  * docker-compose.yml: Файл для оркестрации всех бэкенд сервисов.
* frontend/
  * web-app/: React проект для SPA админ-приложения.
  * mobile-app/: Flutter проект для мобильного приложения.
  * docker-compose.yml: Файл для оркестрации фронтенд приложений.
* docs/api/: Документация для REST API
* README.md: Основная документация проекта.
* .gitignore: Файл для исключения ненужных файлов из Git, может содержаться в каждой из поддиректорий.

### 1.2. Ветки
- **master** — основная ветка, содержащая стабильную версию кода. Прямые коммиты в эту ветку запрещены.
- **develop** — ветка для разработки. Все изменения должны быть слиты сюда перед тем, как попасть в `master`.
- **feature/** — ветки для разработки новых функций. Название ветки должно соответствовать формату `feature/<название-функционала>`.
- **bugfix/** — ветки для исправления багов. Название ветки должно соответствовать формату `bugfix/<описание-бага>`.

Путь кода в `master` ветку:
```
master
└─< develop
    ├─< feature/*
    └─< bugfix/*
```

## 2. Процедура добавления кода
1. Перед началом работы над новой функцией или исправлением создайте новую ветку от `develop`:
2. После завершения работы над веткой создайте Pull Request в ветку `develop`.
3. Название PR должно соответствовать формату:
   `<тип>: <описание>`
   Где `<тип>` может быть:
   * `new` — новая функция
   * `fix` — исправление бага
   * `docs` — изменения в документации
   * `refactor` — рефакторинг кода
   * `test` — добавление или изменение тестов
   * `misc` — вспомогательные изменения (например, обновление зависимостей)

> [!IMPORTANT]
> Минимальный уровень покрытия кода backend-сервисов unit-тестами — 60%.

## 3. Документация
1. Докуметнация внешнего REST API должна находится в `docs\`, оформленная в виде Markdown-файлов с названием `<название-сервиса>.md`
2. Модели для документации должны лежать в `docs\models\`, в виде `<название-модели>.json`
3. Документация каждого сервиса должна содержать:
   * Функциональное название endpoint-а
   * Уровень доступа, ограничения
   * HTTP-метод, URI запроса (возможны пояснения значений переменных пути)
   * Описание необходимых заголовков (при наличии)
   * Тело запроса в виде JSON (при наличии), также может представлять из себя модель из `docs\models\`.
   * Возможные варианты ответа сервиса, каждый из которых включает:
     - HTTP-код ответа
     - Тело ответа (при наличии)

### Пример:
...\
/## 4. Регистрация участника на соревнование\
/### `POST /competitions/{id}/register`\
Уровень доступа: `[organizer]`\
Тело запроса:
```json
{
  "competitor_id": 724560469214
}
```
Ответы:
  * **Успешно зарегистрирован или уже был зарегистрирован**\
    201 Created/200 Ok
    ```json
    {
      "id": 9275023853
      "competition_id": 812031285712
      "competitor_id": 724560469214
      "created_at": "2025-03-10 22:12:44.126597+03"
    }
    ```
  * **Неподходящий уровень доступа**\
    403 Forbidden
    ```json
    {
      "error": "Insufficient access level"
    }
    ```
  * **Соревнования с таким id не существует**\
    404 Not found
    ```json
    {
      "error": "Competition not found"
    }
    ```
  * **Участника с таким id не существует**\
    404 Not found
    ```json
    {
      "error": "Competitor not found"
    }
    ```
  * **Внутренняя ошибка сервера**\
    500 Interanal Server Error\
    `No content`
    
...

## 4. Зоны ответсвенности
* Project Manager — Лебедев Антон
* Backend:
  - app-server — Пиявкин Антон, Козакова Анна, Дудкина София
  - auth-server — Кравченко Никита
  - api-gateway — Кравченко Никита
* Frontend:
  - web — Кравченко Никита
  - mobile — Демиденко Никита, Новохацкий Данил
* CI/CD — Лебедев Антон, Пиявкин Антон

## 5. Полезные материалы
* Графические материалы — https://drive.google.com/drive/folders/1xE3skqLdafwkhKQWRsG8dghXoY-ITycv?usp=sharing
* Схема базы данных основного приложения — https://dbdiagram.io/d/Bow-Competitions-67c6ce5d263d6cf9a0279758
