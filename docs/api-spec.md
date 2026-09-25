# PRD — API Spec Dashboard BPS

- **Project**: Dashboard BPS (bigdata)
- **Sumber struktur DB**: `data/webapi.sql`
- **Status**: Implemented / v2.0
- **Tanggal**: 2026-08-11

---

## 1. Ringkasan

Dokumen ini adalah spesifikasi API (PRD) untuk layanan **Dashboard BPS**. API dibagi menjadi dua kelompok besar:

1. **Master** — endpoint referensi data master (domain, subject, variable, variable turunan, variable vertical, tahun, tahun turunan) dan datacontent, plus endpoint `*-distinct` dan `variable-distinct`.
2. **Insight** — endpoint agregasi/analisis data untuk kebutuhan tampilan (big number, per tahun, per wilayah, per wilayah-turvar).

Semua endpoint mengembalikan **response envelope** yang konsisten (`status`, `data`, `metadata`). Seluruh response (termasuk error & route-not-found) berbentuk **JSON**.

> Catatan: spesifikasi ini disesuaikan dengan implementasi final hasil diskusi/iterasi. Perubahan utama dari draft v1: param filter memakai `filter` (bukan `where`), sorting memakai `sort` + `order` terpisah, format datetime `yyyy-MM-dd HH:mm:ss` (tanpa `T`/`Z`), dan tambahan endpoint-distinct serta endpoint insight baru.

---

## 1.1 Tech Stack

| Komponen | Teknologi | Keterangan |
|----------|-----------|------------|
| Bahasa pemrograman | **Go (Golang)** | Versi 1.21+ |
| Web framework | **Gin** | `github.com/gin-gonic/gin` |
| Database | **PostgreSQL** | Basis data `bigdata`, schema `webapi` |
| ORM (endpoint master) | **GORM** | `gorm.io/gorm` + `gorm.io/driver/postgres` |
| Query manual (endpoint insight) | **`database/sql`** (driver pgx) | Query SQL ditulis manual untuk agregasi/analisis |
| Koneksi raw DB | `database.DB_SQL_POSTGRES` | `*sql.DB` untuk query raw insight |
| Konfigurasi | **env / dotenv** | Kredensial DB, host, port via environment variable |

> Kebijakan: endpoint **master** memakai ORM (GORM), endpoint **insight** memakai **raw query** (`database/sql`), BUKAN GORM.

---

## 2. Konvensi Umum

### 2.1 Base URL

Tanpa prefix `/api`. Contoh:

```
https://<host>/master/domain
https://<host>/insight/big-number
```

Base host menyesuaikan environment (dev / staging / prod).

### 2.2 Format Data & Kode

- HTTP method: seluruh endpoint memakai **GET** (operasi read-only).
- Response content-type: `application/json`.
- Semua output datetime (`get_at`, `last_updated_at`, `last_update_data`, `last_update_pipeline`) memakai format **`yyyy-MM-dd HH:mm:ss`** (tanpa `T` dan `Z`).

### 2.3 Response Envelope Umum

```json
{
  "status": true,
  "data": [],
  "metadata": {}
}
```

| Field      | Tipe | Keterangan |
|------------|------|------------|
| `status`   | bool | `true` sukses, `false` gagal |
| `data`     | array \| object | Hasil data (array untuk list, object untuk single/agregat) |
| `metadata` | object | Informasi tambahan (lihat per kelompok API) |

HTTP status: `200` sukses, `400` bad request / param tidak valid, `404` route-not-found, `500` error server. Error & route-not-found dikembalikan sebagai JSON dengan bentuk envelope yang sama (`status: false`).

---

## 3. API Master

### 3.1 Endpoint

