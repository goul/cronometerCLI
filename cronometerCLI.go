package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/jrmycanady/gocronometer"
)

func connectToMQTT(mqttServer string, mqttUser string, mqttPass string) mqtt.Client {

	opts := mqtt.NewClientOptions()
	opts.AddBroker(mqttServer)
	opts.SetClientID("cronometer-mqtt-link")
	opts.SetUsername(mqttUser)
	opts.SetPassword(mqttPass)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
func extractNamesValsAndUnits(headerNames []string, values []string) ([]string, []string, []string) {
	titles := []string{}
	vals := []string{}
	units := []string{}

	for i := 0; i < len(headerNames); i++ {
		// titles are of the form "blah (sdds) (unit)
		// we want the blah and the unit - lose all the stuff in brackets

		full := headerNames[i]

		measurementStart := strings.LastIndex(full, "(")
		end := strings.LastIndex(full, ")")
		if measurementStart != -1 || end != -1 {
			unit := full[measurementStart+1 : end]
			units = append(units, unit)
			title := full[0 : measurementStart-1]
			title = strings.ReplaceAll(strings.ToLower(title), "(", "")
			title = strings.ReplaceAll(title, ")", "")
			title = strings.ReplaceAll(title, " ", "_")
			title = strings.ReplaceAll(title, "-", "_")
			titles = append(titles, title)
			if values[i] != "" {
				vals = append(vals, values[i])
			} else {
				vals = append(vals, "0")
			}

		} else {
			log.Println("Header was a little odd %s - unable to find units in brackets", full)
		}
	}
	return titles, vals, units
}

func sendDataMessage(mqttClient mqtt.Client, itemNames []string, itemValues []string) {
	topic := "health/cronometer/state"
	message := buildDataMessage(itemNames, itemValues)

	token := mqttClient.Publish(topic, 0, true, message)
	token.Wait()
}
func buildDataMessage(titles []string, values []string) string {
	messageBuilder := strings.Builder{}
	messageBuilder.WriteString("{\n")

	for i := 0; i < len(titles); i++ {
		messageBuilder.WriteString("  \"" + titles[i] + "\": " + values[i])
		if i < len(titles)-1 {
			messageBuilder.WriteString(",\n")
		} else {
			messageBuilder.WriteString("\n")
		}
	}

	messageBuilder.WriteString("}\n")
	return messageBuilder.String()
}
func sendConfigMessage(mqttClient mqtt.Client, title string, unit string) {

	topic := "homeassistant/sensor/cronometer/" + title + "/config"

	messageBuilder := strings.Builder{}
	messageBuilder.WriteString("{\n")
	messageBuilder.WriteString("\"name\": \"" + title + "\",\n")
	messageBuilder.WriteString("\"state_topic\": \"health/cronometer/state\",\n")
	messageBuilder.WriteString("\"unit_of_measurement\": \"" + unit + "\",\n")
	messageBuilder.WriteString("\"value_template\": \"{{ value_json." + title + " }}\",\n")
	messageBuilder.WriteString("\"unique_id\": \"" + title + "\",\n")
	messageBuilder.WriteString("\"device\": {\n")
	messageBuilder.WriteString("		\"identifiers\": [\"cronometer\"],\n")
	messageBuilder.WriteString("		\"name\": \"Cronometer\",\n")
	messageBuilder.WriteString("		\"manufacturer\": \"Custom\",\n")
	messageBuilder.WriteString("		\"model\": \"Food tracker\"\n")
	messageBuilder.WriteString("		}\n")
	messageBuilder.WriteString("}\n")

	token := mqttClient.Publish(topic, 0, true, messageBuilder.String())
	token.Wait()

}

func main() {
	username := flag.String("username", "", "Username to authenticate with")
	password := flag.String("password", "", "Password to authenticate with")
	year := flag.Int("year", 2026, "Year to pull data for")
	month := flag.Int("month", 1, "Month to pull data for")
	day := flag.Int("day", 1, "Day to pull data for")

	mqttServer := flag.String("mqttServer", "tcp://localhost:1883", "MQTT Server URL")
	mqttUser := flag.String("mqttUser", "", "MQTT User")
	mqttPass := flag.String("mqttPassword", "", "MQTT Password")

	flag.Parse()

	// create connection to mqtt server
	mqttClient := connectToMQTT(*mqttServer, *mqttUser, *mqttPass)

	// Create the client to cronometer.
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

	// now the actual data
	rowData, err := csvReader.Read()
	if err != nil {
		log.Fatal(err) // EOF
	}
	itemNames, itemValues, itemUnits := extractNamesValsAndUnits(headers, rowData)

	// get the exercise info
	// get exercise info
	rawCSVData, err = c.ExportExercises(context.Background(), extractDate, extractDate)

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
	// round them...
	exerciseTotalMins = math.Round(exerciseTotalMins*100) / 100
	exerciseTotalCals = math.Round(exerciseTotalCals*100) / 100

	// add these to the general
	itemNames = append(itemNames, "exercise_duration")
	itemValues = append(itemValues, strconv.FormatFloat(exerciseTotalMins, 'f', -1, 64))
	itemUnits = append(itemUnits, "min")

	itemNames = append(itemNames, "exercise_total_calories")
	itemValues = append(itemValues, strconv.FormatFloat(exerciseTotalCals, 'f', -1, 64))
	itemUnits = append(itemUnits, "kcal")

	//fmt.Println(itemNames)
	//fmt.Println(itemUnits)
	//fmt.Println(itemValues)

	for i := 0; i < len(itemNames); i++ {
		sendConfigMessage(mqttClient, itemNames[i], itemUnits[i])
	}
	sendDataMessage(mqttClient, itemNames, itemValues)

	// need a little pause on disconnect to allow messages to all be sent......
	mqttClient.Disconnect(2000)
}
