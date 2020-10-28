package hazard_providers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//"C:\Users\Q0HECWPL\Documents\NSI\NSI_Fathom_depths\NSI_Fathom_depths.csv"
type FrequencyData struct {
	fluvial bool      //false is pluvial
	year    int       //2020, 2050
	values  []float64 //5yr,20yr,100yr,250yr,500yr
}
type Record struct {
	Fd_id          string
	FutureFluvial  FrequencyData
	FuturePluvial  FrequencyData
	CurrentFluvial FrequencyData
	CurrentPluvial FrequencyData
}

func ConvertFile(file string) map[string]Record {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return nil
	}
	scanner := bufio.NewScanner(f)
	if err != nil {
		return nil
	}
	scanner.Scan()
	fmt.Println(scanner.Text()) //header row
	m := make(map[string]Record)
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), ",")
		fd_id := lines[0]
		//fluv_2020_5yr,pluv_2020_5yr,fluv_2020_20yr,pluv_2020_20yr,fluv_2020_100yr,pluv_2020_100yr,fluv_2020_250yr,pluv_2020_250yr,fluv_2020_500yr,pluv_2020_500yr,fluv_2050_5yr,pluv_2050_5yr,fluv_2050_20yr,pluv_2050_20yr,fluv_2050_100yr,pluv_2050_100yr,fluv_2050_250yr,pluv_2050_250yr,fluv_2050_500yr,pluv_2050_500yr
		fluvial := true
		cfvals := make([]float64, 5)
		cpvals := make([]float64, 5)
		ffvals := make([]float64, 5)
		fpvals := make([]float64, 5)
		twentyTwenty := 0
		fpidx := 0
		ffidx := 0
		cpidx := 0
		cfidx := 0
		for i := 1; i < len(lines); i++ {
			if twentyTwenty > 10 {
				//2050
				if fluvial {
					ffvals[ffidx], err = strconv.ParseFloat(lines[i], 64)
					ffidx++
				} else {
					fpvals[fpidx], err = strconv.ParseFloat(lines[i], 64)
					fpidx++
				}
			} else {
				//2020
				if fluvial {
					cfvals[cfidx], err = strconv.ParseFloat(lines[i], 64)
					cfidx++
				} else {
					cpvals[cpidx], err = strconv.ParseFloat(lines[i], 64)
					cpidx++
				}
			}
			fluvial = !fluvial
			twentyTwenty++
		}
		futurefluvial := FrequencyData{fluvial: true, year: 2050, values: ffvals}
		futurepluvial := FrequencyData{fluvial: false, year: 2050, values: fpvals}
		currentfluvial := FrequencyData{fluvial: true, year: 2020, values: cfvals}
		currentpluvial := FrequencyData{fluvial: false, year: 2020, values: cpvals}
		r := Record{Fd_id: fd_id, FutureFluvial: futurefluvial, FuturePluvial: futurepluvial, CurrentFluvial: currentfluvial, CurrentPluvial: currentpluvial}
		m[fd_id] = r
	}
	return m
}