| Endpoint | Keterangan |
|----------|------------|
| `master/domain` | Master domain/wilayah |
| `master/subjek` | Master subjek (kategori indikator) |
| `master/variable` | Master variable/indikator |
| `master/variable-distinct` | Variable yang memiliki `variable_turunan` |
| `master/variable-turunan` | Master variable turunan |
| `master/variable-vertical` | Master variable vertical (wilayah/tipe) |
| `master/tahun` | Master periode/tahun |
| `master/tahun-distinct` | Distinct `tahun_id`, `tahun_name` dari `datacontent` |
| `master/tahun-turunan` | Master periode turunan |
| `master/tahun-turunan-distinct` | Distinct `turtahun_id`, `turtahun_name` dari `datacontent` |
| `master/tahun-turunan-group` | Group periode turunan (`group_turth_id`, `group_turth_name`) |
| `master/tahun-turunan-group-distinct` | Distinct group periode turunan dari `datacontent` |
| `datacontent` | Data nilai content (payload data) |

### 3.2 Parameter Umum (Query) — Endpoint Master

Semua endpoint master mendukung param berikut:

| Param | Contoh | Tipe | Keterangan |
|-------|--------|------|------------|
| `filter` | `filter=[{"field":"domain_id","operator":"eq","value":"0000"}]` | string (URL-encoded JSON array) | Array filter `{"field", "operator", "value"}`. Operator: `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `like`. Field mengikuti nama kolom DB. Bisa berisi beberapa filter. |
| `sort` | `sort=domain_name` | string | Nama kolom pengurutan (tunggal). |
| `order` | `order=desc` | string | Arah urutan: `asc` (default) / `desc`. |
| `per_page` | `per_page=10` | number | Jumlah data per halaman. Default 20. |
| `page` | `page=1` | number | Nomor halaman (1-based). |

Contoh filter:
```
GET /master/domain?filter=[{"field":"domain_id","operator":"eq","value":"0000"}]
GET /master/domain?filter=[{"field":"domain_name","operator":"like","value":"Aceh"}]
GET /master/tahun?filter=[{"field":"tahun_name","operator":"gt","value":"2020"}]
```

### 3.3 Payload Response Master

```json
{
  "status": true,
  "data": [],
  "metadata": {
    "last_update_data": "2026-07-22 11:21:33",
    "last_update_pipeline": "2026-07-22 11:21:33",
    "total_data": 549,
    "total_page": 28,
    "per_page": 20,
    "page": 1
  }
}
```

| Field metadata | Tipe | Keterangan |
|----------------|------|------------|
| `last_update_data` | string | Waktu update data terakhir (`yyyy-MM-dd HH:mm:ss`). |
| `last_update_pipeline` | string | Waktu update pipeline terakhir (`yyyy-MM-dd HH:mm:ss`). |
| `total_data` | number | Total seluruh data. |
| `total_page` | number | Total jumlah halaman. |
| `per_page` | number | Jumlah data per halaman yang diminta. |
| `page` | number | Halaman saat ini. |

> `last_update_pipeline` = `MAX(get_at)` seluruh tabel sumber. `last_update_data` = `MAX(get_at)` untuk seluruh master; untuk `datacontent` dan endpoint-`*distinct`, = `MAX(last_updated_at)` **sesuai filter** yang dikirim.

### 3.4 Detail Endpoint Master

#### 3.4.1 `master/domain`

Mengambil daftar master domain (wilayah BPS).

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `domain_name` | string | Nama domain/wilayah |
| `domain_url` | string | URL website BPS terkait |
| `get_at` | string | Waktu pengambilan data |
| `is_active` | bool | Aktif / non-aktif |

```
GET /master/domain
GET /master/domain?filter=[{"field":"domain_id","operator":"eq","value":"0000"}]
GET /master/domain?sort=domain_name&order=asc&per_page=10&page=1
```

#### 3.4.2 `master/subjek`

Mengambil daftar master subjek indikator.

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `sub_id` | string | Kode subject |
| `sub_name` | string | Nama subject |
| `get_at` | string | Waktu pengambilan |

```
GET /master/subjek
GET /master/subjek?filter=[{"field":"sub_name","operator":"like","value":"Ekonomi"}]
```

#### 3.4.3 `master/variable`

Mengambil daftar master variable/indikator.

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `var_id` | string | Kode variable |
| `var_name` | string | Nama variable |
| `sub_id` | string | Kode subject |
| `sub_name` | string | Nama subject |
| `notes` | string | Catatan/sumber |
| `unit` | string | Satuan |
| `get_at` | string | Waktu pengambilan |

```
GET /master/variable
GET /master/variable?filter=[{"field":"var_id","operator":"eq","value":"70"}]
GET /master/variable?sort=var_name&order=asc&per_page=50
```

#### 3.4.4 `master/variable-distinct`

Sama dengan `master/variable` (params & isi sama), tetapi **hanya menampilkan variable yang memiliki `variable_turunan`**. Relasi: `var_id` dari `master_variable` harus ada di `master_variable_turunan`.

```
GET /master/variable-distinct
GET /master/variable-distinct?filter=[{"field":"domain_id","operator":"eq","value":"0000"}]
```

#### 3.4.5 `master/variable-turunan`

Mengambil daftar master variable turunan.

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `var_id` | string | Kode variable |
| `var_name` | string | Nama variable |
| `turvar_id` | string | Kode variable turunan |
| `turvar_name` | string | Nama variable turunan |
| `get_at` | string | Waktu pengambilan |

```
GET /master/variable-turunan
GET /master/variable-turunan?filter=[{"field":"var_id","operator":"eq","value":"10"}]
```

#### 3.4.6 `master/variable-vertical`

Mengambil daftar master variable vertical (wilayah/tipe penyajian).

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `var_id` | string | Kode variable |
| `var_name` | string | Nama variable |
| `vervar_id` | string | Kode variable vertical |
| `vervar_name` | string | Nama variable vertical |
| `get_at` | string | Waktu pengambilan |

```
GET /master/variable-vertical
GET /master/variable-vertical?filter=[{"field":"vervar_id","operator":"eq","value":"1106"}]
```

#### 3.4.7 `master/tahun`

Mengambil daftar master periode/tahun.

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `tahun_id` | string | Kode tahun |
| `tahun_name` | string | Tahun (label) |
| `get_at` | string | Waktu pengambilan |
| `is_active` | bool | Aktif / non-aktif |

```
GET /master/tahun
GET /master/tahun?filter=[{"field":"tahun_name","operator":"eq","value":"2026"}]
GET /master/tahun?sort=tahun_name&order=desc
```

#### 3.4.8 `master/tahun-distinct`

Distinct `tahun_id`, `tahun_name` dari `webapi.datacontent`. Params sama dgn endpoint master (`filter`, `sort`, `order`, `per_page`, `page`). Field filter/sort yang didukung: `tahun_id`, `tahun_name`, `domain_id`, `var_id`.

```
GET /master/tahun-distinct
GET /master/tahun-distinct?filter=[{"field":"domain_id","operator":"eq","value":"0000"},{"field":"var_id","operator":"eq","value":"286"}]
```

Response `data` (array):
```
[ { "tahun_id": "125", "tahun_name": "2025" } ]
```

#### 3.4.9 `master/tahun-turunan`

Mengambil daftar master periode turunan.

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `turtahun_id` | string | Kode periode turunan |
| `turtahun_name` | string | Nama periode turunan |
| `group_turth_id` | string | Kode group periode |
| `group_turth_name` | string | Nama group periode |
| `get_at` | string | Waktu pengambilan |

```
GET /master/tahun-turunan
GET /master/tahun-turunan?filter=[{"field":"group_turth_name","operator":"eq","value":"Tahunan"}]
```

#### 3.4.10 `master/tahun-turunan-distinct`

Distinct `turtahun_id`, `turtahun_name` dari `webapi.datacontent`. Params sama dgn endpoint master. Field filter/sort yang didukung: `turtahun_id`, `turtahun_name`, `domain_id`, `var_id`.

```
GET /master/tahun-turunan-distinct
GET /master/tahun-turunan-distinct?filter=[{"field":"domain_id","operator":"eq","value":"0000"},{"field":"var_id","operator":"eq","value":"102"}]
```

Response `data` (array):
```
[ { "turtahun_id": "0", "turtahun_name": "Tahun" } ]
```

#### 3.4.11 `master/tahun-turunan-group`

Group periode dari `webapi.master_tahun_turunan`, di-resolve lewat INNER JOIN ke `webapi.datacontent` supaya bisa difilter `domain_id` & `var_id`. Field filter/sort: `turtahun_group_id`, `turtahun_group_name` (alias dari `group_turth_id`/`group_turth_name`), `group_turth_id`, `group_turth_name`, `domain_id`, `var_id`, `turtahun_id`, `turtahun_name`.

Varian ini **non-distinct**: satu baris per baris `master_tahun_turunan` (versi `get_at`), sehingga bisa memuat group yang sama lebih dari sekali — sama seperti endpoint master `tahun-turunan`.

```
GET /master/tahun-turunan-group?filter=[{"field":"domain_id","operator":"eq","value":"0000"},{"field":"var_id","operator":"eq","value":"1"}]
```

Response `data` (array):
```
[ { "turtahun_group_id": "1", "turtahun_group_name": "Bulanan" } ]
```

#### 3.4.12 `master/tahun-turunan-group-distinct`

Distinct group periode dari `webapi.datacontent` ⋈ `webapi.master_tahun_turunan`. Field filter/sort: `turtahun_group_id`, `turtahun_group_name`, `group_turth_id`, `group_turth_name`, `domain_id`, `var_id`, `turtahun_id`, `turtahun_name`.

```
GET /master/tahun-turunan-group-distinct?filter=[{"field":"domain_id","operator":"eq","value":"0000"},{"field":"var_id","operator":"eq","value":"1"}]
```

Response `data` (array):
```
[ { "turtahun_group_id": "1", "turtahun_group_name": "Bulanan" } ]
```

#### 3.4.13 `datacontent`

Mengambil daftar isi data (nilai angka) hasil pipeline. Mengikuti parameter master (`filter`, `sort`, `order`, `per_page`, `page`).

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | Primary key |
| `domain_id` | string | Kode domain |
| `var_id` | string | Kode variable |
| `var_name` | string | Nama variable |
| `unit` | string | Satuan variable (kolom `var_unit`) |
| `sub_name` | string | Nama subject (kolom `var_subject`) |
| `turvar_id` | string | Kode variable turunan |
| `turvar_name` | string | Nama variable turunan |
| `tahun_id` | string | Kode tahun |
| `tahun_name` | string | Label tahun |
| `turtahun_id` | string | Kode tahun turunan |
| `turtahun_name` | string | Label tahun turunan |
| `vervar_id` | string | Kode variable vertical |
| `vervar_name` | string | Nama variable vertical |
| `datacontent_id` | string | Kode content |
| `datacontent_value` | number | Nilai data |
| `last_updated_at` | string | Waktu update data terakhir |
| `get_at` | string | Waktu pengambilan |

Contoh:
```
GET /datacontent
GET /datacontent?filter=[{"field":"var_name","operator":"eq","value":"...PDRB..."}]
GET /datacontent?filter=[{"field":"domain_id","operator":"eq","value":"0000"}]&sort=datacontent_value&order=desc&per_page=10&page=2
```

---

## 4. API Insight

### 4.1 Parameter Wajib

Param wajib per endpoint insight (semuanya kode, diambil dari `webapi.datacontent` / `master_tahun`):

| Endpoint | Params wajib |
|----------|--------------|
| `insight/big-number` | `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id` |
| `insight/per-tahun` | `domain_id`, `var_id`, `turvar_id`, `turtahun_id` |
| `insight/per-wilayah` | `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id` |
| `insight/per-wilayah-turvar` | `domain_id`, `var_id`, `tahun_id`, `turtahun_id` |

Missing param → `400`.

### 4.2 Payload Response Insight

```json
{
  "status": true,
  "data": {},
  "metadata": {
    "last_update_data": "2026-02-19 10:24:06",
    "last_update_pipeline": "2026-07-24 11:54:35"
  }
}
```

| Field metadata | Tipe | Keterangan |
|----------------|------|------------|
| `last_update_data` | string | `MAX(last_updated_at)` dari baris yang cocok (`yyyy-MM-dd HH:mm:ss`). |
| `last_update_pipeline` | string | `MAX(get_at)` dari baris yang cocok (`yyyy-MM-dd HH:mm:ss`). |

Endpoint insight TIDAK memakai pagination.

### 4.3 Detail Endpoint Insight

#### 4.3.1 `insight/big-number`

Angka agregat variable per tahun, dibandingkan dengan tahun sebelumnya.

Params: `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id` (semua wajib).

Response `data` (object):

| Field | Keterangan |
|-------|------------|
| `tahun_sekarang_id` | Kode tahun sekarang |
| `tahun_sekarang_name` | Label tahun sekarang |
| `tahun_sekarang_data.max_value` | Nilai tertinggi (excl. `vervar_name = Indonesia`) |
| `tahun_sekarang_data.max_name` | Nama wilayah nilai tertinggi |
| `tahun_sekarang_data.min_value` | Nilai terendah (excl. Indonesia) |
| `tahun_sekarang_data.min_name` | Nama wilayah nilai terendah |
| `tahun_sekarang_data.avg_value` | Rata-rata (excl. Indonesia) |
| `tahun_sekarang_data.median_value` | Median (excl. Indonesia) |
| `tahun_sekarang_data.indo_value` | Nilai Indonesia (`vervar_name = Indonesia`) |
| `tahun_sekarang_data.jabar_value` | Nilai Jawa Barat (`vervar_id = 3200`) |
| `tahun_sebelumnya_id` | Kode tahun sebelumnya (`tahun_name - 1`) |
| `tahun_sebelumnya_name` | Label tahun sebelumnya |
| `tahun_sebelumnya_data.avg_value` | Rata-rata tahun sebelumnya |
| `tahun_sebelumnya_data.median_value` | Median tahun sebelumnya |
| `tahun_sebelumnya_data.indo_value` | Nilai Indonesia tahun sebelumnya |
| `tahun_sebelumnya_data.jabar_value` | Nilai Jawa Barat tahun sebelumnya |
| `tahun_sebelumnya_data.avg_persen` | `((avg_sekarang - avg_sebelumnya)/avg_sebelumnya) * 100` |
| `tahun_sebelumnya_data.indo_persen` | `((indo_sekarang - indo_sebelumnya)/indo_sebelumnya) * 100` |

> `median_value` dihitung dengan `percentile_cont(0.5)`, memakai filter yang sama dengan `avg_value` (excl. `vervar_name = Indonesia`). `jabar_value` diambil dari baris `vervar_id = 3200` (Jawa Barat) — rumusnya sama dengan `indo_value`, hanya beda kode. Nilai `null` bila barisnya tidak ada (mis. `indo_value` di domain `3200`).

```
GET /insight/big-number?domain_id=0000&var_id=286&turvar_id=530&tahun_id=125&turtahun_id=0
```

```json
{
  "status": true,
  "data": {
    "tahun_sekarang_id": "125",
    "tahun_sekarang_name": "2025",
    "tahun_sekarang_data": {
      "max_value": 3926153.3,
      "max_name": "DKI JAKARTA",
      "min_value": 28377.77,
      "min_name": "PAPUA PEGUNUNGAN",
      "avg_value": 622000.6855263158,
      "median_value": 249279.065,
      "indo_value": 23821103.6,
      "jabar_value": 3038667.95
    },
    "tahun_sebelumnya_id": "124",
    "tahun_sebelumnya_name": "2024",
    "tahun_sebelumnya_data": {
      "avg_value": 579549.7536842105,
      "median_value": 233033.325,
      "indo_value": 22138990.8,
      "jabar_value": 2823451.79,
      "avg_persen": 7.324812334444773,
      "indo_persen": 7.597965124950505
    }
  },
  "metadata": {
    "last_update_data": "2026-02-19 10:24:06",
    "last_update_pipeline": "2026-07-24 11:54:35"
  }
}
```

> Tahun sebelumnya ditentukan via `master_tahun` (`domain_id` + `tahun_name - 1`). Jika tak ada data tahun sebelumnya, field-nya `null`.

#### 4.3.2 `insight/per-tahun`

Nilai data per wilayah (`vervar`) & per tahun. Params: `domain_id`, `var_id`, `turvar_id`, `turtahun_id` (tanpa `tahun_id`).

Response `data` (array), dikelompokkan per wilayah, tiap wilayah berisi daftar tahun:

```json
[
  {
    "vervar_id": "1100",
    "vervar_name": "ACEH",
    "data": [
      { "tahun_id": "125", "tahun_name": "2025",
        "datacontent_id": "11002865301250", "datacontent_value": 257502.43 }
    ]
  }
]
```

```
GET /insight/per-tahun?domain_id=0000&var_id=286&turvar_id=530&turtahun_id=0
```

#### 4.3.3 `insight/per-wilayah`

Nilai data per wilayah + klaster (range) berdasarkan percentile `datacontent_value`. Params: `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id`.

Response `data` (object):

| Field | Keterangan |
|-------|------------|
| `value` | Array per wilayah: `{vervar_id, vervar_name, datacontent_id, datacontent_value, value, color}` |
| `range` | Array klaster: `{from, to, color, total_cluster}` |

Aturan klaster:
- Banyak klaster berdasarkan banyak nilai unik, **min 2 / max 5**.
- Batas `from`/`to` dihitung dari distribusi percentile (pembagian peringkat/quantile) → tiap wilayah masuk tepat 1 klaster.
- `color` diinterpolasi dari `#FFFFFF` (klaster terendah) ke `#4856FF` (klaster tertinggi).
- `total_cluster` = jumlah wilayah dalam rentang tsb.
- `value` wilayah = `datacontent_value` (nilai aslinya); `color` = warna klaster wilayah tsb.

