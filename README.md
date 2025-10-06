# Private Blog

Repositori ini adalah contoh aplikasi web kecil (mirip Telegraph) yang telah direfaktor menjadi struktur project Go yang lebih idiomatik dan menggunakan pola Repository.

Refaktor ini memisahkan tanggung jawab ke dalam paket-paket sehingga lebih mudah mengganti implementasi penyimpanan data, menambahkan pengujian, dan memperluas aplikasi.

---

## Struktur proyek

```
cmd/
  web/                # entrypoint aplikasi (main)
internal/
  handler/            # HTTP handlers + template
    handler.go
  models/             # model domain
    article.go
  repository/         # interface repository + implementasi
    repository.go
    memory.go          # implementasi in-memory
  service/            # logika bisnis
    article_service.go
go.mod
README.md
```

Kenapa susunan ini?
- `cmd/web` berisi entrypoint program. Perintah lain bisa ditambahkan nanti.
- Paket di `internal/*` bersifat privat untuk modul ini sehingga detail implementasi tidak diekspor ke pengguna paket lain.
- `repository` menerapkan pola Repository: kode bergantung pada interface, dan implementasi penyimpanan (in-memory, DB) mengimplementasikannya.

---

## Perubahan setelah refaktor

- Aplikasi single-file `main.go` dipecah menjadi beberapa paket: `models`, `repository` (interface + memory), `service`, dan `handler`.
- `service.ArticleService` menampung logika bisnis (pembuatan ID, sanitasi, operasi create/get/update/delete).
- `repository.Repository` adalah interface; `repository.MemoryRepo` adalah implementasi in-memory yang thread-safe.
- `handler` bertanggung jawab pada endpoint HTTP dan template.
- Entrypoint `cmd/web/main.go` menghubungkan repository -> service -> handler dan menjalankan server HTTP pada :8080.

---

## Cara menjalankan

Pastikan Go terinstal (disarankan Go 1.20+). Dari root proyek jalankan:

Windows (cmd.exe) - jalankan langsung:

```cmd
go run ./cmd/web
```

atau build dan jalankan executable:

```cmd
go build -o web ./cmd/web
.\web
```

Server akan berjalan di `http://localhost:8080`.

---

## Endpoint

- GET `/` — halaman editor (home)
- POST `/create` — membuat artikel
- GET `/view/{id}` — melihat artikel
- GET `/edit/{id}` — halaman edit (memerlukan kepemilikan)
- POST `/update/{id}` — memperbarui artikel (memerlukan kepemilikan)
- POST `/delete/{id}` — menghapus artikel (memerlukan kepemilikan)

Kepemilikan dilacak lewat cookie `user_id` yang dibuat pada kunjungan pertama.

---

## Catatan, keterbatasan, dan rekomendasi
