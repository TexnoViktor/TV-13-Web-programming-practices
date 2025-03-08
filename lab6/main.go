package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

// ElectricalDevice представляє електроприймач
type ElectricalDevice struct {
	Name           string  `json:"name"`
	Efficiency     float64 `json:"efficiency"`     // ηн - коефіцієнт корисної дії
	PowerFactor    float64 `json:"powerFactor"`    // cos φ - коефіцієнт потужності
	Voltage        float64 `json:"voltage"`        // Uн - напруга, кВ
	Quantity       int     `json:"quantity"`       // n - кількість, шт
	Power          float64 `json:"power"`          // Pн - номінальна потужність, кВт
	UsageFactor    float64 `json:"usageFactor"`    // КВ - коефіцієнт використання
	ReactiveFactor float64 `json:"reactiveFactor"` // tgφ - коефіцієнт реактивної потужності
}

// DeviceCalculation містить результати розрахунків для одного електроприймача
type DeviceCalculation struct {
	Device              ElectricalDevice // Вхідні дані пристрою
	NPower              float64          // n * Pн, кВт
	Current             float64          // Ip - розрахунковий струм, A
	NPowerUsage         float64          // n * Pн * КВ, кВт
	NPowerUsageReactive float64          // n * Pн * КВ * tgφ, квар
	NPowerSquared       float64          // n * Pн^2
}

// GroupCalculation містить результати групових розрахунків
type GroupCalculation struct {
	Name                 string              // Назва групи (ШР1, ШР2 тощо)
	Devices              []DeviceCalculation // Пристрої групи
	TotalQuantity        int                 // Загальна кількість пристроїв
	TotalPower           float64             // Сумарна номінальна потужність
	UsageFactor          float64             // КВ - груповий коефіцієнт використання
	EffectiveDeviceCount float64             // ne - ефективна кількість пристроїв
	PowerFactor          float64             // Кр - розрахунковий коефіцієнт активної потужності
	ActivePower          float64             // Pp - розрахункове активне навантаження, кВт
	ReactivePower        float64             // Qp - розрахункове реактивне навантаження, квар
	TotalPowerApparent   float64             // Sp - повна потужність, кВА
	GroupCurrent         float64             // Ip - розрахунковий груповий струм, А
}

// TotalCalculation містить підсумкові результати розрахунків
type TotalCalculation struct {
	Groups               []GroupCalculation // Групи розрахунків (ШР1, ШР2, ШР3)
	TotalQuantity        int                // Загальна кількість пристроїв
	TotalPower           float64            // Сумарна номінальна потужність
	UsageFactor          float64            // КВ - коефіцієнт використання цеху в цілому
	EffectiveDeviceCount float64            // ne - ефективна кількість пристроїв цеху
	PowerFactor          float64            // Кр - розрахунковий коефіцієнт активної потужності
	ActivePower          float64            // Pp - розрахункове активне навантаження, кВт
	ReactivePower        float64            // Qp - розрахункове реактивне навантаження, квар
	TotalPowerApparent   float64            // Sp - повна потужність, кВА
	TotalCurrent         float64            // Ip - розрахунковий груповий струм, А
}

