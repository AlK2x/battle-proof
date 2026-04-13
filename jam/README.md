
# 2. Jam Service

## Контекст

Управляет соревнованиями (чемпионатами, фестивалями) и их участниками. При добавлении участника синхронно валидирует пользователя через GraphQL запрос в User Service. Публикует события в Kafka для downstream-сервисов.

## Модели

```go
type Jam struct {
 ID string
 Name string
 Location string
 Date time.Time
 CreatedBy string // user ID организатора
}

type Participant struct {
 JamID string
 UserID string
 Role JamRole
}

type JamRole string
const (
 EventRoleDancer JamRole = "dancer"
 EventRoleJudge JamRole = "judge"
 EventRoleMedia JamRole = "media"
)
```

## API

```
POST /jams
GET /jams
GET /jams/{id}
POST /jams/{id}/participants
DELETE /jams/{id}/participants/
GET /jams/{id}/participants
GET /health
```

## Функциональные требования

### ФТ-1: Создание соревнования — `POST /jams`

- Принять: `name`, `location`, `date`, `created_by`
- Проверить через GraphQL что `created_by` существует и имеет роль `organizer`
- Сгенерировать UUID, сохранить
- Опубликовать в Kafka
- Вернуть `201 Created`

**Бизнес-правила:**
- `name`, `location`, `date` — обязательны
- `date` не может быть в прошлом
- Только пользователь с ролью `organizer` может создать соревнование

**Kafka message** → топик `jams`:
```json
{
 "event_type": "event.created",
 "event_id": "uuid",
 "name": "Summer Championship 2025",
 "location": "Tallinn",
 "date": "2025-08-01T18:00:00Z",
 "created_by": "user-uuid",
 "occurred_at": "2025-03-13T10:00:00Z"
}
```

### ФТ-2: Получение списка — `GET /jams`

- Вернуть все соревнования
- Пустой список → `[]`, не ошибка
- Код: `200 OK`

### ФТ-3: Получение по ID — `GET /jams/{id}`

- Найти по `id` → `200 OK`
- Не найдено → `404 Not Found`

> Внешний контракт для Battle Service — не ломать.

### ФТ-4: Добавление участника — `POST /jams/{id}/participants`

- Принять: `user_id`, `role`
- Проверить что соревнование существует
- Сделать GraphQL запрос в User Service:
```graphql
query {
 user(id: "uuid") {
 id
 role
 level
 }
}
```
- Сохранить участника
- Опубликовать в Kafka
- Вернуть `201 Created`

**Бизнес-правила:**
- Пользователь должен существовать → иначе `404`
- Один пользователь не может быть добавлен дважды в одно соревнование → `409 Conflict`
- Роль на событии **независима** от роли пользователя (dancer может быть judge на событии)




### ФТ-5: Список участников — `GET /jams/{id}/participants`

- Вернуть всех участников соревнования
- Соревнование не найдено → `404`
- Пустой список → `[]`

> Внешний контракт для Battle Service.

### ФТ-6: Health Check — `GET /health`

```json
{ "status": "ok" }
```

## Нефункциональные требования

### НФТ-1: Отказоустойчивость — User Service недоступен

- Возвращаем `503 Service Unavailable`
- Участник не добавляется, соревнование не создаётся
- Логируем с деталями ошибки

### НФТ-2: Порядок публикации в Kafka

Публикация происходит **строго после** успешного сохранения. Если сохранение упало — в Kafka ничего не летит.

### НФТ-3: Формат ошибок

```json
{ "error": "описание ошибки" }
```

### НФТ-4: Потокобезопасность

`sync.RWMutex` в репозитории.

### НФТ-5: Наблюдаемость

- Логировать каждый запрос: метод, путь, статус, время
- Логировать каждый GraphQL запрос к User Service: успех / ошибка
- Логировать каждую публикацию в Kafka: топик, тип события, успех / ошибка

---