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
	TypeTemperature = "temp"
)

func StartServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", parseHTMLHandler)
	mux.HandleFunc("/convert", convertHandler)

	// Запускаем сервер на порту 8080
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

func parseHTMLHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим HTML-шаблон
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	// Отображаем страницу с пустым результатом
	data := PageData{Result: ""}
	tmpl.Execute(w, data)
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из формы
	typeString := r.FormValue("type")
	valueStr := r.FormValue("value")
	from := r.FormValue("from")
	to := r.FormValue("to")

	// Преобразуем значение в число
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	// Выполняем конвертацию
	result := convert(value, from, to, typeString)
	data := PageData{Result: fmt.Sprintf("%.2f %s is %.2f %s", value, from, result, to)}

	// Парсим HTML-шаблон и отображаем результат
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
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
	case "millimeter":
		mm = value
	case "centimeters":
		mm = value * 10
	case "meters":
		mm = value * 1000
	case "kilometers":
		mm = value * 1000000
	case "inch":
		mm = value * 25.4
	case "foot":
		mm = value * 304.8
	case "yard":
		mm = value * 914.4
	case "mile":
		mm = value * 1609344
	}
	return mm
}

func convertMillimeterToValue(millimeters float64, to string) float64 {
	var finalValue float64
	switch to {
	case "millimeter":
		finalValue = millimeters
	case "centimeters":
		finalValue = millimeters / 10
	case "meters":
		finalValue = millimeters / 1000
	case "kilometers":
		finalValue = millimeters / 1000000
	case "inch":
		finalValue = millimeters / 25.4
	case "foot":
		finalValue = millimeters / 304.8
	case "yard":
		finalValue = millimeters / 914.4
	case "mile":
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
	case "milligram":
		mg = value
	case "gram":
		mg = value * 1000
	case "kilogram":
		mg = value * 1000000
	case "ounce":
		mg = value * 28349.5231
	case "pound":
		mg = value / 0.0000022046
	}
	return mg
}

func convertMilligramsToValue(milligrams float64, to string) float64 {
	var finalValue float64
	switch to {
	case "milligram":
		finalValue = milligrams
	case "gram":
		finalValue = milligrams / 1000
	case "kilogram":
		finalValue = milligrams / 1000000
	case "ounce":
		finalValue = milligrams / 28349.5231
	case "pound":
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
	case "celsius":
		c = value
	case "fahrenheit":
		c = (value - 32) / 1.8000
	case "kelvin":
		c = value - 273.15
	}
	return c
}

func convertCelsiusToValue(celsius float64, to string) float64 {
	var finalValue float64
	switch to {
	case "celsius":
		finalValue = celsius
	case "fahrenheit":
		finalValue = celsius*1.8000 + 32
	case "kelvin":
		finalValue = celsius + 273.15
	}
	return finalValue
}