// CalculatePowerFactor визначає коефіцієнт Kр за таблицями
func CalculatePowerFactor(usageFactor, effectiveDeviceCount float64, isHighLevel bool) float64 {
	// Округлюємо ne до найближчого меншого цілого числа
	ne := math.Floor(effectiveDeviceCount)

	if isHighLevel {
		// Використовуємо таблицю 6.4 для високого рівня (T0 = 2,5 год.)
		if ne >= 50 {
			if usageFactor >= 0.7 {
				return 0.8
			} else if usageFactor >= 0.6 {
				return 0.8
			} else if usageFactor >= 0.5 {
				return 0.75
			} else if usageFactor >= 0.4 {
				return 0.7
			} else if usageFactor >= 0.3 {
				return 0.7
			} else if usageFactor >= 0.2 {
				return 0.65
			} else {
				return 0.65
			}
		} else if ne >= 25 && ne < 50 {
			if usageFactor >= 0.7 {
				return 0.85
			} else if usageFactor >= 0.6 {
				return 0.85
			} else if usageFactor >= 0.5 {
				return 0.8
			} else if usageFactor >= 0.4 {
				return 0.75
			} else if usageFactor >= 0.3 {
				return 0.75
			} else if usageFactor >= 0.2 {
				return 0.75
			} else {
				return 0.75
			}
		} else if ne >= 10 && ne < 25 {
			if usageFactor >= 0.7 {
				return 0.9
			} else if usageFactor >= 0.6 {
				return 0.9
			} else if usageFactor >= 0.5 {
				return 0.85
			} else if usageFactor >= 0.4 {
				return 0.85
			} else if usageFactor >= 0.3 {
				return 0.85
			} else if usageFactor >= 0.2 {
				return 0.8
			} else {
				return 0.8
			}
		} else if ne >= 6 && ne < 10 {
			if usageFactor >= 0.7 {
				return 0.9
			} else if usageFactor >= 0.6 {
				return 0.92
			} else if usageFactor >= 0.5 {
				return 0.93
			} else if usageFactor >= 0.4 {
				return 0.94
			} else if usageFactor >= 0.3 {
				return 0.95
			} else if usageFactor >= 0.2 {
				return 0.96
			} else {
				return 0.96
			}
		} else if ne == 5 {
			if usageFactor >= 0.7 {
				return 0.93
			} else if usageFactor >= 0.6 {
				return 0.94
			} else if usageFactor >= 0.5 {
				return 0.96
			} else if usageFactor >= 0.4 {
				return 0.98
			} else if usageFactor >= 0.3 {
				return 1.0
			} else if usageFactor >= 0.2 {
				return 1.02
			} else {
				return 1.05
			}
		} else if ne == 4 {
			if usageFactor >= 0.7 {
				return 0.97
			} else if usageFactor >= 0.6 {
				return 1.0
			} else if usageFactor >= 0.5 {
				return 1.04
			} else if usageFactor >= 0.4 {
				return 1.06
			} else if usageFactor >= 0.3 {
				return 1.19
			} else if usageFactor >= 0.2 {
				return 1.46
			} else {
				return 1.73
			}
		} else if ne == 3 {
			if usageFactor >= 0.7 {
				return 1.0
			} else if usageFactor >= 0.6 {
				return 1.08
			} else if usageFactor >= 0.5 {
				return 1.14
			} else if usageFactor >= 0.4 {
				return 1.23
			} else if usageFactor >= 0.3 {
				return 1.42
			} else if usageFactor >= 0.2 {
				return 1.8
			} else {
				return 2.17
			}
		} else if ne == 2 {
			if usageFactor >= 0.7 {
				return 1.0
			} else if usageFactor >= 0.6 {
				return 1.11
			} else if usageFactor >= 0.5 {
				return 1.24
			} else if usageFactor >= 0.4 {
				return 1.52
			} else if usageFactor >= 0.3 {
				return 1.9
			} else if usageFactor >= 0.2 {
				return 2.69
			} else {
				return 3.44
			}
		} else if ne == 1 {
			if usageFactor >= 0.7 {
				return 1.14
			} else if usageFactor >= 0.6 {
				return 1.33
			} else if usageFactor >= 0.5 {
				return 1.6
			} else if usageFactor >= 0.4 {
				return 2.0
			} else if usageFactor >= 0.3 {
				return 2.67
			} else if usageFactor >= 0.2 {
				return 4.0
			} else {
				return 5.33
			}
		}
	} else {
		// Використовуємо таблицю 6.3 для низького рівня (T0 = 10 хв.)
		if ne >= 100 {
			return 1.0
		} else if ne >= 80 {
			if usageFactor < 0.2 {
				return 1.0
			} else {
				return 1.16
			}
		} else if ne >= 60 {
			if usageFactor < 0.2 {
				return 1.0
			} else if usageFactor < 0.4 {
				return 1.03
			} else {
				return 1.25
			}
		} else if ne >= 50 {
			if usageFactor < 0.2 {
				return 1.0
			} else if usageFactor < 0.4 {
				return 1.07
			} else {
				return 1.3
			}
		} else if ne >= 40 {
			if usageFactor < 0.2 {
				return 1.0
			} else if usageFactor < 0.4 {
				return 1.13
			} else {
				return 1.4
			}
		} else if ne >= 35 {
			if usageFactor < 0.2 {
				return 1.0
			} else if usageFactor < 0.4 {
				return 1.16
			} else {
				return 1.44
			}
		} else if ne >= 30 {
			if usageFactor < 0.2 {
				return 1.05
			} else if usageFactor < 0.4 {
				return 1.21
			} else {
				return 1.51
			}
		} else if ne >= 25 {
			if usageFactor < 0.2 {
				return 1.1
			} else if usageFactor < 0.4 {
				return 1.27
			} else {
				return 1.6
			}
		} else if ne >= 20 {
			if usageFactor < 0.2 {
				return 1.16
			} else if usageFactor < 0.4 {
				return 1.35
			} else {
				return 1.72
			}
		} else if ne >= 18 {
			if usageFactor < 0.2 {
				return 1.19
			} else if usageFactor < 0.4 {
				return 1.39
			} else {
				return 1.78
			}
		} else if ne >= 16 {
			if usageFactor < 0.2 {
				return 1.23
			} else if usageFactor < 0.4 {
				return 1.43
			} else {
				return 1.85
			}
		} else if ne >= 14 {
			if usageFactor < 0.2 {
				return 1.27
			} else if usageFactor < 0.4 {
				return 1.49
			} else {
				return 1.94
			}
		} else if ne >= 12 {
			if usageFactor < 0.2 {
				return 1.32
			} else if usageFactor < 0.4 {
				return 1.56
			} else {
				return 2.04
			}
		} else if ne >= 10 {
			if usageFactor < 0.2 {
				return 1.39
			} else if usageFactor < 0.4 {
				return 1.65
			} else {
				return 2.18
			}
		} else if ne >= 9 {
			if usageFactor < 0.2 {
				return 1.43
			} else if usageFactor < 0.4 {
				return 1.71
			} else {
				return 2.27
			}
		} else if ne >= 8 {
			if usageFactor < 0.2 {
				return 1.48
			} else if usageFactor < 0.4 {
				return 1.78
			} else {
				return 2.37
			}
		} else if ne >= 7 {
			if usageFactor < 0.2 {
				return 1.54
			} else if usageFactor < 0.4 {
				return 1.86
			} else {
				return 2.49
			}
		} else if ne >= 6 {
			if usageFactor < 0.2 {
				return 1.62
			} else if usageFactor < 0.4 {
				return 1.96
			} else {
				return 2.64
			}
		} else if ne >= 5 {
			if usageFactor < 0.2 {
				return 1.72
			} else if usageFactor < 0.4 {
				return 2.09
			} else {
				return 2.84
			}
		} else if ne >= 4 {
			if usageFactor < 0.2 {
				return 1.91
			} else if usageFactor < 0.4 {
				return 2.35
			} else {
				return 3.24
			}
		} else if ne >= 3 {
			if usageFactor < 0.2 {
				return 2.31
			} else if usageFactor < 0.4 {
				return 2.89
			} else {
				return 4.06
			}
		} else if ne >= 2 {
			if usageFactor < 0.2 {
				return 3.39
			} else if usageFactor < 0.4 {
				return 4.33
			} else {
				return 6.22
			}
		} else if ne >= 1 {
			if usageFactor < 0.2 {
				return 4.0
			} else if usageFactor < 0.4 {
				return 5.33
			} else {
				return 8.0
			}
		}
	}

	// Якщо жодна умова не виконується, повертаємо значення за замовчуванням
	return 1.0
}

