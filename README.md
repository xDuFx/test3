# test3
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
