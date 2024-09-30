package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var (
	receipts          = sync.Map{}
	alphanumericRegex = regexp.MustCompile(`[a-zA-Z0-9]`)
)

func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid receipt format", http.StatusBadRequest)
		return
	}

	// Check for no items
	if len(receipt.Items) == 0 {
		http.Error(w, "Receipt must contain at least one item", http.StatusBadRequest)
		return
	}

	// Check if purchase date is valid
	_, err = time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		http.Error(w, "Invalid purchase date format", http.StatusBadRequest)
		return
	}

	// Check if purchase time is valid
	_, err = time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		http.Error(w, "Invalid purchase time format", http.StatusBadRequest)
		return
	}

	// Check for empty item descriptions
	for _, item := range receipt.Items {
		if strings.TrimSpace(item.ShortDescription) == "" {
			http.Error(w, "Item description cannot be empty", http.StatusBadRequest)
			return
		}
	}

	// If everything is valid, process the receipt
	id := uuid.New().String()
	points := calculatePoints(receipt)
	receipts.Store(id, points)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	value, exists := receipts.Load(id)
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := value.(int)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"points": points})
}

func calculatePoints(receipt Receipt) int {
	points := 0
	points += len(alphanumericRegex.FindAllString(receipt.Retailer, -1))

	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil || total <= 0 {
		return points // safely return points calculated so far
	}

	if math.Mod(total, 1) == 0 {
		points += 50
	}
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		descriptionLength := len(strings.TrimSpace(item.ShortDescription))
		if descriptionLength == 0 {
			continue // skip empty descriptions but continue processing
		}
		if descriptionLength%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil || price <= 0 {
				continue // skip invalid prices
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && date.Day()%2 != 0 {
		points += 6
	}

	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil && purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/receipts/process", processReceiptHandler).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", getPointsHandler).Methods("GET")

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