```
GET /insight/per-wilayah?domain_id=0000&var_id=286&turvar_id=530&tahun_id=125&turtahun_id=0
```

```json
{
  "status": true,
  "data": {
    "value": [
      { "vervar_id": "1100", "vervar_name": "ACEH",
        "datacontent_id": "11002865301250", "datacontent_value": 257502.43,
        "value": 257502.43, "color": "#A4ABFF" }
    ],
    "range": [
      { "from": 28377.77, "to": 81807.57, "color": "#FFFFFF", "total_cluster": 7 },
      { "from": 90513.18, "to": 192951.54, "color": "#D1D5FF", "total_cluster": 8 },
      { "from": 204333.57, "to": 328569.51, "color": "#A4ABFF", "total_cluster": 8 },
      { "from": 349662.6, "to": 889072.55, "color": "#7680FF", "total_cluster": 8 },
      { "from": 936203.98, "to": 23821103.6, "color": "#4856FF", "total_cluster": 8 }
    ]
  },
  "metadata": {
    "last_update_data": "2026-02-19 10:24:06",
    "last_update_pipeline": "2026-07-24 11:54:35"
  }
}
```

#### 4.3.4 `insight/per-wilayah-turvar`

Nilai data per wilayah (`vervar`) & per variable turunan (`turvar`). Params: `domain_id`, `var_id`, `tahun_id`, `turtahun_id` (tanpa `turvar_id`).

