# TODO
- Создать два окружения, чтобы не редактировать сборку в main.go и abstract_test.go
- Перевести остальные репозитории на Postgres
- Добавить автоматическое создание БД
- Добавить выдачу токена авторизации сразу при регистрации

- Пользователи
- - Создание пользователя = регистрация
- - - Шифрование пароля
- - - Проверка логина на существование при добавлении
- - - Проверка DTO регистрации
- - Аутентификация
- - - Проверка логина и пароля
- - - Генерация токена
- - - Авторизация по токены из заголовка
- - - Авторизация по токену из кук
- - - Проверка DTO авторизации
- Заказы
- - Добавление заказа
- - - Проверка авторизации
- - - Проверка номера заказа на повтор
- - - Проверка номера заказа по Луне
- - - Добавление записи в таблицу
- - - Проверка номера заказа на повтор
- - - Ошибка авторизации
- - - Ошибка DTO
- - - Ошибка по дублю заказа от текущего юзера
- - - Ошибка по дублю заказа от другого юзера
- - - Ошибка валидациии номера заказа по Луне
- - Список заказов
- - - Проверка авторизации
- - - Получение записей для юзера из БД
- - - Сортировка записей для юзера из БД
- - - Форматирование даты
- - - Ошибка когда данных не нашлось
- - - Ошибка авторизации
- - Просмотр заказа
- - - Проверка авторизации
- - - Ошибка авторизации
- - - Запрос в БД для получения
- - - Подсчет кол-ва обращений с IP
- - - Ошибка по кол-ву обращений с IP
- Баланс
- - Проверка баланса
- - - Проверка авторизации
- - - Запрос в БД для получения баланса
- - - Автоматическое создание записи в БД, если ее там нет
- - Списание
- - - Проверка авторизации
- - - Ошибка когда пользователь не авторизован
- - - Занесение записи в таблицу списания
- - - Расчет нового баланса после списания
- - - Проверка возможности списания
- - - Ошибка невозможности списания из-за нехватки баланса
- - - Проверка номера заказа на Луну
- - - Ошибка невалидного номер заказа
- - Список списаний
- - - Проверка авторизации
- - - Запрос в БД для получения списаний
- - - Сортировка списаний по дате списания
- - - Форматирование даты
- - - Ошибка когда данных не нашлось
- - - Ошибка когда пользователь не авторизован
+ Энвы и флаги
+ - Адрес и порт для веб сервера
+ - Адрес и порт для БД
+ - Мок адрес для расчета системы начислений

# Шаги

https://go.dev/doc/tutorial/create-module
go mod init github.com/volkoviimagnit/gofermart

https://github.com/sirupsen/logrus
go get github.com/sirupsen/logrus

https://github.com/joho/godotenv
go get github.com/joho/godotenv

https://github.com/caarlos0/env
go get github.com/caarlos0/env/v6

https://github.com/go-chi/chi
go get -u github.com/go-chi/chi/v5

https://github.com/stretchr/testify
go get -u github.com/stretchr/testify

https://github.com/ShiraazMoollatjie/goluhn
go get -u github.com/ShiraazMoollatjie/goluhn

https://github.com/go-resty/resty
go get -u github.com/go-resty/resty/v2

https://github.com/jackc/pgx
go get -u github.com/jackc/pgx/v4
go mod tidy

go build -o main && chmod +x main && LOG_LEVEL=trace DATABASE_URI=databaseEnv RUN_ADDRESS=runEnv ACCRUAL_SYSTEM_ADDRESS=accEnv ./main -a=localhostArg -d=databaseArg -r=accrualArg -ll=debug


cd /Users/volkov_ii/Projects/gofermart/cmd/gophermart && go build -buildvcs=false -o gophermart && \
cd /Users/volkov_ii/Projects/gofermart && chmod -R +x cmd/gophermart/gophermart && \
cd /Users/volkov_ii/Projects/gofermart && SERVER_PORT=$((8000+($RANDOM % 1000))) .tools/gophermarttest-darwin-arm64 \
-test.v -test.run=^TestGophermart$ \
-gophermart-binary-path=cmd/gophermart/gophermart \
-gophermart-host=localhost \
-gophermart-port=8081 \
-gophermart-database-uri="postgres://postgres:postgres@127.0.0.1:5433/gofermart?sslmode=disable" \
-accrual-binary-path=cmd/accrual/accrual_darwin_arm64 \
-accrual-host=localhost \
-accrual-port=8080 \
-accrual-database-uri="postgres://postgres:postgres@127.0.0.1:5433/gofermart?sslmode=disable" \
echo 'ok'