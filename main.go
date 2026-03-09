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

	fmt.Println(rawCSVData)
	r := csv.NewReader(strings.NewReader(rawCSVData))

	// first record has all the headers
	headers, err := r.Read()
	if err != nil {
		log.Fatal(err)
	}
	colIndex := make(map[string]int)
	for i, h := range headers {
		colIndex[h] = i
	}

	// Choose the columns you want by name
	want := []string{"Energy (kcal)", "Protein (g)", "Fibre (g)", "Carbs (g)"}
	for {
		row, err := r.Read()
		if err != nil {
			break // EOF
		}

		// Extract only the named columns
		selectedMap := make(map[string]float64)
		for _, col := range want {
			idx := colIndex[col]
			var f float64
			f, err = strconv.ParseFloat(row[idx], 64)
			if err != nil {
				fmt.Println("Returned value for %s is not a float - %s", row[idx], err)
			}
			selectedMap[col] = f
		}

		fmt.Println(selectedMap)
		var jsonStr []byte
		jsonStr, err = json.Marshal(selectedMap)
		if err != nil {
			fmt.Println("failed to retrieve servings: %s", err)
		}
		fmt.Println(string(jsonStr))

		//rawCSVData, err = c.ExportExercises(context.Background(), extractDate, extractDate)
		//fmt.Println(rawCSVData)
	}

}
