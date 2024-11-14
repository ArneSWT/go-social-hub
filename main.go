package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Contact struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

var (
	db *sql.DB
	mu sync.RWMutex
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./contacts.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS contacts (
        "id" TEXT NOT NULL PRIMARY KEY,
        "name" TEXT,
        "email" TEXT,
        "phone" TEXT
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, email, phone FROM contacts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var contactList []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.Name, &contact.Email, &contact.Phone); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contactList = append(contactList, contact)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contactList)
}

func getContactById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	row := db.QueryRow("SELECT id, name, email, phone FROM contacts WHERE id = ?", id)

	var contact Contact
	if err := row.Scan(&contact.ID, &contact.Name, &contact.Email, &contact.Phone); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Contact not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contact)
}

func createContact(w http.ResponseWriter, r *http.Request) {
	var newContact Contact
	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newContact.ID = uuid.New().String()
	_, err := db.Exec("INSERT INTO contacts (id, name, email, phone) VALUES (?, ?, ?, ?)",
		newContact.ID, newContact.Name, newContact.Email, newContact.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newContact)
}

func main() {

	/*
		CRUD operations:

		GET /contacts: Fetch all contacts.
		GET /contact?id=<id>: Fetch a single contact by its ID.
		POST /contact/create: Create a new contact.

	*/

	initDB()
	defer db.Close()

	http.HandleFunc("/contacts", getContacts)
	http.HandleFunc("/contact", getContactById)
	http.HandleFunc("/contact/create", createContact)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
