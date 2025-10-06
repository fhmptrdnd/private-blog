package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/fhmptrdnd/weather-api-test-web-based/internal/models"
	"github.com/fhmptrdnd/weather-api-test-web-based/internal/service"
)

type Handler struct {
    svc       *service.ArticleService
    templates *template.Template
}

func NewHandler(svc *service.ArticleService) *Handler {
    // Parse templates once and give each template a name
    t := template.New("templates")
    template.Must(t.New("home").Parse(homeTemplate))
    template.Must(t.New("view").Parse(viewTemplate))
    template.Must(t.New("edit").Parse(editTemplate))

    return &Handler{svc: svc, templates: t}
}

// getOrCreateUserID copied from original app
func getOrCreateUserID(w http.ResponseWriter, r *http.Request) string {
    cookie, err := r.Cookie("user_id")
    if err == nil && cookie.Value != "" {
        return cookie.Value
    }
    // create new
    id := generateID()
    http.SetCookie(w, &http.Cookie{
        Name:     "user_id",
        Value:    id,
        Path:     "/",
        MaxAge:   31536000 * 10,
        HttpOnly: true,
        SameSite: http.SameSiteLaxMode,
    })
    return id
}

// generateID same as in service (non-exported duplicate to avoid import cycle)
func generateID() string {
    b := make([]byte, 8)
    rand.Read(b)
    return hex.EncodeToString(b)
}

