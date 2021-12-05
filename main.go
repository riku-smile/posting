package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

// 構造体の定義
type Log struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Body  string `json:"body"`
	CTime int64  `json:"ctime"`
}

// データの保存先
const logFile = "logs.json"

func main() {
	fmt.Println("server - http://localhost:8888")

	// URL
	http.HandleFunc("/", Index)
	http.HandleFunc("/create", Create)
	// サーバーを起動
	http.ListenAndServe(":8888", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	htmlLog := ""
	logs := loadLogs()
	for _, i := range logs {
		htmlLog += fmt.Sprintf(
			"<p>(%d) <span>%s</span>: %s --- %s</p>",
			i.ID,
			html.EscapeString(i.Name),
			html.EscapeString(i.Body),
			time.Unix(i.CTime, 0).Format("2006/01/02 15:04"))
	}

	// htmlLogをIndex.htmlに渡す
	// html出力
	var tmpl *template.Template = template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, htmlLog)
	if err != nil {
		log.Fatalln(err)
	}
}

func Create(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var log Log
	log.Name = r.Form["name"][0]
	log.Body = r.Form["body"][0]
	if log.Name == "" {
		log.Name = "名無しさん"
	}
	logs := loadLogs()
	log.ID = len(logs) + 1
	log.CTime = time.Now().Unix()
	logs = append(logs, log)
	saveLogs(logs)
	http.Redirect(w, r, "/", 301)
}

func loadLogs() []Log {
	text, err := ioutil.ReadFile(logFile)
	if err != nil {
		return make([]Log, 0)
	}
	var logs []Log
	json.Unmarshal([]byte(text), &logs)
	return logs
}

func saveLogs(logs []Log) {
	bytes, _ := json.Marshal(logs)
	ioutil.WriteFile(logFile, bytes, 0644)
}
