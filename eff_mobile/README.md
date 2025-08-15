# Subscription Aggregation Service

![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)
![Docker](https://img.shields.io/badge/Docker-Powered-blue.svg)
![PostgreSQL](https://img.shields.io/badge/Database-PostgreSQL-blue.svg)

Этот проект представляет собой REST-сервис, реализованный на Go, для агрегации данных об онлайн-подписках пользователей. Он построен с использованием фреймворка Echo, драйвера pgx для PostgreSQL и следует принципам чистой архитектуры.

## ✨ Особенности

-   **REST API**: Полный набор CRUDL (Create, Read, Update, Delete, List) операций для управления подписками.
-   **Агрегация данных**: Отдельный эндпоинт для подсчета суммарной стоимости подписок за выбранный период с возможностью фильтрации.
-   **База данных**: Используется PostgreSQL с системой версионирования схемы через `golang-migrate`.
-   **Конфигурация**: Гибкая настройка через YAML-файл.
-   **Логирование**: Структурированное логирование в консоль и файл.
-   **Документация**: Автоматически генерируемая Swagger (OpenAPI) документация.
-   **Контейнеризация**: Полностью готовый к запуску с помощью Docker и Docker Compose.

## 🏗️ Структура проекта

Проект организован в соответствии с принципами чистой архитектуры для четкого разделения ответственности.

```
.
├── cmd/                  # Точка входа в приложение
│   └── main.go
├── config/               # Структуры и загрузчик конфигурации
├── docs/                 # Сгенерированная Swagger-документация
├── internal/             # Внутренняя логика приложения
│   ├── db/               # Управление жизненным циклом подключения к БД
│   ├── handler/          # Обработчики HTTP-запросов (слой API)
│   ├── initial/          # Инициализация зависимостей
│   ├── model/            # Модели данных и ошибки домена
│   ├── repository/       # Реализация доступа к данным (слой репозитория)
│   └── service/          # Бизнес-логика (слой сервисов)
├── migrations/           # Файлы миграций для базы данных
├── pkg/                  # Переиспользуемые пакеты
│   ├── pdb/              # Обёртка над пулом соединений pgx
│   └── service/          # Общий интерфейс для сервисов
├── config.yaml           # Файл конфигурации
├── docker-compose.yaml   # Конфигурация для запуска через Docker Compose
├── Dockerfile            # Dockerfile для сборки приложения
├── go.mod                # Зависимости проекта
└── README.md             # Этот файл
```

## 🚀 Технологии

-   **Язык**: Go
-   **Веб-фреймворк**: [Echo](https://echo.labstack.com/)
-   **База данных**: PostgreSQL
-   **Драйвер БД**: [pgx](https://github.com/jackc/pgx)
-   **Миграции**: [golang-migrate](https://github.com/golang-migrate/migrate)
-   **Документация API**: [swaggo/swag](https://github.com/swaggo/swag)
-   **Контейнеризация**: Docker, Docker Compose

## ⚙️ Начало работы

### Требования

-   Go 1.24+
-   Docker
-   Docker Compose
-   [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) (для локального запуска)

### Запуск с помощью Docker (Рекомендуемый способ)

Это самый простой и надежный способ запустить проект. Все зависимости, включая базу данных и миграции, управляются автоматически.

1.  **Клонируйте репозиторий:**
    ```bash
    git clone <your-repository-url>
    cd eff_mobile
    ```

2.  **Запустите Docker Compose:**
    ```bash
    docker-compose up --build
    ```
    Эта команда выполнит следующие шаги:
    -   Соберет Docker-образ вашего Go-приложения.
    -   Запустит контейнер с PostgreSQL.
    -   Дождется, пока база данных будет готова принимать подключения.
    -   Запустит контейнер с `golang-migrate`, который применит миграции к базе.
    -   После успешного завершения миграций запустит основной контейнер с вашим API.

Сервис будет доступен по адресу `http://localhost:8080`.

### Локальный запуск (без Docker)

1.  **Клонируйте репозиторий.**

2.  **Установите зависимости:**
    ```bash
    go mod tidy
    ```

3.  **Настройте и запустите PostgreSQL** локально.

4.  **Создайте `config.yaml`** или измените существующий, указав корректную строку подключения к вашей локальной базе данных.
    ```yaml
    database:
      connection-string: "postgresql://user:password@localhost:5432/dbname?sslmode=disable"
      # ...
    ```

5.  **Примените миграции базы данных:**
    Выполните команду, подставив вашу строку подключения.
    ```bash
    migrate -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" -path ./migrations up
    ```

6.  **Запустите приложение:**
    ```bash
    go run ./cmd/main.go
    ```

## 📖 Документация API

API сервиса документировано с помощью Swagger (OpenAPI).

После запуска сервиса (любым способом) интерактивная документация будет доступна по адресу:

**[http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)**

### Основные эндпоинты

-   `POST /subscriptions`: Создать новую подписку.
-   `GET /subscriptions`: Получить список всех подписок.
-   `GET /subscriptions/{id}`: Получить подписку по ID.
-   `PUT /subscriptions/{id}`: Обновить подписку по ID.
-   `DELETE /subscriptions/{id}`: Удалить подписку по ID.
-   `GET /subscriptions/sum`: Рассчитать суммарную стоимость подписок за период с фильтрами.

## 🗄️ Миграции базы данных

Управление версиями схемы базы данных осуществляется с помощью `golang-migrate`. Файлы миграций находятся в папке `/migrations`.

-   **Применить все новые миграции:**
    ```bash
    migrate -database "YOUR_DB_URL" -path ./migrations up
    ```
-   **Откатить последнюю примененную миграцию:**
    ```bash
    migrate -database "YOUR_DB_URL" -path ./migrations down 1
    ```