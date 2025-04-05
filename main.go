package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "os"
	"strconv"
	"time"
)

// Структура для хранения цитаты
type Quote struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Category string `json:"category"`
}

var quotes []Quote

const dataFile = "quotes.json"

// loadQuotes загружает цитаты из JSON-файла или инициализирует дефолтный набор, если файл отсутствует
func loadQuotes() {
	file, err := ioutil.ReadFile(dataFile)
	if err != nil {
		quotes = []Quote{
			{ID: 1, Text: "Жизнь — это то, что с тобой происходит, пока ты строишь планы.", Category: "мотивация"},
			{ID: 2, Text: "Смех продлевает жизнь.", Category: "юмор"},
			{ID: 3, Text: "Успех — это умение идти от неудачи к неудаче, не теряя энтузиазма.", Category: "мотивация"},
		}
		return
	}
	err = json.Unmarshal(file, &quotes)
	if err != nil {
		log.Println("Ошибка при загрузке цитат:", err)
		quotes = []Quote{}
	}
}

// saveQuotes сохраняет цитаты в JSON-файл
func saveQuotes() {
	data, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		log.Println("Ошибка при сохранении цитат:", err)
		return
	}
	err = ioutil.WriteFile(dataFile, data, 0644)
	if err != nil {
		log.Println("Ошибка при записи файла:", err)
	}
}

// getNextID возвращает следующий доступный ID для новой цитаты
func getNextID() int {
	maxID := 0
	for _, q := range quotes {
		if q.ID > maxID {
			maxID = q.ID
		}
	}
	return maxID + 1
}

//
// HTML-шаблоны с использованием Bootstrap для современного дизайна
//

// Главная страница ("/")
var indexTmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Генератор вдохновляющих цитат</title>
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
	<style>
		body { background: linear-gradient(to right, #fbc2eb, #a6c1ee); }
		.container { margin-top: 50px; background: rgba(255,255,255,0.9); padding: 30px; border-radius: 8px; }
		.quote-text { font-size: 1.8rem; font-weight: bold; text-align: center; margin-bottom: 10px; }
		.quote-category { font-style: italic; text-align: center; color: #666; margin-bottom: 20px; }
		.links a { margin: 0 5px; }
	</style>
</head>
<body>
<div class="container">
	<h1 class="text-center mb-4">Генератор вдохновляющих цитат</h1>
	{{if .Quote}}
		<div class="quote-text">"{{.Quote.Text}}"</div>
		<div class="quote-category">[{{.Quote.Category}}]</div>
	{{else}}
		<p class="text-center">Нет цитат для отображения (попробуйте другую категорию).</p>
	{{end}}
	<hr>
	<form class="form-inline justify-content-center mb-3" method="get" action="/">
		<input type="text" name="cat" class="form-control mr-2" placeholder="например, мотивация">
		<button type="submit" class="btn btn-success">Применить</button>
	</form>
	<div class="text-center links">
		<a href="/add" class="btn btn-primary">Добавить цитату</a>
		<a href="/export" class="btn btn-info">Экспорт (JSON)</a>
		<a href="/admin" class="btn btn-warning">Админ панель</a>
		<a href="/" class="btn btn-secondary">Обновить цитату</a>
	</div>
</div>
</body>
</html>
`))

// Страница добавления цитаты ("/add")
var addTmpl = template.Must(template.New("add").Parse(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Добавить цитату</title>
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
	<style>
		body { background: linear-gradient(to right, #fbc2eb, #a6c1ee); }
		.container { margin-top: 50px; background: rgba(255,255,255,0.9); padding: 30px; border-radius: 8px; }
	</style>
</head>
<body>
<div class="container">
	<h1 class="text-center mb-4">Добавить новую цитату</h1>
	<form method="post" action="/add">
		<div class="form-group">
			<label for="text">Текст цитаты:</label>
			<textarea name="text" id="text" class="form-control" rows="3"></textarea>
		</div>
		<div class="form-group">
			<label for="category">Категория:</label>
			<input type="text" name="category" id="category" class="form-control" placeholder="например, мотивация">
		</div>
		<button type="submit" class="btn btn-primary">Добавить</button>
		<a href="/" class="btn btn-secondary">На главную</a>
	</form>
</div>
</body>
</html>
`))

// Административная панель ("/admin")
var adminTmpl = template.Must(template.New("admin").Parse(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Админ панель - Цитаты</title>
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
</head>
<body>
<div class="container mt-4">
	<h1 class="mb-4">Административная панель</h1>
	<table class="table table-bordered">
		<thead class="thead-dark">
			<tr>
				<th>ID</th>
				<th>Текст</th>
				<th>Категория</th>
				<th>Действия</th>
			</tr>
		</thead>
		<tbody>
			{{range .Quotes}}
			<tr>
				<td>{{.ID}}</td>
				<td>{{.Text}}</td>
				<td>{{.Category}}</td>
				<td>
					<a href="/edit?id={{.ID}}" class="btn btn-sm btn-primary">Редактировать</a>
					<a href="/delete?id={{.ID}}" class="btn btn-sm btn-danger" onclick="return confirm('Удалить эту цитату?');">Удалить</a>
				</td>
			</tr>
			{{end}}
		</tbody>
	</table>
	<a href="/" class="btn btn-secondary">На главную</a>
</div>
</body>
</html>
`))

// Страница редактирования цитаты ("/edit")
var editTmpl = template.Must(template.New("edit").Parse(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Редактировать цитату</title>
	<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css">
</head>
<body>
<div class="container mt-4">
	<h1 class="mb-4">Редактировать цитату</h1>
	<form method="post" action="/edit">
		<input type="hidden" name="id" value="{{.Quote.ID}}">
		<div class="form-group">
			<label for="text">Текст цитаты:</label>
			<textarea class="form-control" name="text" id="text" rows="3">{{.Quote.Text}}</textarea>
		</div>
		<div class="form-group">
			<label for="category">Категория:</label>
			<input type="text" class="form-control" name="category" id="category" value="{{.Quote.Category}}">
		</div>
		<button type="submit" class="btn btn-primary">Сохранить изменения</button>
		<a href="/admin" class="btn btn-secondary">Отмена</a>
	</form>
</div>
</body>
</html>
`))

//
// Обработчики маршрутов
//

// indexHandler – главная страница, выбирающая случайную цитату с возможным фильтром по категории
func indexHandler(w http.ResponseWriter, r *http.Request) {
	cat := r.URL.Query().Get("cat")
	var filtered []Quote
	for _, q := range quotes {
		if cat == "" || q.Category == cat {
			filtered = append(filtered, q)
		}
	}
	var randomQuote *Quote
	if len(filtered) > 0 {
		rand.Seed(time.Now().UnixNano())
		randomQuote = &filtered[rand.Intn(len(filtered))]
	}
	data := struct {
		Quote *Quote
	}{Quote: randomQuote}
	if err := indexTmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
	}
}

