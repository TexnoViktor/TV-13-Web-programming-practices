package main

import (
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// EquipmentInput represents input parameters for a single piece of equipment
type EquipmentInput struct {
	Name              string  // Найменування ЕП
	Efficiency        float64 // ηн - номінальне значення коефіцієнта корисної дії ЕП
	PowerFactor       float64 // cos φ - коефіцієнт потужності навантаження
	Voltage           float64 // Uн - напруга навантаження, кВ
	Quantity          int     // n - кількість ЕП, шт
	Power             float64 // Рн - номінальна потужність ЕП, кВт
	UsageCoef         float64 // КВ - коефіцієнт використання
	ReactivePowerCoef float64 // tgφ - коефіцієнт реактивної потужності
}

// EquipmentOutput represents calculated parameters for a single piece of equipment
type EquipmentOutput struct {
	EquipmentInput
	PowerTotal     float64 // n·Pн, кВт
	PowerWithUsage float64 // n·Pн·КВ, кВт
	ReactivePower  float64 // n·Pн·КВ·tgφ, квар
	PowerSquared   float64 // n·Pн², кВт²
	CurrentRated   float64 // Ip - розрахунковий струм, А
}

// BusOutput represents calculated parameters for a distribution bus (ШР)
type BusOutput struct {
	Name                     string
	Equipment                []EquipmentOutput
	UsageCoefGroup           float64 // КВ - груповий коефіцієнт використання
	EffectiveQuantity        float64 // ne - ефективна кількість ЕП
	EffectiveQuantityRounded int     // ne rounded
	PowerCoef                float64 // Kр - розрахунковий коефіцієнт активної потужності
	ActivePower              float64 // Pp - розрахункове активне навантаження, кВт
	ReactivePower            float64 // Qp - розрахункове реактивне навантаження, квар
	ApparentPower            float64 // Sp - повна потужність, кВА
	BusCurrent               float64 // Ip - розрахунковий груповий струм, А
}

// WorkshopOutput represents calculated parameters for the entire workshop
type WorkshopOutput struct {
	Buses                         []BusOutput
	LargeEquipment                []EquipmentOutput
	UsageCoefTotal                float64 // КВ - загальний коефіцієнт використання
	EffectiveQuantityTotal        float64 // ne - загальна ефективна кількість ЕП
	EffectiveQuantityTotalRounded int     // ne rounded
	PowerCoefTotal                float64 // Kр - загальний розрахунковий коефіцієнт активної потужності
	ActivePowerTotal              float64 // Pp - загальне розрахункове активне навантаження, кВт
	ReactivePowerTotal            float64 // Qp - загальне розрахункове реактивне навантаження, квар
	ApparentPowerTotal            float64 // Sp - загальна повна потужність, кВА
	TotalCurrent                  float64 // Ip - загальний розрахунковий струм, А
}

// CalculateEquipmentOutput calculates all parameters for a piece of equipment
func CalculateEquipmentOutput(input EquipmentInput) EquipmentOutput {
	output := EquipmentOutput{EquipmentInput: input}

	// Calculate n·Pн
	output.PowerTotal = float64(input.Quantity) * input.Power

	// Calculate n·Pн·КВ
	output.PowerWithUsage = output.PowerTotal * input.UsageCoef

	// Calculate n·Pн·КВ·tgφ
	output.ReactivePower = output.PowerWithUsage * input.ReactivePowerCoef

	// Calculate n·Pн²
	output.PowerSquared = float64(input.Quantity) * math.Pow(input.Power, 2)

	// Calculate Ip - розрахунковий струм
	output.CurrentRated = (float64(input.Quantity) * input.Power) / (math.Sqrt(3) * input.Voltage * input.PowerFactor * input.Efficiency)

	return output
}

// LookupPowerCoef determines the power coefficient from the lookup tables
func LookupPowerCoef(usageCoef float64, effectiveQuantity int, isForTransformer bool) float64 {
	// This is a simplified implementation of the table lookup
	// In a real application, you would implement the complete lookup tables

	if isForTransformer {
		// Table 6.4 for transformers (T0 = 2.5 hour)
		if effectiveQuantity <= 5 {
			if usageCoef <= 0.2 {
				return 1.0
			} else if usageCoef <= 0.4 {
				return 0.9
			} else {
				return 0.8
			}
		} else if effectiveQuantity <= 50 {
			if usageCoef <= 0.2 {
				return 0.8
			} else {
				return 0.7
			}
		} else {
			return 0.65
		}
	} else {
		// Table 6.3 for normal lines (T0 = 10 min)
		if effectiveQuantity <= 10 {
			if usageCoef <= 0.2 {
				return 1.25
			} else if usageCoef <= 0.4 {
				return 1.2
			} else if usageCoef <= 0.6 {
				return 1.1
			} else {
				return 1.0
			}
		} else if effectiveQuantity <= 50 {
			if usageCoef <= 0.2 {
				return 1.1
			} else {
				return 1.0
			}
		} else {
			return 1.0
		}
	}
}

// CalculateBusOutput calculates all parameters for a distribution bus (ШР)
func CalculateBusOutput(name string, equipment []EquipmentOutput) BusOutput {
	bus := BusOutput{
		Name:      name,
		Equipment: equipment,
	}

	// Calculate total values
	var totalPowerNominal float64 = 0
	var totalPowerWithUsage float64 = 0
	var totalReactivePower float64 = 0
	var totalPowerSquared float64 = 0

	for _, eq := range equipment {
		totalPowerNominal += eq.PowerTotal
		totalPowerWithUsage += eq.PowerWithUsage
		totalReactivePower += eq.ReactivePower
		totalPowerSquared += eq.PowerSquared
	}

	// Calculate КВ - груповий коефіцієнт використання
	bus.UsageCoefGroup = totalPowerWithUsage / totalPowerNominal

	// Calculate ne - ефективна кількість ЕП
	bus.EffectiveQuantity = math.Pow(totalPowerNominal, 2) / totalPowerSquared
	bus.EffectiveQuantityRounded = int(bus.EffectiveQuantity)

	// Lookup Kр - розрахунковий коефіцієнт активної потужності
	bus.PowerCoef = LookupPowerCoef(bus.UsageCoefGroup, bus.EffectiveQuantityRounded, false)

	// Calculate Pp - розрахункове активне навантаження
	bus.ActivePower = bus.PowerCoef * totalPowerWithUsage

	// Calculate Qp - розрахункове реактивне навантаження
	// For ne <= 10, Qp = 1.1 * total reactive power, otherwise Qp = total reactive power
	if bus.EffectiveQuantityRounded <= 10 {
		bus.ReactivePower = 1.1 * totalReactivePower
	} else {
		bus.ReactivePower = totalReactivePower
	}

	// Calculate Sp - повна потужність
	bus.ApparentPower = math.Sqrt(math.Pow(bus.ActivePower, 2) + math.Pow(bus.ReactivePower, 2))

	// Calculate Ip - розрахунковий груповий струм (using voltage of the first equipment for simplicity)
	if len(equipment) > 0 {
		bus.BusCurrent = bus.ActivePower / equipment[0].Voltage
	}

	return bus
}

// CalculateWorkshopOutput calculates parameters for the entire workshop
func CalculateWorkshopOutput(buses []BusOutput, largeEquipment []EquipmentOutput) WorkshopOutput {
	workshop := WorkshopOutput{
		Buses:          buses,
		LargeEquipment: largeEquipment,
	}

	// Calculate totals across all buses and large equipment
	var totalPowerNominal float64 = 0
	var totalPowerWithUsage float64 = 0
	var totalPowerSquared float64 = 0

	for _, bus := range buses {
		for _, eq := range bus.Equipment {
			totalPowerNominal += eq.PowerTotal
			totalPowerWithUsage += eq.PowerWithUsage
			totalPowerSquared += eq.PowerSquared
		}
	}

	for _, eq := range largeEquipment {
		totalPowerNominal += eq.PowerTotal
		totalPowerWithUsage += eq.PowerWithUsage
		totalPowerSquared += eq.PowerSquared
	}

	// Calculate КВ - загальний коефіцієнт використання
	workshop.UsageCoefTotal = totalPowerWithUsage / totalPowerNominal

	// Calculate ne - загальна ефективна кількість ЕП
	workshop.EffectiveQuantityTotal = math.Pow(totalPowerNominal, 2) / totalPowerSquared
	workshop.EffectiveQuantityTotalRounded = int(workshop.EffectiveQuantityTotal)

	// Lookup Kр - загальний розрахунковий коефіцієнт активної потужності (for transformer)
	workshop.PowerCoefTotal = LookupPowerCoef(workshop.UsageCoefTotal, workshop.EffectiveQuantityTotalRounded, true)

	// Calculate reactive power for all large equipment
	var totalLargeReactivePower float64 = 0
	for _, eq := range largeEquipment {
		totalLargeReactivePower += eq.ReactivePower
	}

	// Calculate total active and reactive power from buses
	var totalBusActivePower float64 = 0
	var totalBusReactivePower float64 = 0
	for _, bus := range buses {
		totalBusActivePower += bus.ActivePower
		totalBusReactivePower += bus.ReactivePower
	}

	// Calculate Pp - загальне розрахункове активне навантаження
	workshop.ActivePowerTotal = workshop.PowerCoefTotal * (totalPowerWithUsage)

	// Calculate Qp - загальне розрахункове реактивне навантаження
	totalReactivePowerWithUsage := 0.0
	for _, bus := range buses {
		for _, eq := range bus.Equipment {
			totalReactivePowerWithUsage += eq.ReactivePower
		}
	}
	for _, eq := range largeEquipment {
		totalReactivePowerWithUsage += eq.ReactivePower
	}
	workshop.ReactivePowerTotal = workshop.PowerCoefTotal * totalReactivePowerWithUsage

	// Calculate Sp - загальна повна потужність
	workshop.ApparentPowerTotal = math.Sqrt(math.Pow(workshop.ActivePowerTotal, 2) + math.Pow(workshop.ReactivePowerTotal, 2))

	// Calculate Ip - загальний розрахунковий струм (using voltage of the first equipment for simplicity)
	if len(buses) > 0 && len(buses[0].Equipment) > 0 {
		workshop.TotalCurrent = workshop.ActivePowerTotal / buses[0].Equipment[0].Voltage
	}

	return workshop
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Cannot parse form", http.StatusBadRequest)
			return
		}

		// Example: Process data for the first bus (ШР1)
		sr1Equipment := []EquipmentInput{
			{
				Name:              "Шліфувальний верстат (1-4)",
				Efficiency:        parseFloat(r.FormValue("efficiency_1"), 0.92),
				PowerFactor:       parseFloat(r.FormValue("power_factor_1"), 0.9),
				Voltage:           parseFloat(r.FormValue("voltage_1"), 0.38),
				Quantity:          parseInt(r.FormValue("quantity_1"), 4),
				Power:             parseFloat(r.FormValue("power_1"), 20),
				UsageCoef:         parseFloat(r.FormValue("usage_coef_1"), 0.15),
				ReactivePowerCoef: parseFloat(r.FormValue("reactive_power_coef_1"), 1.33),
			},
			{
				Name:              "Свердлильний верстат (5-6)",
				Efficiency:        parseFloat(r.FormValue("efficiency_2"), 0.92),
				PowerFactor:       parseFloat(r.FormValue("power_factor_2"), 0.9),
				Voltage:           parseFloat(r.FormValue("voltage_2"), 0.38),
				Quantity:          parseInt(r.FormValue("quantity_2"), 2),
				Power:             parseFloat(r.FormValue("power_2"), 14),
				UsageCoef:         parseFloat(r.FormValue("usage_coef_2"), 0.12),
				ReactivePowerCoef: parseFloat(r.FormValue("reactive_power_coef_2"), 1.0),
			},
			// Add all other equipment for ШР1...
			{
				Name:              "Фугувальний верстат (9-12)",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          4,
				Power:             42,
				UsageCoef:         0.15,
				ReactivePowerCoef: 1.33,
			},
			{
				Name:              "Циркулярна пила (13)",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          1,
				Power:             36,
				UsageCoef:         0.3,
				ReactivePowerCoef: 1.52,
			},
			{
				Name:              "Прес (16)",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          1,
				Power:             20,
				UsageCoef:         0.5,
				ReactivePowerCoef: 0.75,
			},
			{
				Name:              "Полірувальний верстат (24)",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          1,
				Power:             40,
				UsageCoef:         0.2,
				ReactivePowerCoef: 1.0,
			},
			{
				Name:              "Фрезерний верстат (26-27)",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          2,
				Power:             32,
				UsageCoef:         0.2,
				ReactivePowerCoef: 1.0,
			},
			{
				Name:              "Вентилятор (36)",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          1,
				Power:             20,
				UsageCoef:         0.65,
				ReactivePowerCoef: 0.75,
			},
		}

		// Process data for large equipment connected directly to the transformer
		largeEquipment := []EquipmentInput{
			{
				Name:              "Зварювальний трансформатор",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          2,
				Power:             100,
				UsageCoef:         0.2,
				ReactivePowerCoef: 3.0,
			},
			{
				Name:              "Сушильна шафа",
				Efficiency:        0.92,
				PowerFactor:       0.9,
				Voltage:           0.38,
				Quantity:          2,
				Power:             120,
				UsageCoef:         0.8,
				ReactivePowerCoef: 0.0, // No reactive power
			},
		}

		// Process equipment for ШР1
		sr1EquipmentOutput := make([]EquipmentOutput, len(sr1Equipment))
		for i, eq := range sr1Equipment {
			sr1EquipmentOutput[i] = CalculateEquipmentOutput(eq)
		}

		// Process large equipment
		largeEquipmentOutput := make([]EquipmentOutput, len(largeEquipment))
		for i, eq := range largeEquipment {
			largeEquipmentOutput[i] = CalculateEquipmentOutput(eq)
		}

		// Calculate ШР1 parameters
		sr1Output := CalculateBusOutput("ШР1", sr1EquipmentOutput)

		// For simplicity, assume ШР2 and ШР3 are identical to ШР1
		sr2Output := sr1Output
		sr2Output.Name = "ШР2"
		sr3Output := sr1Output
		sr3Output.Name = "ШР3"

		buses := []BusOutput{sr1Output, sr2Output, sr3Output}

		// Calculate workshop parameters
		workshopOutput := CalculateWorkshopOutput(buses, largeEquipmentOutput)

		// Render results
		tmpl, err := template.ParseFiles("templates/results.html")
		if err != nil {
			http.Error(w, "Template error", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, workshopOutput)
		if err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
			return
		}

		return
	}

	// If not POST, show the form
	tmpl, err := template.ParseFiles("templates/form.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Template execution error", http.StatusInternalServerError)
		return
	}
}

// Helper function to parse float values with default
func parseFloat(value string, defaultValue float64) float64 {
	if value == "" {
		return defaultValue
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}

	return f
}

// Helper function to parse int values with default
func parseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return i
}

func main() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handle the calculator page
	http.HandleFunc("/", handleCalculate)

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
