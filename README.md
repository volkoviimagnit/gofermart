# go-musthave-diploma-tpl

Шаблон репозитория для индивидуального дипломного проекта курса «Go-разработчик»

# Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере.
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` — адрес вашего репозитория на GitHub без
   префикса `https://`) для создания модуля

# Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона, выполните команду:

```
git remote add -m master template https://github.com/yandex-praktikum/go-musthave-diploma-tpl.git
```

Для обновления кода автотестов выполните команду:

```
git fetch template && git checkout template/master .github
```

Затем добавьте полученные изменения в свой репозиторий.


# Схема БД

user
- id 
- login string not null
- password string not null
- token string null

user_order
- user.id string not null
- order_id string not null
- status_id string not null
- uploaded_at timestamp not null
- accrual int null
unique key user.id + order_id
unique key order_id

user_balance
- user.id
- balance float not null
- withdrawn float not null

user_withdrawn
- user.id
- order_id string not null
- sum float not null
- processed_at timestamp not null

