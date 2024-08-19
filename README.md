# test3
**Используемые технологии:**

- Go
- JWT
- PostgreSQL

**Итог:**

Написана часть сервиса аутентификации.

Два REST маршрута:

- Первый маршрут выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
- Второй маршрут выполняет Refresh операцию на пару Access, Refresh токенов

**Требования:**

Access токен тип JWT, алгоритм SHA512.

Refresh токен  формат передачи base64, хранится в базе исключительно в виде bcrypt хеша,  защищен от изменения на стороне клиента и попыток повторного использования.

Access, Refresh токены обоюдно связаны, Refresh операцию для Access токена можно выполнить только тем Refresh токеном который был выдан вместе с ним.

Payload токенов содержет сведения об ip адресе клиента, которому он был выдан. В случае, если ip адрес изменился, при рефреш операции нужно послать email warning на почту юзера (для упрощения можно использовать моковые данные).
В базе данных хранится две таблички:
create table Users (
	guid text,
	email text  not null default 'nothing'
);
и
create table SessionUsers (
	id serial primary key,
	refreshToken text,
	ip text  
);
alter  table SessionUsers add column guid text UNIQUE references Users(guid);
Конфигурация задаётся через json файл. Под одним GUID может залогинен только один пользователь(в моменте не придумал, как отслеживать пользователей, если они будут под одним GUID).