Response `data` (array), dikelompokkan per wilayah, tiap wilayah berisi daftar turvar:

```json
[
  {
    "vervar_id": "1100",
    "vervar_name": "ACEH",
    "data": [
      { "turvar_id": "530", "turvar_name": "Harga Berlaku",
        "datacontent_id": "11002865301250", "datacontent_value": 257502.43 },
      { "turvar_id": "531", "turvar_name": "Harga Konstan 2010",
        "datacontent_id": "11002865311250", "datacontent_value": 158340.39 }
    ]
  }
]
```

```
GET /insight/per-wilayah-turvar?domain_id=0000&var_id=286&tahun_id=125&turtahun_id=0
```

---

## 5. Ringkasan Endpoint

### API Master (param: `filter`, `sort`, `order`, `per_page`, `page`)

| Endpoint | Payload `data` |
|----------|----------------|
| `master/domain` | array |
| `master/subjek` | array |
| `master/variable` | array |
| `master/variable-distinct` | array |
| `master/variable-turunan` | array |
| `master/variable-vertical` | array |
| `master/tahun` | array |
| `master/tahun-distinct` | array |
| `master/tahun-turunan` | array |
| `master/tahun-turunan-distinct` | array |
| `datacontent` | array |

### API Insight (raw query)

| Endpoint | Params wajib | Payload `data` |
|----------|--------------|----------------|
| `insight/big-number` | `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id` | object `{}` |
| `insight/per-tahun` | `domain_id`, `var_id`, `turvar_id`, `turtahun_id` | array `[]` |
| `insight/per-wilayah` | `domain_id`, `var_id`, `turvar_id`, `tahun_id`, `turtahun_id` | object `{}` (value + range) |
| `insight/per-wilayah-turvar` | `domain_id`, `var_id`, `tahun_id`, `turtahun_id` | array `[]` |

