package server

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type PageData struct {
	Result string
}

const (
	TypeLength      = "length"
	TypeWeight      = "weight"
	TypeTemperature = "temperature"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", lengthHandler)
	mux.HandleFunc("/convert/length", lengthHandler)
	mux.HandleFunc("/convert/weight", weightHandler)
	mux.HandleFunc("/convert/temperature", temperatureHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func lengthHandler(w http.ResponseWriter, r *http.Request) {
	handlerPages(w, r, TypeLength)
}

func weightHandler(w http.ResponseWriter, r *http.Request) {
	handlerPages(w, r, TypeWeight)
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	handlerPages(w, r, TypeTemperature)
}

func handlerPages(w http.ResponseWriter, r *http.Request, typeOfConvert string) {
	if r.Method == http.MethodPost {
		// Обработка POST-запроса
		typeString := typeOfConvert // Тип конвертации (вес)
		valueStr := r.FormValue("value")
		from := r.FormValue("from")
		to := r.FormValue("to")

		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, "Invalid value", http.StatusBadRequest)
			return
		}

		result := convert(value, from, to, typeString)
		data := PageData{Result: fmt.Sprintf("%.2f %s is %.2f %s", value, from, result, to)}

		renderTemplate(w, "frontend/"+typeOfConvert+".html", data)
	} else if r.Method == http.MethodGet {
		renderTemplate(w, "frontend/"+typeOfConvert+".html", PageData{Result: ""})
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	// Парсим HTML-шаблон
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Отображаем страницу с данными
	t.Execute(w, data)
}

func convert(value float64, from, to, typeString string) float64 {
	var convertedData float64
	switch typeString {
	case TypeLength:
		convertedData = convertLength(value, from, to)
	case TypeWeight:
		convertedData = convertWeight(value, from, to)
	case TypeTemperature:
		convertedData = convertTemperature(value, from, to)
	}

	return convertedData

}

func convertLength(value float64, from, to string) float64 {
	convertedToMilimeters := convertToMillimeter(value, from)
	convertToNeededValue := convertMillimeterToValue(convertedToMilimeters, to)
	return convertToNeededValue
}

func convertToMillimeter(value float64, from string) float64 {
	var mm float64
	switch from {
	case "millimeters":
		mm = value
	case "centimeters":
		mm = value * 10
	case "meters":
		mm = value * 1000
	case "kilometers":
		mm = value * 1000000
	case "inches":
		mm = value * 25.4
	case "foots":
		mm = value * 304.8
	case "yards":
		mm = value * 914.4
	case "miles":
		mm = value * 1609344
	}
	return mm
}

func convertMillimeterToValue(millimeters float64, to string) float64 {
	var finalValue float64
	switch to {
	case "millimeters":
		finalValue = millimeters
	case "centimeters":
		finalValue = millimeters / 10
	case "meters":
		finalValue = millimeters / 1000
	case "kilometers":
		finalValue = millimeters / 1000000
	case "inches":
		finalValue = millimeters / 25.4
	case "foots":
		finalValue = millimeters / 304.8
	case "yards":
		finalValue = millimeters / 914.4
	case "miles":
		finalValue = millimeters / 1609344
	}
	return finalValue
}

func convertWeight(value float64, from, to string) float64 {
	convertedToMilligrams := convertToMilligram(value, from)
	convertedToNeededValue := convertMilligramsToValue(convertedToMilligrams, to)
	return convertedToNeededValue
}

func convertToMilligram(value float64, from string) float64 {
	var mg float64
	switch from {
	case "milligrams":
		mg = value
	case "grams":
		mg = value * 1000
	case "kilograms":
		mg = value * 1000000
	case "ounces":
		mg = value * 28349.5231
	case "pounds":
		mg = value / 0.0000022046
	}
	return mg
}

func convertMilligramsToValue(milligrams float64, to string) float64 {
	var finalValue float64
	switch to {
	case "milligrams":
		finalValue = milligrams
	case "grams":
		finalValue = milligrams / 1000
	case "kilograms":
		finalValue = milligrams / 1000000
	case "ounces":
		finalValue = milligrams / 28349.5231
	case "pounds":
		finalValue = milligrams * 0.0000022046
	}
	return finalValue
}

func convertTemperature(value float64, from, to string) float64 {
	convertedToCelsius := convertToCelcius(value, from)
	convertedToNeededValue := convertCelsiusToValue(convertedToCelsius, to)
	return convertedToNeededValue
}

func convertToCelcius(value float64, from string) float64 {
	var c float64
	switch from {
	case "Celsius":
		c = value
	case "Fahrenheit":
		c = (value - 32) / 1.8000
	case "Kelvin":
		c = value - 273.15
	}
	return c
}

func convertCelsiusToValue(celsius float64, to string) float64 {
	var finalValue float64
	switch to {
	case "Celsius":
		finalValue = celsius
	case "Fahrenheit":
		finalValue = celsius*1.8000 + 32
	case "Kelvin":
		finalValue = celsius + 273.15
	}
	return finalValue
}
