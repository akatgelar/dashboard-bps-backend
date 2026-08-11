# Dashboard BPS — Dashboard Backend

REST API backend untuk **Dashboard BPS (bigdata)** — menyediakan endpoint referensi **master** dan agregasi **insight** dari database `bigdata` (schema `webapi`). Dokumentasi API lengkap: `docs/api-spec.md` dan swagger (`/docs`).

---

## Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Bahasa | **Go 1.24** |
| Web framework | **Gin** (`github.com/gin-gonic/gin`) |
| Database | **PostgreSQL** (schema `webapi`) |
| ORM (endpoint master) | **GORM** (`gorm.io/driver/postgres`) |
| Query manual (endpoint insight) | **`database/sql`** + pgx (raw query) |
| Dokumentasi API | **Swagger** (`swaggo/swag`, `gin-swagger`) |
| Auto reload (dev) | **Air** (`air-verse/air`) |
| Konfigurasi | env / `godotenv` (`.env`) |

> Kebijakan: endpoint **master** memakai GORM, endpoint **insight** memakai **raw query** (bukan GORM).

---

## Prasyarat

- Go 1.24+
- PostgreSQL (database `bigdata`, schema `webapi`) yang sudah berisi data (lihat `data/webapi.sql`)
- Docker & Docker Compose (opsional)

---

## Konfigurasi Environment

Salin `.env.example` menjadi `.env` lalu isi nilainya:

```bash
cp .env.example .env
```

Variabel yang dipakai:

| Variabel | Keterangan |
|----------|------------|
| `GIN_MODE` | `debug` / `release` |
| `DB_POSTGRES_HOST` | Host database PostgreSQL |
| `DB_POSTGRES_PORT` | Port database PostgreSQL |
| `DB_POSTGRES_USER` | User database |
| `DB_POSTGRES_PASS` | Password database |
| `DB_POSTGRES_NAME` | Nama database (`bigdata`) |
| `APP_PORT` | Port HTTP server (mis. `8081`) |

---

## Menjalankan

### 1. Local (langsung)

```bash
go mod download      # install dependency
go run main.go
```

Server berjalan di `http://localhost:<APP_PORT>`.

### 2. Auto reload (Air)

```bash
go install github.com/air-verse/air@latest
air
```

### 3. Docker

```bash
cd .docker
docker compose up -d --build
```

> `.docker/docker-compose.yml` memakai `env_file` (`../.env`) sehingga konfigurasi DB & port mengikuti `.env`. Untuk produksi, pastikan `DB_POSTGRES_HOST` mengarah ke host DB yang bisa diakses container.

---

## Dokumentasi API (Swagger)

Setelah server jalan:

```
http://localhost:<APP_PORT>/docs/index.html
```

Untuk mengenerate ulang swagger setelah menambah/merubah endpoint:

```bash
swag init -g main.go -d . -o services/swagger --generatedTime=false
```

---

## Daftar Endpoint

### Master (param: `filter`, `sort`, `order`, `per_page`, `page`)

| Endpoint | Keterangan |
|----------|------------|
| `GET /master/domain` | Master domain/wilayah |
| `GET /master/subjek` | Master subjek |
| `GET /master/variable` | Master variable |
| `GET /master/variable-distinct` | Variable yang punya `variable_turunan` |
| `GET /master/variable-turunan` | Master variable turunan |
| `GET /master/variable-vertical` | Master variable vertical |
| `GET /master/tahun` | Master tahun |
| `GET /master/tahun-distinct` | Distinct `tahun_id`, `tahun_name` dari `datacontent` |
| `GET /master/tahun-turunan` | Master tahun turunan |
| `GET /master/tahun-turunan-distinct` | Distinct `turtahun_id`, `turtahun_name` dari `datacontent` |
| `GET /datacontent` | Data payload / nilai content |

Contoh query master:
```
GET /master/domain?filter=[{"field":"domain_id","operator":"eq","value":"0000"}]
GET /master/variable?sort=var_name&order=asc&per_page=20&page=1
```

Operator filter: `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `like`.

### Insight (raw query, params wajib)

| Endpoint | Params wajib | Keterangan |
|----------|--------------|------------|
| `GET /insight/big-number` | `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id` | Big number (tahun sekarang vs sebelumnya) |
| `GET /insight/per-tahun` | `domain_id`, `var_id`, `turvar_id`, `turtahun_id` | Data per wilayah & per tahun |
| `GET /insight/per-wilayah` | `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id` | Data per wilayah + klaster percentile |
| `GET /insight/per-wilayah-turvar` | `domain_id`, `var_id`, `tahun_id`, `turtahun_id` | Data per wilayah & per variable turunan |

Contoh:
```
GET /insight/per-wilayah?domain_id=0000&var_id=286&turvar_id=530&tahun_id=125&turtahun_id=0
```

---

## Struktur Direktori

```
.
├── main.go                 # entry point
├── database/               # koneksi GORM & raw *sql.DB
├── models/
│   ├── data/               # model per tabel
│   └── ...                 # base response/metadata/filter
├── routes/                 # definisi route (master, insight, index)
├── services/
│   ├── master/             # logika & query master
│   ├── insight/            # logika & query insight (raw SQL)
│   └── swagger/            # hasil generate swag (docs.go)
├── helpers/                # helpers (mis. tipe waktu)
├── docs/                   # api-spec.md, webapi.sql
└── .docker/                # Dockerfile & docker-compose.yml
```

---

## Referensi

- Spesifikasi API: `docs/api-spec.md`
- Struktur DB: `data/webapi.sql`
