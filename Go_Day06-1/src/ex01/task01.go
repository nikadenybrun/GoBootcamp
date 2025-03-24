package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"src/ratelimiter"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

type Article struct {
	Id      int
	Name    string
	Article string
}
type Credentials struct {
	user     string
	dbname   string
	password string
}

func getCredentials() Credentials {
	file, err := os.Open("archive/admin_credentials.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var credentials Credentials
	_, err = fmt.Fscanf(file, "dbname=%s\nuser=%s\npassword=%s\n", &credentials.dbname, &credentials.user, &credentials.password)
	if err != nil {
		log.Fatal("here")
	}
	return credentials

}
func main() {
	driverName := "postgres"
	credentials := getCredentials()
	dataSourceName := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", credentials.user, credentials.dbname, credentials.password)
	var err error

	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()
	rl := &ratelimiter.RateLimiter{
		Counters:    make(map[string]int),
		Window:      time.Second,
		MaxRequests: 100,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rl.Limit(http.HandlerFunc(showPage)).ServeHTTP(w, r)
	})
	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		rl.Limit(http.HandlerFunc(postArticle)).ServeHTTP(w, r)
	})
	http.HandleFunc("/posting", func(w http.ResponseWriter, r *http.Request) {
		rl.Limit(http.HandlerFunc(addArticle)).ServeHTTP(w, r)
	})
	http.HandleFunc("/article/", func(w http.ResponseWriter, r *http.Request) {
		rl.Limit(http.HandlerFunc(showArticle)).ServeHTTP(w, r)
	})

	styleHandler := http.FileServer(http.Dir("archive/css"))
	http.Handle("/css/", http.StripPrefix("/css", styleHandler))

	imageHandler := http.FileServer(http.Dir("archive/image"))
	http.Handle("/image/", http.StripPrefix("/image", imageHandler))

	err = http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		log.Fatalln(err)
	}
}

func showPage(w http.ResponseWriter, r *http.Request) {
	var page, nextPage, prevPage int
	var limit = 3

	if strings.HasPrefix(r.URL.RawQuery, "page=") {
		num, err := strconv.Atoi(strings.TrimPrefix(r.URL.RawQuery, "page="))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		page = num
	} else {
		page = 1
	}

	total := getCountOfArticles()
	articles := getArticles(limit, limit*(page-1))

	if page*limit < total {
		nextPage = page + 1
	} else {
		nextPage = 0
	}

	prevPage = page - 1

	data := struct {
		Articles []Article
		Next     int
		Previous int
	}{Articles: articles, Next: nextPage, Previous: prevPage}

	tmpl, _ := template.ParseFiles("archive/templates/index.html")

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Fatalln(err)
	}
}
func showArticle(w http.ResponseWriter, r *http.Request) {
	ids := strings.Split(r.URL.Path, "/")[2]
	id, err := strconv.Atoi(ids)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	article, err := getArticleById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	data := struct {
		Article Article
	}{Article: article}

	tmpl, err := template.ParseFiles("archive/html/article.html")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func getCountOfArticles() int {
	query := "SELECT count(*) FROM blog"
	var total int

	err := db.QueryRow(query).Scan(&total)
	if err != nil {
		log.Fatalln(err)
	}

	return total
}
func getArticleById(id int) (Article, error) {
	query := "SELECT id, name, article FROM blog WHERE id = $1"
	row := db.QueryRow(query, id)

	if row == nil {
		return Article{}, errors.New("article not found")
	}

	article := Article{}
	err := row.Scan(&article.Id, &article.Name, &article.Article)
	if err != nil {
		return Article{}, err
	}
	return article, nil
}

func getArticles(limit, offset int) []Article {
	query := "SELECT * FROM blog ORDER BY id LIMIT " + strconv.Itoa(limit) +
		" OFFSET " + strconv.Itoa(offset)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalln(err)
	}

	defer rows.Close()

	var articles []Article
	for rows.Next() {
		article := Article{}
		rows.Scan(&article.Id, &article.Name, &article.Article)
		articles = append(articles, article)
	}

	return articles
}

func postArticle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "archive/html/admin.html")
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	article := r.FormValue("article")
	name := r.FormValue("name")
	query := "INSERT INTO blog (name, article) VALUES ($1, $2)"

	_, err := db.Exec(query, name, article)
	if err != nil {
		log.Fatalln(err)
	}

	http.ServeFile(w, r, "archive/html/posting.html")
}