// CalculateDeviceCurrent розраховує струм для одного пристрою
func CalculateDeviceCurrent(device ElectricalDevice) float64 {
	return float64(device.Quantity) * device.Power / (math.Sqrt(3) * device.Voltage * device.PowerFactor * device.Efficiency)
}

// CalculateGroupData обчислює групові дані для електроприймачів
func CalculateGroupData(devices []ElectricalDevice, groupName string) GroupCalculation {
	var group GroupCalculation
	group.Name = groupName
	group.Devices = make([]DeviceCalculation, len(devices))

	totalQuantity := 0
	totalPower := 0.0
	totalPowerUsage := 0.0
	totalPowerUsageReactive := 0.0
	totalPowerSquared := 0.0

	// Розрахунок даних для кожного пристрою в групі
	for i, device := range devices {
		calc := DeviceCalculation{
			Device:              device,
			NPower:              float64(device.Quantity) * device.Power,
			Current:             CalculateDeviceCurrent(device),
			NPowerUsage:         float64(device.Quantity) * device.Power * device.UsageFactor,
			NPowerUsageReactive: float64(device.Quantity) * device.Power * device.UsageFactor * device.ReactiveFactor,
			NPowerSquared:       float64(device.Quantity) * device.Power * device.Power,
		}

		totalQuantity += device.Quantity
		totalPower += calc.NPower
		totalPowerUsage += calc.NPowerUsage
		totalPowerUsageReactive += calc.NPowerUsageReactive
		totalPowerSquared += calc.NPowerSquared

		group.Devices[i] = calc
	}

	group.TotalQuantity = totalQuantity
	group.TotalPower = totalPower

	// Розрахунок групового коефіцієнта використання
	if totalPower > 0 {
		group.UsageFactor = totalPowerUsage / totalPower
	}

	// Розрахунок ефективної кількості пристроїв
	if totalPowerSquared > 0 {
		group.EffectiveDeviceCount = math.Pow(totalPower, 2) / totalPowerSquared
	}

	// Визначення розрахункового коефіцієнта активної потужності
	group.PowerFactor = CalculatePowerFactor(group.UsageFactor, group.EffectiveDeviceCount, false)

	// Розрахунок розрахункового активного навантаження
	if group.EffectiveDeviceCount <= 10 {
		group.ActivePower = group.PowerFactor * totalPowerUsage
	} else {
		group.ActivePower = totalPowerUsage
	}

	// Розрахунок розрахункового реактивного навантаження
	if group.EffectiveDeviceCount <= 10 {
		group.ReactivePower = 1.1 * totalPowerUsageReactive
	} else {
		group.ReactivePower = totalPowerUsageReactive
	}

	// Розрахунок повної потужності
	group.TotalPowerApparent = math.Sqrt(math.Pow(group.ActivePower, 2) + math.Pow(group.ReactivePower, 2))

	// Розрахунок групового струму
	if devices[0].Voltage > 0 {
		group.GroupCurrent = group.ActivePower / devices[0].Voltage
	}

	return group
}

