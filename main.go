package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Contact struct {
	ID string "json:id"
	Name string "json:name"
	Email string "json:email"
	Phone string "json:phone"
}

var (
	contacts = make(map[string]Contact)
	contactID = 1
	// mutual exclusion, prevents concurrent access to a ressource
	mu = sync.RWMutex{}
)

// get all contacts
func getContacts(w http.ResponseWriter, r *http.Request) {
	// lock the read lock, so no goroutines can write to the map
	mu.RLock()
	// unlocks the read lock
	defer mu.RUnlock()

	// create and fill a contact list
	var contactList []Contact
	for _, contact := range contacts {
		contactList = append(contactList, contact)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contactList)
}

func getContactById(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	mu.RLock()
	defer mu.RUnlock()

	if contact, exists := contacts[id]; exists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contact)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

// Create a new contact
func createContact(w http.ResponseWriter, r *http.Request) {
	var newContact Contact
	if err := json.NewDecoder(r.Body).Decode(&newContact); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	newContact.ID = fmt.Sprintf("%d", contactID)
	contacts[newContact.ID] = newContact
	contactID++

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

	http.HandleFunc("/contacts", getContacts)
	http.HandleFunc("/contact", getContactById)
	http.HandleFunc("/contact/create", createContact)

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}