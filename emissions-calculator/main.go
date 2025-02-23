package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type CalculationInput struct {
	CoalAmount       float64 `json:"coalAmount"`
	FuelOilAmount    float64 `json:"fuelOilAmount"`
	NaturalGasAmount float64 `json:"naturalGasAmount"`
}

type CalculationResult struct {
	CoalEmission       EmissionResult `json:"coalEmission"`
	FuelOilEmission    EmissionResult `json:"fuelOilEmission"`
	NaturalGasEmission EmissionResult `json:"naturalGasEmission"`
}

type EmissionResult struct {
	EmissionFactor float64 `json:"emissionFactor"`
	TotalEmission  float64 `json:"totalEmission"`
}

func calculateEmissions(input CalculationInput) CalculationResult {
	// Constants for calculation
	const (
		coalHeatValue       = 20.47 // MJ/kg
		fuelOilHeatValue    = 40.40 // MJ/kg
		naturalGasHeatValue = 33.08 // MJ/m3
		filterEfficiency    = 0.985
	)

	// Coal emission calculation
	coalEmissionFactor := (25.20 * 0.8 * 1000) / coalHeatValue * (1 - filterEfficiency)
	coalTotalEmission := coalEmissionFactor * input.CoalAmount * coalHeatValue / 1000000

	// Fuel oil emission calculation
	fuelOilEmissionFactor := (0.15 * 1.0 * 1000) / fuelOilHeatValue * (1 - filterEfficiency)
	fuelOilTotalEmission := fuelOilEmissionFactor * input.FuelOilAmount * fuelOilHeatValue / 1000000

	// Natural gas has no particulate emissions
	return CalculationResult{
		CoalEmission: EmissionResult{
			EmissionFactor: coalEmissionFactor,
			TotalEmission:  coalTotalEmission,
		},
		FuelOilEmission: EmissionResult{
			EmissionFactor: fuelOilEmissionFactor,
			TotalEmission:  fuelOilTotalEmission,
		},
		NaturalGasEmission: EmissionResult{
			EmissionFactor: 0,
			TotalEmission:  0,
		},
	}
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input CalculationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := calculateEmissions(input)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	http.HandleFunc("/calculate", handleCalculate)
	http.Handle("/", http.FileServer(http.Dir("static")))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