// CalculateTotalData обчислює загальні дані для всіх груп
func CalculateTotalData(groups []GroupCalculation) TotalCalculation {
	var total TotalCalculation
	total.Groups = groups

	totalQuantity := 0
	totalPower := 0.0
	totalPowerUsage := 0.0
	totalPowerUsageReactive := 0.0
	totalPowerSquared := 0.0

	for _, group := range groups {
		totalQuantity += group.TotalQuantity
		totalPower += group.TotalPower

		for _, device := range group.Devices {
			totalPowerUsage += device.NPowerUsage
			totalPowerUsageReactive += device.NPowerUsageReactive
			totalPowerSquared += device.NPowerSquared
		}
	}

	total.TotalQuantity = totalQuantity
	total.TotalPower = totalPower

	// Розрахунок коефіцієнта використання цеху в цілому
	if totalPower > 0 {
		total.UsageFactor = totalPowerUsage / totalPower
	}

	// Розрахунок ефективної кількості пристроїв цеху
	if totalPowerSquared > 0 {
		total.EffectiveDeviceCount = math.Pow(totalPower, 2) / totalPowerSquared
	}

	// Визначення розрахункового коефіцієнта активної потужності цеху
	total.PowerFactor = CalculatePowerFactor(total.UsageFactor, total.EffectiveDeviceCount, true)

	// Розрахунок розрахункового активного навантаження цеху
	total.ActivePower = total.PowerFactor * totalPowerUsage

	// Розрахунок розрахункового реактивного навантаження цеху
	total.ReactivePower = total.PowerFactor * totalPowerUsageReactive

	// Розрахунок повної потужності цеху
	total.TotalPowerApparent = math.Sqrt(math.Pow(total.ActivePower, 2) + math.Pow(total.ReactivePower, 2))

	// Розрахунок загального струму цеху
	if len(groups) > 0 && len(groups[0].Devices) > 0 {
		voltage := groups[0].Devices[0].Device.Voltage
		if voltage > 0 {
			total.TotalCurrent = total.ActivePower / voltage
		}
	}

	return total
}

