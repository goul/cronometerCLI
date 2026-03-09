package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jrmycanady/gocronometer"
)

func main() {
	username := flag.String("username", "", "Username to authenticate with")
	password := flag.String("password", "", "Password to authenticate with")
	year := flag.Int("year", 2026, "Year to pull data for")
	month := flag.Int("month", 1, "Month to pull data for")
	day := flag.Int("day", 1, "Day to pull data for")

	flag.Parse()

	// Create the client.
	c := gocronometer.NewClient(nil)

	// Login to cronometer.
	err := c.Login(context.Background(), *username, *password)
	if err != nil {
		fmt.Println("failed to login with valid creds: %s", err)
	}

	// this is the date to extract or infor for.
	extractDate := time.Date(*year, time.Month(*month), *day, 0, 0, 0, 0, time.UTC)

	// Retrieve the export data.
	var rawCSVData string
	rawCSVData, err = c.ExportDailyNutrition(context.Background(), extractDate, extractDate)
	if err != nil {
		fmt.Println("failed to retrieve servings: %s", err)
	}

	//	fmt.Println(rawCSVData)
	csvReader := csv.NewReader(strings.NewReader(rawCSVData))

	// first record has all the headers
	headers, err := csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}
	rowData, err := csvReader.Read()
	if err != nil {
		log.Fatal(err) // EOF
	}

	// Extract only the named columns
	dataMap := make(map[string]float64)

	for col := 0; col < len(headers); col++ {
		var f float64
		f, err = strconv.ParseFloat(rowData[col], 64)
		//if err != nil {
		//	fmt.Println("Returned value for %s is not a float - %s", rowData[col], err)
		//}

		dataMap[headers[col]] = f
	}
	//	fmt.Println(dataMap)

	// get exercise info
	rawCSVData, err = c.ExportExercises(context.Background(), extractDate, extractDate)
	//	fmt.Println(rawCSVData)

	var exerciseTotalMins float64
	var exerciseTotalCals float64

	csvReader = csv.NewReader(strings.NewReader(rawCSVData))

	// ignore headers
	headers, err = csvReader.Read()

	//Row data is in form Day,Group,Exercise,Minutes,Calories Burned
	for {
		rowData, err = csvReader.Read()
		if err != nil {
			break // exit
		}
		var val float64

		val, err = strconv.ParseFloat(rowData[3], 64)
		if err != nil {
			fmt.Println("Returned value for %s is not a float - %s", rowData[3], err)
		}
		exerciseTotalMins += val
		val, err = strconv.ParseFloat(rowData[4], 64)
		if err != nil {
			fmt.Println("Returned value for %s is not a float - %s", rowData[3], err)
		}
		exerciseTotalCals += val

	}

	dataMap["Exercise Minute"] = exerciseTotalMins
	dataMap["Exercise Cals"] = exerciseTotalCals

	var jsonStr []byte
	jsonStr, err = json.Marshal(dataMap)
	if err != nil {
		fmt.Println("failed to marshal our map into json: %s", err)
	}

	fmt.Println(string(jsonStr))

}
