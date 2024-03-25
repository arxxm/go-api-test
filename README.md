Сервис получает по апи ФИО и с помощью открытых апи обогащает
информацию о человеке наиболее вероятными возрастом, полом и национальностью.<br/>
Затем данные сохраняются в БД. 

#### Реализованные методы:

| Method   | URL                   | Description                                          | Example                                                                                                                              |
|----------|-----------------------|------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------|
| `GET`    | `/api/v1/persons`     | Получение данных с различными фильтрами и пагинацией | `/api/v1/persons`                                                                                                                    |
| `PATCH`  | `/api/v1/persons/:id` | Изменение сущности                                   | `/api/v1/persons/1`                                                                                                                  |
| `POST`   | `/api/v1/persons`     | Добавление новых сущностей                           | `/api/v1/persons`<br/>`body:`<br/>`{`<br/>`"name": "Dmitriy",`<br/>`"surname": "Ushakov",`<br/>`"patronymic": "Vasilevich" `<br/>`}` |
| `DELETE` | `/api/v1/persons/:id` | Удаление по идентификатору                           | `/api/v1/persons/1`                                                                                                                  |

###
####  Корректное сообщение перед записью в БД обогощается:
- Возрастом - https://api.agify.io/?name=Dmitriy
- Полом - https://api.genderize.io/?name=Dmitriy
- Национальностью - https://api.nationalize.io/?name=Dmitriy

Структура БД создана путем миграций(./migrations)<br/>
Конфигурационные данные вынесены в .env
