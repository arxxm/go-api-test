Сервис, предоставляющий API для работы с данными пользователей.<br/>
В качестве хранилища данных используется Postgres.<br/>
Запросы на получение списка кешируются в Redis на 5 минут.<br/>
Данные на получение юзера кешируются в памяти на 5 минут.

#### Реализованные методы:

| Method   | URL                  | Description                                          | Example                                                                                                                                                                                                                 |
|----------|----------------------|------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `GET`    | `/api/v1/users/`     | Получение данных с различными фильтрами и пагинацией | `/api/v1/users`                                                                                                                                                                                                         |
| `PUT`    | `/api/v1/users/{id}` | Изменение сущности                                   | `/api/v1/users/1`                                                                                                                                                                                                       |
| `POST`   | `/api/v1/users/`     | Добавление новой сущности                            | `/api/v1/users`<br/>`body:`<br/>`{`<br/>`"name": "Jack",`<br/>`"last_name": "Lee",`<br/>`"surname": "Smith",`<br/>`"gender": "male",`<br/>`"status": "active",`<br/>`"date_of_birth": "1990-01-01T00:00:00Z" `<br/>`}`  |
| `DELETE` | `/api/v1/users/`     | Удаление по идентификатору                           | `/api/v1/users/`<br\> `body:`<br/>`{`<br/>`"id": "1"`<br/>`}`                                                                                                                                                           |
| `GET`    | `/api/v1/users/{id}` | Получение по идентификатору                          | `/api/v1/users/1`                                                                                                                                                                                                       |



Конфигурационные данные вынесены в .env
