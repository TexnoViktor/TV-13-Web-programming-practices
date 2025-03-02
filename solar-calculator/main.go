package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
)

// Дані для розрахунку від клієнта
type Request struct {
	Pc     float64 `json:"pc"`
	Sigma1 float64 `json:"sigma1"`
	Sigma2 float64 `json:"sigma2"`
	Price  float64 `json:"price"`
}

// Результати розрахунку
type Response struct {
	DeltaW1       float64 `json:"deltaW1"`
	DeltaW2       float64 `json:"deltaW2"`
	ProfitBefore  float64 `json:"profitBefore"`  // Додано
	PenaltyBefore float64 `json:"penaltyBefore"` // Додано
	ProfitAfter   float64 `json:"profitAfter"`   // Додано
	PenaltyAfter  float64 `json:"penaltyAfter"`  // Додано
	TotalProfit   float64 `json:"totalProfit"`
}

// Функція помилок (error function)
func erf(x float64) float64 {
	// Апроксимація з використанням формули
	a1 := 0.254829592
	a2 := -0.284496736
	a3 := 1.421413741
	a4 := -1.453152027
	a5 := 1.061405429
	p := 0.3275911

	t := 1.0 / (1.0 + p*math.Abs(x))
	y := 1.0 - (((((a5*t+a4)*t)+a3)*t+a2)*t+a1)*t*math.Exp(-x*x)

	if x < 0 {
		return -y
	}
	return y
}

// Нормальний CDF
func normalCDF(x, mean, sigma float64) float64 {
	return 0.5 * (1 + erf((x-mean)/(sigma*math.Sqrt(2))))
}

// Розрахунок частки без небалансів
func calculateDeltaW(sigma, pc float64) float64 {
	lower := pc * 0.95
	upper := pc * 1.05
	return normalCDF(upper, pc, sigma) - normalCDF(lower, pc, sigma)
}

// Обробник розрахунків
func calculateHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невірний формат даних", http.StatusBadRequest)
		return
	}

	deltaW1 := calculateDeltaW(req.Sigma1, req.Pc)
	deltaW2 := calculateDeltaW(req.Sigma2, req.Pc)

	// Розрахунок прибутків та штрафів
	profitBefore := req.Pc * 24 * deltaW1 * req.Price * 1000
	penaltyBefore := req.Pc * 24 * (1 - deltaW1) * req.Price * 1000

	profitAfter := req.Pc * 24 * deltaW2 * req.Price * 1000
	penaltyAfter := req.Pc * 24 * (1 - deltaW2) * req.Price * 1000

	totalProfit := (profitAfter - penaltyAfter) - (profitBefore - penaltyBefore)

	response := Response{
		DeltaW1:       deltaW1 * 100,
		DeltaW2:       deltaW2 * 100,
		ProfitBefore:  profitBefore,
		PenaltyBefore: penaltyBefore,
		ProfitAfter:   profitAfter,
		PenaltyAfter:  penaltyAfter,
		TotalProfit:   totalProfit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Сервіс статичних файлів
func serveStatic(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./static"))
	fs.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", serveStatic)
	http.HandleFunc("/calculate", calculateHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Сервер запущено на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
