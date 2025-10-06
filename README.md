# Private Blog

Repositori ini adalah contoh aplikasi web sederhana (mirip Telegraph) yang sudah direfaktor menjadi struktur project Go yang lebih rapi dan mengikuti pola Repository.

Tujuan refaktor: memisahkan tanggung jawab (separation of concerns), mempermudah pengujian, dan membuat titik ekstensi (mis. ganti penyimpanan ke DB) lebih mudah.

---

## Struktur proyek (singkat)

```
cmd/
  web/                # entrypoint aplikasi (main)
internal/
  handler/            # HTTP handlers + template
  models/             # model domain
  repository/         # interface repository + implementasi (memory)
  service/            # logika bisnis
go.mod
README.md
```

- `cmd/web` : `main.go` yang menginisialisasi repo/service/handler dan menjalankan server.
- `internal/models` : definisi struct `Article`.
- `internal/repository` : interface `Repository` + `MemoryRepo` (implementasi in-memory).
- `internal/service` : `ArticleService` yang mengenkapsulasi logika (create/get/update/delete).
- `internal/handler` : HTTP handlers dan template rendering.

---

## Cara menjalankan (development)

1. Pastikan Go 1.20+ terpasang.
2. Dari folder project jalankan (Windows cmd):

```cmd
go run ./cmd/web
```

atau build executable lalu jalankan:

```cmd
go build -o web ./cmd/web
.\web
```

Buka `http://localhost:8080` di browser.

---

## Endpoints & Contoh (curl)

- GET `/` — halaman editor (form HTML)
- POST `/create` — buat artikel (form POST)
- GET `/view/{id}` — lihat artikel
- GET `/edit/{id}` — tampilkan halaman edit (harus pemilik)
- POST `/update/{id}` — update artikel (harus pemilik)
- POST `/delete/{id}` — hapus artikel (harus pemilik)

Contoh membuat artikel via curl (form POST):

```bash
curl -X POST \
  -F "title=Contoh Judul" \
  -F "author=Fahmi" \
  -F "content=Isi artikel\nBaris kedua" \
  http://localhost:8080/create -v
```

Contoh hapus (POST):

```bash
curl -X POST http://localhost:8080/delete/{id} -v
```

Catatan: kepemilikan artikel ditentukan oleh cookie `user_id` yang dibuat pada kunjungan pertama. Untuk melakukan edit/delete lewat curl, sertakan cookie yang sama (browser otomatis menyimpannya).

---

## Penjelasan teknis singkat

- ID dihasilkan secara acak (hex) saat membuat artikel.
- Konten disanitasi sederhana: newline -> `<br>` (contoh minimal). Untuk produksi perlu sanitizer yang lebih baik.
- `MemoryRepo` menggunakan `sync.RWMutex` agar aman untuk akses bersamaan.

---

## Menambahkan unit test (saran)

Saya dapat membuat test untuk:

- `internal/repository` (MemoryRepo): test Create/Get/Update/Delete dan error path.
- `internal/service` (ArticleService): test Create, Get (view increment), Update (pemilik vs bukan pemilik), Delete.

Perintah menjalankan test:

```cmd
go test ./... -v
```

---

## Migrasi ke SQLite — langkah ringkas

1. Pilih driver SQLite:
   - `modernc.org/sqlite` (pure Go)
   - `github.com/mattn/go-sqlite3` (memerlukan CGO)
2. Tambah file `internal/repository/sqlite.go` yang mengimplementasikan interface `Repository`.
3. Buat tabel `articles` dengan tipe kolom sesuai model (lihat contoh SQL di bawah).
4. Di `cmd/web/main.go`, ubah inisialisasi repo menjadi SQLiteRepo jika konfigurasi/flag menginginkan persistent DB.

Contoh SQL sederhana untuk membuat tabel:

```sql
CREATE TABLE IF NOT EXISTS articles (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL,
  author TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  views INTEGER NOT NULL,
  owner_id TEXT
);
```

---

## Debugging & troubleshooting cepat

- Jika server tidak jalan: jalankan `go run ./cmd/web` dan periksa pesan error di terminal.
- Jika template error: periksa parsing template di `internal/handler/handler.go`.
- Jika cookie tidak muncul: cek konfigurasi browser/extensions yang memblok cookie.
