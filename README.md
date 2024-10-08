## Установка

### Требования

- Go 1.16 или новее
- Postgresql

### Склонируйте репозиторий

```
git clone git@github.com:Sirokko25/Cybersecurity-Center.git
cd final
go mod tidy
```

## Настройка
### Параметры окружения
Приложение использует три параметра окружения:

TODO_PORT: Порт, на котором будет работать сервер. По умолчанию используется порт 7070.
TIME_FORMAT: Формат времени для создания заметки
DSN: Строка для подключения к Postgre. В ней нужно указать свои параметры подключения.

## Инициализация базы данных

Приложение автоматически создаст файл базы данных и необходимые таблицы при первом запуске, если файл базы данных не существует.

## Запуск приложения
### Запуск сервера

Перейти в директорию cmd.

```
go run main.go
```
Сервер будет доступен по адресу http://localhost:7070

### Тестирование
Для тестирования имеется Postman коллекция.

### Структура проекта
- сmd/: Директория содержит главный файл приложения - точка входа сервера.
- internal/: Директория с логикой проекта:
    - auth/: Пакет для аунтефикации.
    - handlers/: Пакет с обработчиками API запросов.
    - storage/: Пакет для инициализации и работы с базой данных.
    - server/: Пакет для запуска сервера.
- models/: Пакет с необходимыми структурами и методами.
