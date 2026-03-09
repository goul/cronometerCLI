#cronometerCLI

Simple command line wrapper around https://github.com/jrmycanady/gocronometer that pulls data from cronometer into a json response.
Intention is to then use this data from withing homeassistant to automate some nutrition tracking tools.

Usage :

-username <cronometer user email> -password=<cronometer password> 
-year=2026 -month=3 -day=9

year/month/day is the day of data to pull. This returns a json structure with all nutrition info within it and in addition exercise duration and exercise calories for the day.

example response :

{"Added Sugars (g)":0,"Alcohol (g)":0,"B1 (Thiamine) (mg)":0.13,"B12 (Cobalamin) (µg)":0.49,"B2 (Riboflavin) (mg)":0.34,"B3 (Niacin) (mg)":1.27,"B5 (Pantothenic Acid) (mg)":0.98,"B6 (Pyridoxine) (mg)":0.34,"Caffeine (mg)":0,"Calcium (mg)":74.27,"Carbs (g)":119.11,"Cholesterol (mg)":164.12,"Completed":0,"Copper (mg)":0.21,"Cystine (g)":0.23,"Date":0,"Energy (kcal)":996.06,"Exercise Cals":-135.26,"Exercise Minute":36.00204013387362,"Fat (g)":35.94,"Fiber (g)":14.24,"Folate (µg)":56.26,"Histidine (g)":0.24,"Iron (mg)":2.51,"Isoleucine (g)":0.48,"Leucine (g)":0.8,"Lysine (g)":0.62,"Magnesium (mg)":75.91,"Manganese (mg)":0.86,"Methionine (g)":0.26,"Monounsaturated (g)":4.28,"Net Carbs (g)":104.76,"Omega-3 (g)":1.91,"Omega-6 (g)":1.09,"Phenylalanine (g)":0.54,"Phosphorus (mg)":247.64,"Polyunsaturated (g)":5.43,"Potassium (mg)":346.58,"Protein (g)":58.16,"Saturated (g)":14.62,"Selenium (µg)":19.89,"Sodium (mg)":678.27,"Starch (g)":14.04,"Sugars (g)":31.87,"Threonine (g)":0.42,"Trans-Fats (g)":0.01,"Tryptophan (g)":0.15,"Tyrosine (g)":0.36,"Valine (g)":0.56,"Vitamin A (µg)":144.25,"Vitamin C (mg)":72.75,"Vitamin D (IU)":38.28,"Vitamin E (mg)":1.35,"Vitamin K (µg)":3.64,"Water (g)":108.3,"Zinc (mg)":1.62}
