# Документация к API

## 1. Обзор

Этот API предназначен для работы с заметками (добавление, редактирование, удаление, просмотр). Поддерживает аутентификацию, управление сессиями. Все данные (записки, пользователи, сессии) хранятся в БД MySQL.<br>
Для аутентификации используется JWT-токен. У токена есть время действия (выдается на 3 минуты), по истечению времени сессия становится недействительной (удаляется) и необходимо заново авторизоваться. Регистрация не предусмотрена, существует только два пользователь по умолчанию (`username: admin, password: admin` и `username: user123, password: user321`).<br>
Сервис предоставляет API для входа, создания записки, редактирования, удаления и просмотра, соответственно по `POST /api/login`, `POST /api/note`, `PUT /api/update/{id}`, `DELETE /api/delete/{id}`, `GET /api/note`.<br>
Реализован паттерн репозиторий (разделенные слои репозитория, логики, и хендлеров). Аутентификация и логирование реализованы с помощью middleware. Конфигурацию адресса и порта MySQL можно настроить в `./config/config.yml`. <br>
<br>
Для удобства, в корне репозитория лежит Postman коллекция `notes.postman_collection.json`. Там же приведены примеры запросов и возможные ошибки.

## 2. Предварительные требования

Перед запуском убедитесь, что у вас установлен Docker.

## 3. Быстрый запуск

1. Клонируйте репозиторий:

```bash
git clone https://github.com/tarbeevms/intern-kode
```

2. Войдите в скопированную директорию в папку ./build:

```bash
cd intern-kode/build
```

3. Запустите Docker Compose:

```bash
docker compose up 
```

После запуска docker compose запустятся два контейнера: MySQL и Notes с БД и API соответственно. У приложения есть свой Dockerfile.
API будет запущено на `localhost:8080`, а MySQL на `mysql:3306`. Порт и адресс, на котором запускается Mysql можно поменять в `./config/config.yml`.

## 4. API Эндпоинты

При попытке отправки неккоректного ответа и неккоректного чтения запроса json формата во всех эндпоинтах предусмотрен хендлинг ошибок.

