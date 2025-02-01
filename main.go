package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func logAccess(r *http.Request) {
	ip := getIP(r)
	userAgent := r.UserAgent()
	accessedURL := r.URL.Path
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	logEntry := fmt.Sprintf("Time: %s | IP: %s | User-Agent: %s | URL: %s\n", timestamp, ip, userAgent, accessedURL)

	f, err := os.OpenFile("access_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening log file:", err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(logEntry); err != nil {
		log.Println("Error writing to log file:", err)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	logAccess(r) // Ghi log khi có người truy cập vào trang chủ
	renderTemplate(w, "index", nil)
}

func projectsHandler(w http.ResponseWriter, r *http.Request) {
	// write log when someone access to projects page
	logAccess(r)
	projects := []struct {
		Title  string
		Detail string
	}{
		{"Dự án A", "Chi tiết dự án A"},
		{"Dự án B", "Chi tiết dự án B"},
	}
	renderTemplate(w, "projects", projects)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	logAccess(r) // Ghi log khi có người truy cập vào trang About
	tmpl, _ := template.ParseFiles("templates/about.html")
	tmpl.Execute(w, nil)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/projects", projectsHandler)
	http.HandleFunc("/about", aboutHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