// addHandler – обработка страницы добавления новой цитаты
func addHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := addTmpl.Execute(w, nil); err != nil {
			http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}
		text := r.FormValue("text")
		category := r.FormValue("category")
		if text == "" || category == "" {
			http.Error(w, "Необходимо заполнить все поля", http.StatusBadRequest)
			return
		}
		newQuote := Quote{ID: getNextID(), Text: text, Category: category}
		quotes = append(quotes, newQuote)
		saveQuotes()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// exportHandler – экспорт цитат в формате JSON
func exportHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		http.Error(w, "Ошибка экспорта", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// adminHandler – отображение административной панели со списком всех цитат
func adminHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Quotes []Quote
	}{Quotes: quotes}
	if err := adminTmpl.Execute(w, data); err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
	}
}

// editHandler – редактирование выбранной цитаты
func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		var found *Quote
		for i := range quotes {
			if quotes[i].ID == id {
				found = &quotes[i]
				break
			}
		}
		if found == nil {
			http.Error(w, "Цитата не найдена", http.StatusNotFound)
			return
		}
		data := struct {
			Quote Quote
		}{Quote: *found}
		if err := editTmpl.Execute(w, data); err != nil {
			http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка обработки формы", http.StatusBadRequest)
			return
		}
		idStr := r.FormValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}
		text := r.FormValue("text")
		category := r.FormValue("category")
		if text == "" || category == "" {
			http.Error(w, "Все поля обязательны", http.StatusBadRequest)
			return
		}
		updated := false
		for i := range quotes {
			if quotes[i].ID == id {
				quotes[i].Text = text
				quotes[i].Category = category
				updated = true
				break
			}
		}
		if !updated {
			http.Error(w, "Цитата не найдена", http.StatusNotFound)
			return
		}
		saveQuotes()
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// deleteHandler – удаление цитаты по ID
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}
	indexToDelete := -1
	for i, q := range quotes {
		if q.ID == id {
			indexToDelete = i
			break
		}
	}
	if indexToDelete == -1 {
		http.Error(w, "Цитата не найдена", http.StatusNotFound)
		return
	}
	quotes = append(quotes[:indexToDelete], quotes[indexToDelete+1:]...)
	saveQuotes()
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func main() {
	// Загружаем цитаты из JSON-файла
	loadQuotes()

	// Регистрируем маршруты
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/export", exportHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
