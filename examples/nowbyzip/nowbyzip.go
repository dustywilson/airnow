package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/dustywilson/airnow"
)

func main() {
	an := airnow.New(os.Getenv("AIRNOW"))
	ob, err := an.NowByZIP("98501", 25)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	t := template.Must(template.New("out").Parse(`{{"" -}}
Time:     {{ .Time }}
Area:     {{ .Area }}
State:    {{ .State }}
LatLng:   {{ .LatLng.Latitude }},{{ .LatLng.Longitude }}
AQI:      {{ .AQI }}
Category: {{ .Category.Num }}: {{ .Category.Name }} [{{ .Category.Color }}]
  {{- ""}}`))
	err = t.Execute(os.Stdout, ob)
	fmt.Println()
	if err != nil {
		fmt.Println(err)
	}
}
