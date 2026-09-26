# Задание 1

**Домен: Управление данными пользователя**
Поддомен: Управление личным кабинетом (профилем)

- Контекст - микросервис - Управление личными данными - Управление данными ЛК (личная информация, устройтва пользователя и т.д.)

Поддомен: Управление УЗ

- Контекст - микросервис - Управление УЗ - Регистрация, создание ЛК, восстановление доступа

**Домен: Управление контентом**
Поддомен: Управление каталогом контента

- Контекст - микросервис - Управление каталогом - Публикация каталога контента. Фильмов, видео и т.д. Установление зависимости контента от устройства

- Контекст - микросервис - Управление интеграцией с *сервисом* - Получение и обновление контента, который предоставлет сервис. Установление типа контента

- Контекст - микросервис - Управление персонализацией - Управление функциями истории просмотров, избранного, оценками и т.д.

- Контекст - микросервис - Управление правами - Проверка прав доступа на контент, сверка с данными внешнего сервиса

Поддомен: Управление просмотром

- Контекст - Управление просмотром - Воспроизведение контента на разных устройствах

**Домен: Управление монетизацией**
Поддомен: Управление платежами

- Контекст - микросервис - Платежный сервис - Управление жизненным циклом платежей, возвратами, подтверждениями и т.д.

Поддомен: Управление ценами и скидками

- Контекст - микросервис - Управление каталогом-продуктов - Тарифные планы и что в них входит (источники контента, устройства и т.д.), ведение каталога цен, доступности.

- Контекст - микросервис - Управление скидками - Управление жизненным циклом скидок и других акций

Поддомен: Управление подписками и правами доступа

- Контекст - микросервис - Управление подпиской - Управление жизненным циклом подписки пользователя. Установление права доступа к контенту по подписке

- Контекст - микросервис - Управление правами доступа - Ведение реестра источников прав и предоставление информации о том, что доступно конкретному пользователю

Поддомен: Управление партнерской монетизацией

- Контекст - микросервис - Управление программами лояльности - Баллы, статусы, кэшбэки и прочее, что оформляется и оплачивается не в нашем ПО

- Контекст - микросервис - Управление программами лояльности *сервис-партнер* - *сервис-партнеры* предлагающие оплату нашей подписки (по программе лояльности) у себя в ПО + любая друга бизнес-логика. 1 сервис на 1 интеграцию с партнером

[Карта контекстов (без взаимодействий)](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/diagrams/c4-domains-view/c4-domains-drawio-png.png)

[Диаграмма контейнеров](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/diagrams/c4-container-view/c4-container-view.puml)

Сервис Movie (который сейчас пользуется общей БД с монолитом) - это будущий catalogManagement

# Задание 2

- Тесты запускал 4 раза. Были ошибки с createSubscription

- Был установлен newman (репорты есть в папке тестов)

### 1. Proxy
С кодом помогал агент

[Скрин-тестов](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/postman-tests.png)

### 2. Kafka
С кодом помогал агент

[Брокер](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kafka-brokers.png)
[Сообщения в топике movie](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kafka-movie-messages.png)
[Сообщения в топике payment](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kafka-payment-messages.png)
[Сообщения в топике user](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kafka-user-messages.png)
[Топики](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kafka-topics.png)

# Задание 3

[/api/movies](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kube/api-movies.png)
[tests](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kube/tests.png)
[event-service](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/kube/event-service-tail.png)

# Задание 4

[deploy](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/helm/depl.png)
[curl](https://github.com/killtoyz/architecture-cinemaabyss/blob/cinema/screenshot/helm/curl.png)
