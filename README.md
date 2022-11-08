<!DOCTYPE html>
<header>
    <h1>
        wb-service
    </h1>
    <h3>
        Задание:
    </h3>
</header>
<body>
    <h3>
        В БД:
    </h3>
    <ul>
        <li>Развернуть локально postgresql</li>
        <li>Создать свою бд</li>
        <li>Настроить своего пользователя. </li>
        <li>Создать таблицы для хранения полученных данных.</li>
    </ul> 
    <h3>
        В сервисе:
    </h3>
    <ol>
        <li>Подключение и подписка на канал в nats-streaming</li>
        <li>Полученные данные писать в Postgres</li>
        <li>Так же полученные данные сохранить in memory в сервисе (Кеш)</li>
        <li>В случае падения сервиса восстанавливать Кеш из Postgres</li>
        <li>Поднять http сервер и выдавать данные по id из кеша</li>
        <li>Сделать простейший интерфейс отображения полученных данных, для их запроса по id</li>
    </ol>
    <h4>
        Доп инфо: 
    </h4>
    <ul>
        <li>Данные статичны, исходя из этого подумайте насчет модели хранения в Кеше и в pg. Модель в файле model.json</li>
        <li>В канал могут закинуть что угодно, подумайте как избежать проблем из-за этого</li>
        <li>Чтобы проверить работает ли подписка онлайн, сделайте себе отдельный скрипт, для публикации данных в канал</li>
        <li>Подумайте как не терять данные в случае ошибок или проблем с сервисом</li>
        <li>Nats-streaming разверните локально ( не путать с Nats )</li>
    </ul>
</body>