### 1. Аутентификация (POST /api/login)
Этот эндпоинт выполняет аутентификацию пользователя и возвращает JWT-токен.<br>
Можно повторно авторизоваться (в таком случае, токен сессии в БД для данного юзера перезапишется, обновится).
Пример запроса: </br>
```json
POST /api/login HTTP/1.1
Host: localhost:8080
Content-Type: application/json
{
  "username": "admin",
  "password": "presale"
}
```
Ответ, статус `200 Ok`: 
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiYWRtaW4iLCJleHAiOjE3MjQxMTExNDV9.0kn4u7X-JhO-eHZf8IOc_zWKONL42-tvFMbKkD1fibo"
}
```

При отсылке некорректных данных получим ошибку и статус `401 Unauthorized`:
```json
{
    "message": "Invalid username or password"
}
```
При ошибке создании (перезаписи) сессии в БД, получим ошибку и статус `401 Unauthorized`:
```json
{
    "message": "Failed to create (rewrite) session"
}
```
При попытке отправки ззапроса к выключенной БД получим ошибку, статус `401 Unauthorized`
```json
{
    "message": "failed to query user: dial tcp: lookup mysql on 127.0.0.11:53: server misbehaving"
}
```
### 2. Создание записки (POST /api/note)
Этот эндпоинт позволяет авторизованным пользователям создавать записки, которые будут храниться в MySQL. Исправляет орфографические ошибки за счет интеграции с Яндекс.Спеллер.
Пример запроса:
```json
POST /api/note HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiYWRtaW4iLCJleHAiOjE3MjQxMTExNDV9.0kn4u7X-JhO-eHZf8IOc_zWKONL42-tvFMbKkD1fibo
Content-Type: application/json
{
    "text": "пошол дожд"
}
}
```
Ответ, статус `200 Ok`:
```json
{
    "id": "5f0f03b5-60dc-4214-966d-0682d99d646f",
    "text": "пошел дождь",
    "author": "admin",
    "created_at": "2024-08-26T11:28:39.442048339Z",
    "updated_at": "2024-08-26T11:28:39.442048389Z"
}
```
При попытке отправки неккоректного JWT токена получим ошибку, статус `401 Unauthorized`:
```json
{
    "message": "Not authorized: signature is invalid"
}
```
При попытке отправки неккоректного Authorization заголовка получим ошибку, статус `401 Unauthorized`:
```json
{
    "message": "Not authorized, wrong header format"
}
```
При попытке отправки просроченого JWT токена получим ошнибку, статус `401 Unauthorized`  (просроченная сессия удаляется из БД и необходимо заново авторизоваться):
```json
{
    "message": "Not authorized: session expired"
}
```
При попытке отправки несуществующей (или же уже удаленной) сессии получим ошибку, статус `401 Unauthorized`:
```json
{
    "message": "Not authorized: session not found"
}
```
При попытке отправки плохого текста (не строчного типа) получим ошибку, статус `400 Bad Request`.
Пример запроса:
```json
{
    "text": 1
}
```
Ответ:
```json
{
    "message": "Invalid request payload"
}
```
При записи в выключенную БД, выдастся ошибка, статус `401 Bad Request`:
```json
{
    "message": "Not authorized: failed to query session: dial tcp: lookup mysql on 127.0.0.11:53: server misbehaving"
}
```
### 3. Обновление записки (PUT /api/note/{id})
Этот эндпоинт позволяет авторизованным пользователям обновлять (свои) записки по ID.
Пример запроса:
```json
PUT /api/note/5f0f03b5-60dc-4214-966d-0682d99d646f HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiYWRtaW4iLCJleHAiOjE3MjQxMTExNDV9.0kn4u7X-JhO-eHZf8IOc_zWKONL42-tvFMbKkD1fibo
Content-Type: application/json
{
    "text": "дожд прекротился"
}
```
Ответ (при том, что в БД есть записка с таким ID), статус `200 Ok`:
```json
{
    "id": "5f0f03b5-60dc-4214-966d-0682d99d646f",
    "text": "дождь прекратился",
    "author": "admin",
    "created_at": "2024-08-26T11:28:39Z",
    "updated_at": "2024-08-26T11:35:24.682051953Z"
}
```
При вводе некорректного ID получим ошибку, статус `400 Bad Request`:
```json
{
    "message": "Invalid note ID"
}
```
При попытке чужого юзера (например `user123`) изменить записку юзера (например `admin`), получим ошибку, статус `400 Bad Request`:
```json
{
    "message": "No permission to update note"
}
```
При записи в выключенную БД, выдастся ошибка, статус `401 Unauthorized`:
```json
{
    "message": "Not authorized: failed to query session: dial tcp: lookup mysql on 127.0.0.11:53: server misbehaving"
}
```
При попытке отправки неккоректного JWT токена получим ошибку, статус `401 Unauthorized`:
```json
{
    "message": "Not authorized: signature is invalid"
}
```
При попытке отправки неккоректного Authorization заголовка получим ошибку, статус `401 Unauthorized`:
```json
{
    "message": "Not authorized, wrong header format"
}
```
При попытке отправки просроченого JWT токена получим ошнибку, статус `401 Unauthorized`  (просроченная сессия удаляется из БД и необходимо заново авторизоваться):
```json
{
    "message": "Not authorized: session expired"
}
```
При попытке отправки несуществующей (или же уже удаленной) сессии получим ошибку, статус `401 Unauthorized`:
```json
{
    "message": "Not authorized: session not found"
}
```
### 4. Получение всех записок юзера (GET /api/note)
Этот эндпоинт позволяет авторизованным пользователям получать все (свои) записки.
Пример запроса:
```json
GET /api/note HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiYWRtaW4iLCJleHAiOjE3MjQxMTExNDV9.0kn4u7X-JhO-eHZf8IOc_zWKONL42-tvFMbKkD1fibo
```
Ответ (несколько записей), статус `200 Ok`:
```json
[
    {
        "id": "9b4bde12-c406-442f-a432-d53addc74748",
        "text": "forget to push the door",
        "author": "admin",
        "created_at": "2024-08-26T11:23:51Z",
        "updated_at": "2024-08-26T11:23:51Z"
    },
    {
        "id": "5f0f03b5-60dc-4214-966d-0682d99d646f",
        "text": "пошел дождь",
        "author": "admin",
        "created_at": "2024-08-26T11:28:39Z",
        "updated_at": "2024-08-26T11:28:39Z"
    }
]
```
Ответ (ноль записей), статус `200 Ok`:
```json
[]
```
Ошибки связанные с сессией/токеном/авторизацией и БД аналогичные.

### 4. Удаление записки юзера по ID (DELETE /api/note/{id})
Этот эндпоинт позволяет авторизованным пользователям удалять (свою) записку по ID.
Пример запроса:
```json
DELETE /api/note/5f0f03b5-60dc-4214-966d-0682d99d646f HTTP/1.1
Host: localhost:8080
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJOYW1lIjoiYWRtaW4iLCJleHAiOjE3MjQxMTExNDV9.0kn4u7X-JhO-eHZf8IOc_zWKONL42-tvFMbKkD1fibo
```
Ответ (при том, что в БД есть записка с таким ID), статус `200 Ok`:
```json
{
    "message": "Deleted sucessfully"
}
```
При вводе некорректного ID получим ошибку, статус `400 Bad Request`:
```json
{
    "message": "Invalid note ID"
}
```
При попытке чужого юзера (например `user123`) удалить записку юзера (например `admin`), получим ошибку, статус `400 Bad Request`:
```json
{
    "message": "No permission to update note"
}
```

Ошибки связанные с сессией/токеном/авторизацией и БД аналогичные.



