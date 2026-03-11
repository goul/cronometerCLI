

Simple command line wrapper around [gocronometer](https://github.com/jrmycanady/gocronometer) that pulls data from cronometer and pushes it to home assistant.

The sensor relies on home assistant being connected to an MQTT instance. If it is a new device will appear under MQTT called cronometer.
This device will contain sensors for all of the nutritional/fitness data pulled from cronometer. 


Usage :
```
cronometerCLI -username <cronometer user email> -password=<cronometer password>  -year=YYYY -month=M -day=D -mqttServer=<mqttServer> -mqttUser=<mqttUser> -mqttPassword=<mqqtPassword>
```
year/month/day is the day of data to pull. This returns a json structure with all nutrition info within it and in addition exercise duration and exercise calories for the day.

```
eg. cronometerCLI -username paul@goulbourn.com -password=myPassword  -year=2026 -month=3 -day=9 -mqttServer=tcp://homeassistant.local:1883/ -mqttUser=demo -mqttPassword=demo
```

