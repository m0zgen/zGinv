# zGinv

📦 **zGinv** is a CLI tool for centralized inventory and management of VPS servers.  
It stores data in SQLite, supports filtering, export to Ansible, SSH config import, and group management.  
Also includes an HTTP API server built with Fiber v3.

---

## 🛠️ Features

- Add servers manually (`add`)
- Import hosts from SSH configs (`import-ssh`)
- List servers with filters (`list`)
- Search by name, tag, or group (`find`)
- Edit server records (`edit`)
- Export to formats: `ansible`, `csv`, `json`, `yaml`
- Dynamic inventory support (`inventory`)
- Launch Web API with `serve`
- API documentation via Swagger UI (`/swagger/index.html`)
- Group inspection and analysis (`groups`)

---

## 📡 Running the API server

```bash
zGinv serve
```

### Parameters:
- `--port`, `-p` — specify port manually
- `ZGINV_PORT=9090` — set port via environment variable

Example:
```bash
zGinv serve -p 7070
```

The API will be available at: `http://localhost:7070/api/`

### 📘 Example HTTP requests:

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

## 📚 Swagger / OpenAPI documentation

Install Swagger CLI tools:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

Swagger UI will be available at:
```http
http://localhost:8080/swagger/index.html
```

Swagger JSON spec (for Postman, etc.):
```bash
curl http://localhost:8080/swagger/doc.json
```

Example annotations:
```go
// @Summary Get list of servers
// @Tags Servers
// @Produce json
// @Success 200 {array} db.Server
// @Router /servers [get]
```

Model `Server` is documented using `json`, `example`, `description`, `format`:
```go
type Server struct {
	Name string `json:"name" example:"dns-kz-1" description:"Server name"`
	// ...
}
```

Swagger UI is powered by [`github.com/Flussen/swagger-fiber-v3`](https://github.com/Flussen/swagger-fiber-v3) — compatible with Fiber v3.

---

## 📋 CLI Usage Examples

### ✅ Add a server manually
```bash
zGinv add \
  --name dns-kz-1 \
  --address 185.100.100.1 \
  --project openbld \
  --region kz \
  --tags dns,edge \
  --group base-dns \
  --comment "New DNS server"
```

### ✏️ Edit an existing server
```bash
zGinv edit dns-kz-1 \
  --address 185.100.100.99 \
  --user ubuntu \
  --tags dns,edge,updated
```

### 🔍 Search servers
```bash
zGinv find --name "dns-kz*"
zGinv find --group base-dns
zGinv find --tag edge
```
→ Supports `*` wildcard and filtering by group or tags.

### 📥 Import from SSH config
```bash
zGinv import-ssh --file ~/path/to/ssh/configs/bld.conf
```
→ Group name is derived automatically from the file name (`bld`).

### 📄 List servers
```bash
zGinv list
zGinv list --project openbld
zGinv list --group base-dns
```

### 📤 Export to different formats

#### 🔹 Ansible INI format:
```bash
zGinv export --format ansible > hosts.ini
zGinv export --group base-dns --format ansible > core.ini
```

#### 🔹 CSV format:
```bash
zGinv export --format csv > servers.csv
```

#### 🔹 JSON format:
```bash
zGinv export --format json --group bld > bld.json
```

#### 🔹 YAML (Ansible-style group structure)
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

### ⚡ Dynamic Ansible inventory

Create a wrapper script:
```bash
echo '#!/bin/bash
zGinv inventory' > inventory.sh
chmod +x inventory.sh
```

Set it in your `ansible.cfg`:
```ini
[defaults]
inventory = ./inventory.sh
```

Now you can run:
```bash
ansible all -m ping
ansible-playbook site.yml
```

---

### 🧭 View all groups
```bash
zGinv groups
```
Displays a list of groups and the number of servers in each.

---

## ⚖️ License

CC0 License © 2025 [Yvgeniy Goncharov](https://openbld.net)
