# Project Golang-catalog

```text
Golang, ElasticSearch, Kafka
```

## Getting Started

Для начальной настройки проекта необходимо выполнить команду:
```bash
make install
```
Она добавит данные в etc/hosts, добавит нужжные файлы и запустит проект. После запуска перейти на роут **"/createIndexes"** для создания индексов в ElasticSearch.
Для заполнения тестовыми данными через кафку, нужно перейти на роуты (Прежде чем отправить данные, нужно создать топики в кафке, по 3 партиции в каждом):
```text
/sendCategories - отправка тестовых категорий в кафку
/sendProducts - отправка тестовых продуктов в кафку
/sendProperties - отправка тестовых свойств в кафку
```

### Api
```text
http://catalog.local:8080/v1/...
```

### Kafka
```text
http://ui.catalog-kafka.local:8090/ - UI kafka
http://catalog-kafka.local:9092 - api
```

### ElasticSearch
```text
http://kibana:5601 - Kibana
http://elasticsearch:9200 - api
```

## MakeFile

Установка проекта и его запуск в докере
```bash
make install
```

Запуск билда проекта:
```bash
make build
```

Запуск проекта без билда:
```bash
make up
```

Остановка проекта:
```bash
make down
```

Удаление бинарного файла:
```bash
make clean
```
