# 🏦 Yobank — Telegram Mini App

Yobank — это современное банковское приложение, работающее внутри Telegram как Mini App. Проект реализован на двух основных технологиях: Go + Gin для backend и React + Typescript для frontend, использует PostgreSQL и поддерживает быструю сборку/развёртывание через Docker Compose.

---

## 📂 Структура репозитория

```

├── api/         # Backend: Go + Gin, бизнес-логика, REST API, работа с PostgreSQL и Kafka, Telegram бот
├── client/      # Frontend: React (Typescript), Telegram UI, Redux
├── docker-compose.yml
├── README.md
```

---

## 🚀 Быстрый старт

1. **Клонирование репозитория**

   ```
    git clone https://github.com/yourname/yobank.git
    cd yobank
   ```

2. **Запуск всех сервисов через Docker Compose**

   ```
    docker-compose up --build
   ```

   После сборки сервисы будут доступны по адресам (см. порты в docker-compose.yml):

    * Frontend (клиент): [http://localhost:5173](http://localhost:5173)
    * Backend (API): [http://localhost:8080](http://localhost:8080)

3. **Настройка переменных окружения**

   Для backend и frontend доступны примеры `.env.example`.
   Скопируй их и заполни своими данными:

   ```
    cp api/.env.example api/.env
    cp client/.env.example client/.env
   ```

---

## 📦 Описание сервисов

### Backend (`/api`)

* Go + Gin
* JWT-авторизация, поддержка Telegram Mini App и авторизация по email-коду
* P2P-переводы между пользователями
* История транзакций, кошельки, публичные username
* Асинхронные уведомления (email, Telegram)
* Интеграция с PostgreSQL и Kafka
* Dockerfile: `api/Dockerfile`

### Frontend (`/client`)

* React + Typescript
* Telegram UI Kit (адаптация под WebApp)
* Авторизация через Telegram и email
* Redux для хранения состояния пользователя, кошельков и истории
* Интерактивные экраны: баланс, переводы, история, пополнение, профиль
* Dockerfile: `client/Dockerfile`

---

## 🛠️ Полезные команды

**Backend:**

```
cd api
go run main.go          # Запуск API сервиса локально
```

**Frontend:**

```
cd client
npm install
npm run dev             # Запуск фронта (React) локально
```

---

## ⚙️ Конфигурация

* Все переменные окружения указываются в `.env` файлах в `api/` и `client/`.
* Для production добавь свои секретные ключи и доменные адреса.

---

## 🌐 Разворачивание и продакшн

* Запуск через Docker Compose: `docker-compose up --build`
* Для HTTPS и публичного доступа рекомендуется проксировать frontend/backend через nginx (пример конфигов см. в репозитории или спрашивай в issues).
* Для локальной отладки Telegram Mini App используй [ngrok](https://ngrok.com/) или SSH reverse tunnel.

---

## 💡 Ключевые возможности

* Авторизация через Telegram Mini App и email-код (без пароля)
* Уникальные публичные username для переводов между пользователями
* Баланс, пополнение, перевод, история транзакций
* Поддержка уведомлений (email/Telegram)
* Продвинутая архитектура на Go и React
* Быстрый запуск и масштабирование через Docker Compose

---