// IndexHandler показує головну сторінку
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// CalculateHandler обробляє POST-запит з даними для розрахунку
func CalculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсинг форми
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отримання кількості пристроїв
	deviceCountStr := r.FormValue("deviceCount")
	deviceCount, err := strconv.Atoi(deviceCountStr)
	if err != nil {
		deviceCount = 0
	}

	// Створення слайсу пристроїв для трьох груп (ШР1, ШР2, ШР3)
	devicesGroup1 := make([]ElectricalDevice, deviceCount)

	// Заповнення даних пристроїв для ШР1
	for i := 0; i < deviceCount; i++ {
		efficiency, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("efficiency_%d", i)), 64)
		powerFactor, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("powerFactor_%d", i)), 64)
		voltage, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("voltage_%d", i)), 64)
		quantity, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("quantity_%d", i)))
		power, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("power_%d", i)), 64)
		usageFactor, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("usageFactor_%d", i)), 64)
		reactiveFactor, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("reactiveFactor_%d", i)), 64)

		devicesGroup1[i] = ElectricalDevice{
			Name:           r.FormValue(fmt.Sprintf("name_%d", i)),
			Efficiency:     efficiency,
			PowerFactor:    powerFactor,
			Voltage:        voltage,
			Quantity:       quantity,
			Power:          power,
			UsageFactor:    usageFactor,
			ReactiveFactor: reactiveFactor,
		}
	}

	// Для спрощення приймаємо, що ШР2 і ШР3 ідентичні ШР1
	devicesGroup2 := make([]ElectricalDevice, len(devicesGroup1))
	devicesGroup3 := make([]ElectricalDevice, len(devicesGroup1))
	copy(devicesGroup2, devicesGroup1)
	copy(devicesGroup3, devicesGroup1)

	// Отримання даних для крупних ЕП
	largeDeviceCountStr := r.FormValue("largeDeviceCount")
	largeDeviceCount, err := strconv.Atoi(largeDeviceCountStr)
	if err != nil {
		largeDeviceCount = 0
	}

	largeDevices := make([]ElectricalDevice, largeDeviceCount)

	// Заповнення даних крупних ЕП
	for i := 0; i < largeDeviceCount; i++ {
		efficiency, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("large_efficiency_%d", i)), 64)
		powerFactor, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("large_powerFactor_%d", i)), 64)
		voltage, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("large_voltage_%d", i)), 64)
		quantity, _ := strconv.Atoi(r.FormValue(fmt.Sprintf("large_quantity_%d", i)))
		power, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("large_power_%d", i)), 64)
		usageFactor, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("large_usageFactor_%d", i)), 64)
		reactiveFactor, _ := strconv.ParseFloat(r.FormValue(fmt.Sprintf("large_reactiveFactor_%d", i)), 64)

		largeDevices[i] = ElectricalDevice{
			Name:           r.FormValue(fmt.Sprintf("large_name_%d", i)),
			Efficiency:     efficiency,
			PowerFactor:    powerFactor,
			Voltage:        voltage,
			Quantity:       quantity,
			Power:          power,
			UsageFactor:    usageFactor,
			ReactiveFactor: reactiveFactor,
		}
	}

	// Розрахунок даних для груп
	group1 := CalculateGroupData(devicesGroup1, "ШР1")
	group2 := CalculateGroupData(devicesGroup2, "ШР2")
	group3 := CalculateGroupData(devicesGroup3, "ШР3")
	largeGroup := CalculateGroupData(largeDevices, "Крупні ЕП")

	// Об'єднання груп для загального розрахунку
	groups := []GroupCalculation{group1, group2, group3, largeGroup}
	total := CalculateTotalData(groups)

	// Відправка результатів у форматі JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(total)
}

func main() {
	// Вказуємо де шукати статичні файли
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Обробники запитів
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/calculate", CalculateHandler)

	// Запуск сервера
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
