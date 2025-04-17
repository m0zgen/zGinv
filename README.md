# zGinv

📦 **zGinv** — CLI-инструмент для централизованной инвентаризации и управления 
VPS-серверами. Хранит данные в SQLite, 
поддерживает фильтрацию, экспорт в Ansible, импорт из SSH-конфигов и работу с группами. Также доступен HTTP API сервер на Fiber v3.

---

## 🛠️ Возможности

- Добавление серверов вручную (`add`)
- Импорт хостов из SSH-конфигов (`import-ssh`)
- Просмотр серверов с фильтрами (`list`)
- Поиск по имени, тегу, группе (`find`)
- Редактирование серверов (`edit`)
- Экспорт в форматы: `ansible`, `csv`, `json`, `yaml`
- Поддержка динамического инвентаря (`inventory`)
- Запуск Web API через `serve`
- Документация через Swagger UI (`/swagger/index.html`)
- Просмотр и анализ групп (`groups`)

---

## 📡 Запуск API-сервера

```bash
zGinv serve
```

### Параметры:
- `--port`, `-p` — указать порт вручную
- `ZGINV_PORT=9090` — задать порт через переменную окружения

Пример:
```bash
zGinv serve -p 7070
```

API будет доступен на: `http://localhost:7070/api/`

### 📘 Примеры HTTP-запросов:

```bash
curl "http://localhost:8080/api/servers?group=base-dns&project=openbld"
curl -X POST http://localhost:8080/api/servers \
  -H "Content-Type: application/json" \
  -d '{"name": "test", "address": "1.2.3.4"}'

curl -X PUT http://localhost:8080/api/servers/test \
  -H "Content-Type: application/json" \
  -d '{"address": "5.6.7.8"}'

curl -X DELETE http://localhost:8080/api/servers/test

curl http://localhost:8080/api/groups
curl http://localhost:8080/api/servers/export
```

---

## 📚 Swagger / OpenAPI документация

Установка:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

Swagger UI будет доступен по адресу:
```http
http://localhost:8080/swagger/index.html
```

Если ты используешь `swag init`, он сгенерирует файл `docs/swagger.json`.

Примеры аннотаций в коде:
```go
// @Summary Получить список серверов
// @Tags Servers
// @Produce json
// @Success 200 {array} db.Server
// @Router /servers [get]
```

Модель `Server` оформлена с `json`, `example`, `description`, `format`:
```go
type Server struct {
	Name string `json:"name" example:"dns-kz-1" description:"Имя сервера"`
	... // и т.д.
}
```

Direct call with `curl`:
```bash
curl http://localhost:8080/swagger/doc.json
```

Swagger UI использует [`github.com/Flussen/swagger-fiber-v3`](https://github.com/Flussen/swagger-fiber-v3) — поддержка Fiber v3.

---

## 📋 Примеры использования CLI

### ✅ Добавление сервера вручную
```bash
zGinv add \
  --name dns-kz-1 \
  --address 185.100.100.1 \
  --project openbld \
  --region kz \
  --tags dns,edge \
  --group base-dns \
  --comment "Новый DNS сервер"
```

### ✏️ Редактирование существующего сервера
```bash
zGinv edit dns-kz-1 \
  --address 185.100.100.99 \
  --user ubuntu \
  --tags dns,edge,updated
```

### 🔍 Поиск серверов
```bash
zGinv find --name "dns-kz*"
zGinv find --group base-dns
zGinv find --tag edge
```
→ Поддерживаются шаблоны `*` и фильтрация по группе или тегам.

### 📥 Импорт из SSH-конфига
```bash
zGinv import-ssh --file ~/path/to/ssh/configs/bld.conf
```
→ Группа будет автоматически определена из имени файла (`bld`).

### 📄 Просмотр серверов
```bash
zGinv list
zGinv list --project openbld
zGinv list --group base-dns
```

### 📤 Экспорт в разные форматы

#### 🔹 Ansible INI формат:
```bash
zGinv export --format ansible > hosts.ini
zGinv export --group base-dns --format ansible > core.ini
```

#### 🔹 CSV формат:
```bash
zGinv export --format csv > servers.csv
```

#### 🔹 JSON формат:
```bash
zGinv export --format json --group bld > bld.json
```

#### 🔹 YAML (структура групп для Ansible)
```bash
zGinv export --format yaml > all-hosts.yml
```
```yaml
all:
  children:
    base-dns:
      hosts:
        dns-kz-1: {}
        dns-kz-2: {}
    edge:
      hosts:
        edge-1: {}
```

---

### ⚡ Динамический инвентарь для Ansible

Создай скрипт:
```bash
echo '#!/bin/bash
zGinv inventory' > inventory.sh
chmod +x inventory.sh
```

В `ansible.cfg`:
```ini
[defaults]
inventory = ./inventory.sh
```

Теперь можешь:
```bash
ansible all -m ping
ansible-playbook site.yml
```

---

### 🧭 Просмотр всех групп
```bash
zGinv groups
```
Выведет список групп и количество серверов в каждой.

---

## ⚖️ Лицензия

CC0 License © 2025 [Евгений Гончаров](https://openbld.net)