---

## 6. Contoh Request-Response Lengkap

### 6.1 Master

```bash
# Ambil domain, urut alfabet, pagination
curl "https://<host>/master/domain?sort=domain_name&order=asc&per_page=10&page=1"

# Filter domain_id = 0000
curl "https://<host>/master/domain?filter=%5B%7B%22field%22%3A%22domain_id%22%2C%22operator%22%3A%22eq%22%2C%22value%22%3A%220000%22%7D%5D"

# Filter datacontent by variable + sort value desc
curl "https://<host>/datacontent?filter=%5B%7B%22field%22%3A%22domain_id%22%2C%22operator%22%3A%22eq%22%2C%22value%22%3A%220000%22%7D%5D&sort=datacontent_value&order=desc&per_page=20&page=1"
```

### 6.2 Insight

```bash
# Big number
curl "https://<host>/insight/big-number?domain_id=0000&var_id=286&turvar_id=530&tahun_id=125&turtahun_id=0"

# Per tahun
curl "https://<host>/insight/per-tahun?domain_id=0000&var_id=286&turvar_id=530&turtahun_id=0"

# Per wilayah (dgn klaster)
curl "https://<host>/insight/per-wilayah?domain_id=0000&var_id=286&turvar_id=530&tahun_id=125&turtahun_id=0"

# Per wilayah - turvar
curl "https://<host>/insight/per-wilayah-turvar?domain_id=0000&var_id=286&tahun_id=125&turtahun_id=0"
```

---

## 7. Catatan / Asumsi

1. Semua endpoint bersifat read-only (GET).
2. Endpoint **master** pakai GORM; endpoint **insight** pakai raw query (`database/sql`).
3. Datetime (`get_at`, `last_updated_at`, metadata) diformat `yyyy-MM-dd HH:mm:ss`, tanpa `T`/`Z`.
4. Metadata:
   - `last_update_pipeline` = `MAX(get_at)`.
   - `last_update_data` = `MAX(get_at)` (master) atau `MAX(last_updated_at)` sesuai filter (datacontent & endpoint-distinct).
5. Param filter master memakai `filter` (array) — bukan `where`. Operator: `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `like`.
6. Sorting memakai `sort` (satu kolom) + `order` (asc/desc), bukan `sort=kolom:asc`.
7. `master/variable-distinct` = variable yang memiliki `variable_turunan`.
8. `insight/big-number`: tahun sebelumnya = `master_tahun.tahun_name - 1`; `avg_persen`/`indo_persen` = perbandingan tahun sekarang terhadap tahun sebelumnya.
9. `insight/per-wilayah`: klaster percentile, min 2 / max 5, warna `#FFFFFF` → `#4856FF`.
