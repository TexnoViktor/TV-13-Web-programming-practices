package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Структури даних
type Task1Data struct {
	H, C, S, N, O, W, A float64
}

type Task2Data struct {
	CComb, HComb, OComb, SComb, W, A, V float64
	QComb                               float64
}

type Results struct {
	Task1 struct {
		Kpc, Kpg        float64
		DryComposition  Task1Data
		CombComposition Task1Data
		Qpn, Qdn, Qan   float64
	}
	Task2 struct {
		CWork, HWork, OWork, SWork, AWork, VWork float64
		QWork                                    float64
	}
}

var results Results

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/task1", task1Handler)
	http.HandleFunc("/task2", task2Handler)
	http.HandleFunc("/results", resultsHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Сервер запущено на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Обробники
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func task1Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/task1.html"))

	if r.Method == http.MethodPost {
		data := Task1Data{}

		data.H, _ = strconv.ParseFloat(r.FormValue("H"), 64)
		data.C, _ = strconv.ParseFloat(r.FormValue("C"), 64)
		data.S, _ = strconv.ParseFloat(r.FormValue("S"), 64)
		data.N, _ = strconv.ParseFloat(r.FormValue("N"), 64)
		data.O, _ = strconv.ParseFloat(r.FormValue("O"), 64)
		data.W, _ = strconv.ParseFloat(r.FormValue("W"), 64)
		data.A, _ = strconv.ParseFloat(r.FormValue("A"), 64)

		// Розрахунки для Task1
		Kpc := 100 / (100 - data.W)
		Kpg := 100 / (100 - data.W - data.A)

		results.Task1.Kpc = Kpc
		results.Task1.Kpg = Kpg

		// Склад сухої маси
		results.Task1.DryComposition = Task1Data{
			H: data.H * Kpc,
			C: data.C * Kpc,
			S: data.S * Kpc,
			N: data.N * Kpc,
			O: data.O * Kpc,
			A: data.A * Kpc,
		}

		// Склад горючої маси
		results.Task1.CombComposition = Task1Data{
			H: data.H * Kpg,
			C: data.C * Kpg,
			S: data.S * Kpg,
			N: data.N * Kpg,
			O: data.O * Kpg,
		}

		// Теплота згоряння
		Qpn := 339*data.C + 1030*data.H - 108.8*(data.O-data.S) - 25*data.W
		results.Task1.Qpn = Qpn / 1000
		results.Task1.Qdn = (Qpn + 0.025*data.W) * 100 / (100 - data.W) / 1000
		results.Task1.Qan = (Qpn + 0.025*data.W) * 100 / (100 - data.W - data.A) / 1000

		http.Redirect(w, r, "/results", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}

func task2Handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/task2.html"))

	if r.Method == http.MethodPost {
		data := Task2Data{}

		data.CComb, _ = strconv.ParseFloat(r.FormValue("CComb"), 64)
		data.HComb, _ = strconv.ParseFloat(r.FormValue("HComb"), 64)
		data.OComb, _ = strconv.ParseFloat(r.FormValue("OComb"), 64)
		data.SComb, _ = strconv.ParseFloat(r.FormValue("SComb"), 64)
		data.W, _ = strconv.ParseFloat(r.FormValue("W"), 64)
		data.A, _ = strconv.ParseFloat(r.FormValue("A"), 64)
		data.V, _ = strconv.ParseFloat(r.FormValue("V"), 64)
		data.QComb, _ = strconv.ParseFloat(r.FormValue("QComb"), 64)

		// Розрахунки для Task2
		results.Task2.CWork = data.CComb * (100 - data.W - data.A) / 100
		results.Task2.HWork = data.HComb * (100 - data.W - data.A) / 100
		results.Task2.OWork = data.OComb * (100 - data.W - data.A) / 100
		results.Task2.SWork = data.SComb * (100 - data.W - data.A) / 100
		results.Task2.AWork = data.A * (100 - data.W) / 100
		results.Task2.VWork = data.V * (100 - data.W) / 100
		results.Task2.QWork = data.QComb * (100 - data.W - data.A) / 100

		http.Redirect(w, r, "/results", http.StatusSeeOther)
		return
	}

	tmpl.Execute(w, nil)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/results.html"))
	tmpl.Execute(w, results)
}