// Templates (copied from original single-file app)
const homeTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Telegraph Clone</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Georgia', serif;
            background: #f7f7f7;
            color: #333;
            line-height: 1.6;
        }

        .header {
            background: white;
            border-bottom: 1px solid #e0e0e0;
            padding: 20px 0;
            position: sticky;
            top: 0;
            z-index: 100;
        }

        .container {
            max-width: 720px;
            margin: 0 auto;
            padding: 0 20px;
        }

        .logo {
            font-size: 1.8em;
            font-weight: bold;
            color: #333;
            text-decoration: none;
        }

        .editor {
            background: white;
            margin: 40px auto;
            padding: 60px 80px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }

        input[type="text"] {
            width: 100%;
            border: none;
            font-size: 2.5em;
            font-family: 'Georgia', serif;
            margin-bottom: 20px;
            outline: none;
        }

        input[type="text"]::placeholder {
            color: #ccc;
        }

        .author-input {
            font-size: 1.1em !important;
            margin-bottom: 30px;
        }

        textarea {
            width: 100%;
            min-height: 400px;
            border: none;
            font-size: 1.2em;
            font-family: 'Georgia', serif;
            line-height: 1.8;
            resize: vertical;
            outline: none;
        }

        textarea::placeholder {
            color: #ccc;
        }

        .btn {
            background: #333;
            color: white;
            border: none;
            padding: 12px 30px;
            font-size: 16px;
            cursor: pointer;
            border-radius: 4px;
            transition: background 0.3s;
        }

        .btn:hover {
            background: #555;
        }

        .btn-container {
            text-align: right;
            margin-top: 20px;
        }

        .footer {
            text-align: center;
            padding: 40px 20px;
            color: #999;
            font-size: 0.9em;
        }

        @media (max-width: 768px) {
            .editor {
                padding: 40px 20px;
            }

            input[type="text"] {
                font-size: 1.8em;
            }

            textarea {
                font-size: 1.1em;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="container">
            <a href="/" class="logo">Telegraph</a>
        </div>
    </div>

    <div class="container">
        <div class="editor">
            <form method="POST" action="/create">
                <input type="text" name="title" placeholder="Judul" required>
                <input type="text" name="author" class="author-input" placeholder="Nama Penulis" required>
                <textarea name="content" placeholder="Ceritakan kisahmu..." required></textarea>
                
                <div class="btn-container">
                    <button type="submit" class="btn">Publikasikan</button>
                </div>
            </form>
        </div>
    </div>

    <div class="footer">
        Telegraph Clone - Buat artikel dengan mudah
    </div>
</body>
</html>
`

const viewTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Article.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Georgia', serif;
            background: #f7f7f7;
            color: #333;
            line-height: 1.6;
        }

        .header {
            background: white;
            border-bottom: 1px solid #e0e0e0;
            padding: 20px 0;
        }

        .container {
            max-width: 720px;
            margin: 0 auto;
            padding: 0 20px;
        }

        .logo {
            font-size: 1.8em;
            font-weight: bold;
            color: #333;
            text-decoration: none;
        }

        article {
            background: white;
            margin: 40px auto;
            padding: 60px 80px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }

        h1 {
            font-size: 2.5em;
            margin-bottom: 20px;
            line-height: 1.2;
        }

        .meta {
            color: #999;
            font-size: 0.95em;
            margin-bottom: 40px;
            padding-bottom: 20px;
            border-bottom: 1px solid #f0f0f0;
        }

        .content {
            font-size: 1.2em;
            line-height: 1.8;
        }

        .content p {
            margin-bottom: 1em;
        }

        .stats {
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #f0f0f0;
            color: #999;
            font-size: 0.9em;
        }

        .btn-home {
            display: inline-block;
            margin-top: 20px;
            padding: 10px 20px;
            background: #333;
            color: white;
            text-decoration: none;
            border-radius: 4px;
            transition: background 0.3s;
        }

        .btn-home:hover {
            background: #555;
        }

        .owner-actions {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #f0f0f0;
            display: flex;
            gap: 10px;
        }

        .btn-edit {
            background: #4CAF50;
            color: white;
            padding: 10px 20px;
            text-decoration: none;
            border-radius: 4px;
            transition: background 0.3s;
            display: inline-block;
        }

        .btn-edit:hover {
            background: #45a049;
        }

        .btn-delete {
            background: #f44336;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            transition: background 0.3s;
        }

        .btn-delete:hover {
            background: #da190b;
        }

        @media (max-width: 768px) {
            article {
                padding: 40px 20px;
            }

            h1 {
                font-size: 1.8em;
            }

            .content {
                font-size: 1.1em;
            }

            .owner-actions {
                flex-direction: column;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="container">
            <a href="/" class="logo">Telegraph</a>
        </div>
    </div>

    <div class="container">
        <article>
            <h1>{{.Article.Title}}</h1>
            
            <div class="meta">
                Oleh <strong>{{.Article.Author}}</strong> · {{.Article.CreatedAt.Format "2 January 2006"}}
            </div>

            <div class="content">
                {{.Article.Content}}
            </div>

            <div class="stats">
                👁️ {{.Article.Views}} tayangan
            </div>

            {{if .IsOwner}}
            <div class="owner-actions">
                <a href="/edit/{{.Article.ID}}" class="btn-edit">✏️ Edit Artikel</a>
                <form method="POST" action="/delete/{{.Article.ID}}" style="display: inline;" onsubmit="return confirm('Yakin ingin menghapus artikel ini?');">
                    <button type="submit" class="btn-delete">🗑️ Hapus Artikel</button>
                </form>
            </div>
            {{end}}

            <a href="/" class="btn-home">Buat Artikel Baru</a>
        </article>
    </div>
</body>
</html>
`

const editTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Edit - {{.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Georgia', serif;
            background: #f7f7f7;
            color: #333;
            line-height: 1.6;
        }

        .header {
            background: white;
            border-bottom: 1px solid #e0e0e0;
            padding: 20px 0;
        }

        .container {
            max-width: 720px;
            margin: 0 auto;
            padding: 0 20px;
        }

        .logo {
            font-size: 1.8em;
            font-weight: bold;
            color: #333;
            text-decoration: none;
        }

        .editor {
            background: white;
            margin: 40px auto;
            padding: 60px 80px;
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }

        .edit-label {
            color: #4CAF50;
            font-size: 0.9em;
            margin-bottom: 20px;
            display: block;
        }

        input[type="text"] {
            width: 100%;
            border: none;
            font-size: 2.5em;
            font-family: 'Georgia', serif;
            margin-bottom: 20px;
            outline: none;
        }

        .author-input {
            font-size: 1.1em !important;
            margin-bottom: 30px;
        }

        textarea {
            width: 100%;
            min-height: 400px;
            border: none;
            font-size: 1.2em;
            font-family: 'Georgia', serif;
            line-height: 1.8;
            resize: vertical;
            outline: none;
        }

        .btn {
            background: #333;
            color: white;
            border: none;
            padding: 12px 30px;
            font-size: 16px;
            cursor: pointer;
            border-radius: 4px;
            transition: background 0.3s;
            margin-right: 10px;
        }

        .btn:hover {
            background: #555;
        }

        .btn-cancel {
            background: #999;
        }

        .btn-cancel:hover {
            background: #777;
        }

        .btn-container {
            text-align: right;
            margin-top: 20px;
        }

        @media (max-width: 768px) {
            .editor {
                padding: 40px 20px;
            }

            input[type="text"] {
                font-size: 1.8em;
            }

            textarea {
                font-size: 1.1em;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="container">
            <a href="/" class="logo">Telegraph</a>
        </div>
    </div>

    <div class="container">
        <div class="editor">
            <span class="edit-label">✏️ Mode Edit</span>
            <form method="POST" action="/update/{{.ID}}">
                <input type="text" name="title" value="{{.Title}}" required>
                <input type="text" name="author" class="author-input" value="{{.Author}}" required>
                <textarea name="content" required>{{.ContentRaw}}</textarea>
                
                <div class="btn-container">
                    <a href="/view/{{.ID}}" class="btn btn-cancel">Batal</a>
                    <button type="submit" class="btn">Simpan Perubahan</button>
                </div>
            </form>
        </div>
    </div>
</body>
</html>
`

type viewData struct {
    Article *models.Article
    IsOwner bool
}

type editData struct {
    ID         string
    Title      string
    Author     string
    Content    string
    ContentRaw string
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }
    getOrCreateUserID(w, r)
    h.templates.ExecuteTemplate(w, "home", nil)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    title := r.FormValue("title")
    author := r.FormValue("author")
    content := r.FormValue("content")
    if title == "" || author == "" || content == "" {
        http.Error(w, "Semua field harus diisi", http.StatusBadRequest)
        return
    }
    owner := getOrCreateUserID(w, r)
    a, err := h.svc.Create(title, author, content, owner)
    if err != nil {
        http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+a.ID, http.StatusSeeOther)
}

func (h *Handler) View(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/view/")
    if id == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    owner := getOrCreateUserID(w, r)
    a, err := h.svc.Get(id)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    data := viewData{Article: a, IsOwner: a.OwnerID == owner}
    h.templates.ExecuteTemplate(w, "view", data)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/edit/")
    if id == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    owner := getOrCreateUserID(w, r)
    a, err := h.svc.GetNoIncrement(id)
    if err != nil {
        http.NotFound(w, r)
        return
    }
    if a.OwnerID != owner {
        http.Error(w, "Anda tidak memiliki akses untuk mengedit artikel ini", http.StatusForbidden)
        return
    }
    data := editData{
        ID:         a.ID,
        Title:      a.Title,
        Author:     a.Author,
        Content:    a.Content,
        ContentRaw: strings.ReplaceAll(a.Content, "<br>", "\n"),
    }
    h.templates.ExecuteTemplate(w, "edit", data)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    id := strings.TrimPrefix(r.URL.Path, "/update/")
    if id == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    owner := getOrCreateUserID(w, r)
    title := r.FormValue("title")
    author := r.FormValue("author")
    content := r.FormValue("content")
    if title == "" || author == "" || content == "" {
        http.Error(w, "Semua field harus diisi", http.StatusBadRequest)
        return
    }
    _, err := h.svc.Update(id, title, author, content, owner)
    if err != nil {
        http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/"+id, http.StatusSeeOther)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    id := strings.TrimPrefix(r.URL.Path, "/delete/")
    if id == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    owner := getOrCreateUserID(w, r)
    if err := h.svc.Delete(id, owner); err != nil {
        http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
