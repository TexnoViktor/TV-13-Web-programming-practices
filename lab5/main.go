package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageData struct {
	Title                string
	ErrorMessage         string
	SingleCircuitResults *SingleCircuitResults
	DoubleCircuitResults *DoubleCircuitResults
	LossesResults        *LossesResults
}

type SingleCircuitResults struct {
	FrequencyOfFailures          float64
	AvgRecoveryTime              float64
	EmergencyDowntimeCoefficient float64
	PlannedDowntimeCoefficient   float64
}

type DoubleCircuitResults struct {
	SimultaneousFailureFreq        float64
	TotalFrequencyOfFailures       float64
	AvgRecoveryTimeDC              float64
	EmergencyDowntimeCoefficientDC float64
}

type LossesResults struct {
	EmergencyNonSupply float64
	PlannedNonSupply   float64
	TotalLosses        float64
}

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Define routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/calculate-reliability", calculateReliabilityHandler)
	http.HandleFunc("/calculate-losses", calculateLossesHandler)

	// Start server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Веб-калькулятор надійності систем електропередачі",
	}
	renderTemplate(w, "index.html", data)
}

func calculateReliabilityHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	// Parse single circuit system elements failure rates
	w_B110kV, err1 := strconv.ParseFloat(r.FormValue("w_B110kV"), 64)
	w_PL110kV, err2 := strconv.ParseFloat(r.FormValue("w_PL110kV"), 64)
	w_T110_10kV, err3 := strconv.ParseFloat(r.FormValue("w_T110_10kV"), 64)
	w_SW10kV, err4 := strconv.ParseFloat(r.FormValue("w_SW10kV"), 64)
	w_Conn10kV, err5 := strconv.ParseFloat(r.FormValue("w_Conn10kV"), 64)
	numConnections, err6 := strconv.ParseFloat(r.FormValue("numConnections"), 64)

	// Parse repair times
	t_B110kV, err7 := strconv.ParseFloat(r.FormValue("t_B110kV"), 64)
	t_PL110kV, err8 := strconv.ParseFloat(r.FormValue("t_PL110kV"), 64)
	t_T110_10kV, err9 := strconv.ParseFloat(r.FormValue("t_T110_10kV"), 64)
	t_SW10kV, err10 := strconv.ParseFloat(r.FormValue("t_SW10kV"), 64)
	t_Conn10kV, err11 := strconv.ParseFloat(r.FormValue("t_Conn10kV"), 64)

	// Parse sectional switch parameters
	w_SectSW, err12 := strconv.ParseFloat(r.FormValue("w_SectSW"), 64)

	data := PageData{
		Title: "Результати розрахунків надійності систем електропередачі",
	}

	// Check for errors in parsing
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil ||
		err7 != nil || err8 != nil || err9 != nil || err10 != nil || err11 != nil || err12 != nil {
		data.ErrorMessage = "Помилка при введенні даних. Будь ласка, перевірте правильність введених значень."
		renderTemplate(w, "index.html", data)
		return
	}

	// Calculate Single Circuit System Reliability
	singleResults := calculateSingleCircuitReliability(
		w_B110kV, w_PL110kV, w_T110_10kV, w_SW10kV, w_Conn10kV, numConnections,
		t_B110kV, t_PL110kV, t_T110_10kV, t_SW10kV, t_Conn10kV,
	)

	// Calculate Double Circuit System Reliability
	doubleResults := calculateDoubleCircuitReliability(
		singleResults.FrequencyOfFailures, singleResults.EmergencyDowntimeCoefficient,
		singleResults.PlannedDowntimeCoefficient, w_SectSW,
	)

	data.SingleCircuitResults = singleResults
	data.DoubleCircuitResults = doubleResults

	renderTemplate(w, "reliability_results.html", data)
}

func calculateLossesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	r.ParseForm()

	// Parse loss parameters
	specificLossesEmergency, err1 := strconv.ParseFloat(r.FormValue("specificLossesEmergency"), 64)
	specificLossesPlanned, err2 := strconv.ParseFloat(r.FormValue("specificLossesPlanned"), 64)
	transformerFailureRate, err3 := strconv.ParseFloat(r.FormValue("transformerFailureRate"), 64)
	avgRecoveryTime, err4 := strconv.ParseFloat(r.FormValue("avgRecoveryTime"), 64)
	plannedDowntimeCoef, err5 := strconv.ParseFloat(r.FormValue("plannedDowntimeCoef"), 64)
	maxPower, err6 := strconv.ParseFloat(r.FormValue("maxPower"), 64)
	utilHours, err7 := strconv.ParseFloat(r.FormValue("utilHours"), 64)

	data := PageData{
		Title: "Результати розрахунків збитків від перерв електропостачання",
	}

	// Check for errors in parsing
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || err6 != nil || err7 != nil {
		data.ErrorMessage = "Помилка при введенні даних. Будь ласка, перевірте правильність введених значень."
		renderTemplate(w, "index.html", data)
		return
	}

	// Calculate losses
	lossesResults := calculateLosses(
		specificLossesEmergency, specificLossesPlanned,
		transformerFailureRate, avgRecoveryTime, plannedDowntimeCoef,
		maxPower, utilHours,
	)

	data.LossesResults = lossesResults

	renderTemplate(w, "losses_results.html", data)
}

func calculateSingleCircuitReliability(
	w_B110kV, w_PL110kV, w_T110_10kV, w_SW10kV, w_Conn10kV, numConnections,
	t_B110kV, t_PL110kV, t_T110_10kV, t_SW10kV, t_Conn10kV float64) *SingleCircuitResults {

	// Total failure frequency of single circuit system
	w_SC := w_B110kV + w_PL110kV + w_T110_10kV + w_SW10kV + w_Conn10kV*numConnections

	// Average recovery time calculation
	sumWT := w_B110kV*t_B110kV + w_PL110kV*t_PL110kV + w_T110_10kV*t_T110_10kV +
		w_SW10kV*t_SW10kV + (w_Conn10kV*numConnections)*t_Conn10kV
	t_SC := sumWT / w_SC

	// Emergency downtime coefficient
	k_a_SC := (w_SC * t_SC) / 8760

	// Planned downtime coefficient (using transformer as max)
	k_p_max := 1.0 * 43.0 / 8760.0 // Assuming μ=1 and t_p=43 for transformer from table
	k_p_SC := 1.2 * k_p_max

	return &SingleCircuitResults{
		FrequencyOfFailures:          w_SC,
		AvgRecoveryTime:              t_SC,
		EmergencyDowntimeCoefficient: k_a_SC,
		PlannedDowntimeCoefficient:   k_p_SC,
	}
}

func calculateDoubleCircuitReliability(w_SC, k_a_SC, k_p_SC, w_SectSW float64) *DoubleCircuitResults {
	// Frequency of simultaneous failures of both circuits
	w_2C := 2 * w_SC * (k_a_SC + 0.5*k_p_SC)

	// Total frequency of failures including the sectional switch
	w_DC := w_2C + w_SectSW

	// Average recovery time and emergency downtime for double circuit
	// These would need more complete formulas based on the provided examples
	// Here is a simplified approach
	t_DC := 2.0 // Simplified assumption
	k_a_DC := (w_DC * t_DC) / 8760

	return &DoubleCircuitResults{
		SimultaneousFailureFreq:        w_2C,
		TotalFrequencyOfFailures:       w_DC,
		AvgRecoveryTimeDC:              t_DC,
		EmergencyDowntimeCoefficientDC: k_a_DC,
	}
}

func calculateLosses(
	specificLossesEmergency, specificLossesPlanned,
	transformerFailureRate, avgRecoveryTime, plannedDowntimeCoef,
	maxPower, utilHours float64) *LossesResults {

	// Calculate expected emergency non-supply of electricity
	emergencyNonSupply := transformerFailureRate * maxPower * utilHours

	// Calculate expected planned non-supply of electricity
	plannedNonSupply := plannedDowntimeCoef * maxPower * utilHours

	// Calculate total losses
	totalLosses := specificLossesEmergency*emergencyNonSupply + specificLossesPlanned*plannedNonSupply

	return &LossesResults{
		EmergencyNonSupply: emergencyNonSupply,
		PlannedNonSupply:   plannedNonSupply,
		TotalLosses:        totalLosses,
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
