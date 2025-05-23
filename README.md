# Приложение для организации соревнований по стрельбе из лука

## 1. Организация репозитория

### 1.1. Структура

```
archery-competitions-app/
│
├── backend/
│   ├── api-gateway/
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   └── ...
│   ├── app-server/
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   └── ...
│   └── auth-server/
│       ├── Dockerfile
│       ├── docker-compose.yml
│       └── ...
│
├── frontend/
│   ├── web-app/
│   │   ├── Dockerfile
│   │   ├── docker-compose.yml
│   │   └── ...
│   └── mobile-app/
│
├── docs/
│   ├── api/
│   ├── enums/
│   ├── models/
│   ├── policies/
│   └── requests/
│
├── README.md
└── LICENSE
```

- backend/:
  - api-gateway/: Содержит Node.js сервер в качестве API Gateway.
  - app-server/: Go проект для монолитного сервиса, реализующего бизнес-логику.
  - auth-server/: Java проект для сервера аутентификации.
- frontend/
  - web-app/: React проект для SPA админ-приложения.
  - mobile-app/: Flutter проект для мобильного приложения.
- docs/
  - api/: Документация для REST API
  - enums/: Используемые перечисления
  - models/: Используемые модели
  - policies/: Используемые при работе с REST API политики
  - requests/: Используемые ответы API
- README.md: Основная документация проекта.
- LICENSE: Файл лицензии.

### 1.2. Ветки

- **master** — основная ветка, содержащая стабильную версию кода. Прямые коммиты в эту ветку запрещены.
- **develop** — ветка для разработки. Все изменения должны быть слиты сюда перед тем, как попасть в `master`.
- **feature/** — ветки для разработки новых функций. Название ветки должно соответствовать формату `feature/<название-функционала>`.

### 1.3. Процедура добавления кода

1. Перед началом работы над новой функцией или исправлением создайте новую ветку от `develop`:
2. После завершения работы над веткой создайте Pull Request в ветку `develop`.
3. После закрытия PR можно продолжать так же пользоваться веткой, подтягивая другие изменения из `develop`.

## 2. Документация

1. Докуметнация REST API должна находится в `docs/api/`, оформленная в виде Markdown-файлов с названием `<название-сервиса>.md`
2. Модели для документации должны лежать в `docs/models/`, в виде `<название-модели>.json`
3. Документация каждого сервиса должна содержать:
   - Функциональное название endpoint-а
   - Уровень доступа, ограничения
   - HTTP-метод, URI запроса (возможны пояснения значений переменных пути)
   - Описание необходимых заголовков (при наличии)
   - Тело запроса в виде JSON (при наличии), также может представлять из себя модель из `docs/models/`.
   - Возможные варианты ответа сервиса, каждый из которых включает:
     - HTTP-код ответа
     - Тело ответа (при наличии)

## 3. Команда разработки

- DevBow Team:
  - [Демиденко Никита](https://github.com/TheOneFoxAgo)
  - [Дудкина София](https://github.com/sssidkn)
  - [Козакова Анна](https://github.com/nutochk)
  - [Кравченко Никита](https://github.com/Glaz0k)
  - [Лебедев Антон](https://github.com/IastWish)
  - [Новохацкий Данил](https://github.com/katagiriwhy)
  - [Пиявкин Антон](https://github.com/Piyavva)
