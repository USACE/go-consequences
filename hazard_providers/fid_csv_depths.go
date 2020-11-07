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
	Values  []float64 //5yr,20yr,100yr,250yr,500yr
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
	count := 0
	for scanner.Scan() {
		lines := strings.Split(scanner.Text(), ",")
		fd_id := lines[0]
		//fmt.Println(fd_id)
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
			if twentyTwenty >= 10 {
				//2050
				if fluvial {
					ffvals[ffidx], err = strconv.ParseFloat(lines[i], 64)
					ffvals[ffidx] = ffvals[ffidx] / 30.48 //centimeters to feet
					ffidx++
				} else {
					fpvals[fpidx], err = strconv.ParseFloat(lines[i], 64)
					fpvals[fpidx] = fpvals[fpidx] / 30.48 //centimeters to feet
					fpidx++
				}
			} else {
				//2020
				if fluvial {
					cfvals[cfidx], err = strconv.ParseFloat(lines[i], 64)
					//fmt.Println("current fluvial")
					cfvals[cfidx] = cfvals[cfidx] / 30.48 //centimeters to feet
					cfidx++
				} else {
					cpvals[cpidx], err = strconv.ParseFloat(lines[i], 64)
					cpvals[cpidx] = cpvals[cpidx] / 30.48 //centimeters to feet
					//fmt.Println("current pluvial")
					cpidx++
				}
			}
			fluvial = !fluvial
			twentyTwenty++
		}
		futurefluvial := FrequencyData{fluvial: true, year: 2050, Values: ffvals}
		futurepluvial := FrequencyData{fluvial: false, year: 2050, Values: fpvals}
		currentfluvial := FrequencyData{fluvial: true, year: 2020, Values: cfvals}
		currentpluvial := FrequencyData{fluvial: false, year: 2020, Values: cpvals}
		if hasNonZeroValues(ffvals, fpvals, cfvals, cpvals) {
			r := Record{Fd_id: fd_id, FutureFluvial: futurefluvial, FuturePluvial: futurepluvial, CurrentFluvial: currentfluvial, CurrentPluvial: currentpluvial}
			m[fd_id] = r
			count++
		} else {
			//skipping.
		}

	}
	fmt.Println(count)
	return m
}
func hasNonZeroValues(ffvals []float64, fpvals []float64, cfvals []float64, cpvals []float64) bool {
	for i := 0; i < 5; i++ {
		if ffvals[i] > 0 {
			return true
		}
		if fpvals[i] > 0 {
			return true
		}
		if cfvals[i] > 0 {
			return true
		}
		if cpvals[i] > 0 {
			return true
		}
	}
	return false
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func WriteBackToDisk(m map[string]Record) {
	f, err := os.Create("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv")
	check(err)
	defer f.Close()
	//write header.
	//FD_ID,fluv_2020_5yr,pluv_2020_5yr,fluv_2020_20yr,pluv_2020_20yr,fluv_2020_100yr,pluv_2020_100yr,fluv_2020_250yr,pluv_2020_250yr,fluv_2020_500yr,pluv_2020_500yr,fluv_2050_5yr,pluv_2050_5yr,fluv_2050_20yr,pluv_2050_20yr,fluv_2050_100yr,pluv_2050_100yr,fluv_2050_250yr,pluv_2050_250yr,fluv_2050_500yr,pluv_2050_500yr
	w := bufio.NewWriter(f)
	w.WriteString("FD_ID,fluv_2020_5yr,pluv_2020_5yr,fluv_2020_20yr,pluv_2020_20yr,fluv_2020_100yr,pluv_2020_100yr,fluv_2020_250yr,pluv_2020_250yr,fluv_2020_500yr,pluv_2020_500yr,fluv_2050_5yr,pluv_2050_5yr,fluv_2050_20yr,pluv_2050_20yr,fluv_2050_100yr,pluv_2050_100yr,fluv_2050_250yr,pluv_2050_250yr,fluv_2050_500yr,pluv_2050_500yr\n")
	w.Flush()
}
