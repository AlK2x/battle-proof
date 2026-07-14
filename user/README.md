# 1. User Service

## Контекст

Центральный справочник пользователей. Выставляет **GraphQL API** — другие сервисы запрашивают его синхронно для валидации. В Kafka ничего не публикует.

## Модели

```go
type User struct {
ID string
Name string
Email string
Level *DancerLevel
CreatedAt time.Time
}


type DancerLevel string
const (
LevelBeginner DancerLevel = "beginner"
LevelAmateur DancerLevel = "amateur"
LevelPro DancerLevel = "pro"
)
```

## GraphQL Schema

```graphql
type User {
id: ID!
name: String!
email: String!
level: DancerLevel # null если не dancer
createdAt: String!
}

enum UserRole { dancer judge media organizer }
enum DancerLevel { beginner amateur pro }

input CreateUserInput {
name: String!
email: String!
level: DancerLevel
}

type Query {
user(id: ID!): User
users: [User!]!
}

type Mutation {
createUser(input: CreateUserInput!): User!
}
```

## Функциональные требования

### ФТ-1: Создание пользователя — `createUser` mutation

- Сгенерировать уникальный `ID` (UUID)
- Проставить `createdAt` автоматически на стороне сервера
- Сохранить пользователя в хранилище
- Вернуть созданного пользователя

**Бизнес-правила:**
- Если `role = dancer` — поле `level` **обязательно**; отсутствие → GraphQL error
- Если `role != dancer` — поле `level` **запрещено**; передача → GraphQL error
- `email` должен быть уникальным → повтор возвращает GraphQL error с кодом `CONFLICT`
- `name` и `email` не могут быть пустыми

### ФТ-2: Получение пользователя — `user(id)` query

- Найти пользователя по `id`
- Если не найден → вернуть `null` (стандартное поведение GraphQL)

### ФТ-3: Получение списка — `users` query

- Вернуть массив всех пользователей
- Если пользователей нет → вернуть `[]`

## Нефункциональные требования

### НФТ-0: Health Check — `GET /health`

```json
{ "status": "ok" }
``


### НФТ-1: Формат ошибок GraphQL

```json
{
"errors": [{ "message": "level is required for dancer role" }]
}
```

### НФТ-2: Потокобезопасность

`sync.RWMutex` в репозитории — `RLock` для чтения, `Lock` для записи.

### НФТ-3: Наблюдаемость

- Логировать каждый входящий GraphQL запрос: операция, статус, время выполнения
- Ошибки валидации → уровень `WARN`
- Внутренние ошибки → уровень `ERROR`


---