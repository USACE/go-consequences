package hazard_providers

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/USACE/go-consequences/hazards"
)

//"C:\Users\Q0HECWPL\Documents\NSI\NSI_Fathom_depths\NSI_Fathom_depths.csv"
type FrequencyData struct {
	fluvial bool      //false is pluvial
	year    int       //2020, 2050
	Values  []float64 //5yr,20yr,100yr,250yr,500yr
}
type DataSet struct {
	Data map[string]Record
}
type Record struct {
	Fd_id          string
	FutureFluvial  FrequencyData
	FuturePluvial  FrequencyData
	CurrentFluvial FrequencyData
	CurrentPluvial FrequencyData
}
type FathomEvent struct {
	Year      int
	Fluvial   bool
	Frequency int //5,20,100,250,500
}
type FathomQuery struct {
	Fd_id string
	FathomEvent
}

func (ds DataSet) ProvideHazard(args interface{}) (interface{}, error) {
	fd_id, ok := args.(FathomQuery)
	if ok {
		r, found := ds.Data[fd_id.Fd_id]
		if found {
			if fd_id.Fluvial {
				if fd_id.Year == 2020 {
					return generateDepthEvent(fd_id.Frequency, r.CurrentFluvial)
				} else if fd_id.Year == 2050 {
					return generateDepthEvent(fd_id.Frequency, r.FutureFluvial)
				} else {
					//throw error?
					return nil, HazardError{Input: "Bad Year Argument"}
				}

			} else {
				if fd_id.Year == 2020 {
					return generateDepthEvent(fd_id.Frequency, r.CurrentPluvial)
				} else if fd_id.Year == 2050 {
					return generateDepthEvent(fd_id.Frequency, r.FuturePluvial)
				} else {
					//throw error?
					return nil, HazardError{Input: "Bad Year Argument"}
				}
			}
		} else {
			return nil, NoHazardFoundError{Input: fd_id.Fd_id}
		}
	} else {
		return nil, HazardError{Input: "Unable to parse args to hazard_providers.FathomQuery"}
	}
}
func generateDepthEvent(frequency int, data FrequencyData) (hazards.DepthEvent, error) {
	switch frequency {
	case 5:
		return hazards.DepthEvent{Depth: data.Values[0]}, nil
	case 20:
		return hazards.DepthEvent{Depth: data.Values[1]}, nil
	case 100:
		return hazards.DepthEvent{Depth: data.Values[2]}, nil
	case 250:
		return hazards.DepthEvent{Depth: data.Values[3]}, nil
	case 500:
		return hazards.DepthEvent{Depth: data.Values[4]}, nil
	default:
		return hazards.DepthEvent{}, NoFrequencyFoundError{Input: fmt.Sprintf("%v", frequency)} //bad frequency
	}
}
func ConvertFile(file string) DataSet {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return DataSet{}
	}
	scanner := bufio.NewScanner(f)
	if err != nil {
		return DataSet{}
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
	ds := DataSet{Data: m}
	return ds
}
func ReadFeetFile(file string) DataSet {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		return DataSet{}
	}
	scanner := bufio.NewScanner(f)
	if err != nil {
		return DataSet{}
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
		futurefluvial := FrequencyData{fluvial: true, year: 2050, Values: ffvals}
		futurepluvial := FrequencyData{fluvial: false, year: 2050, Values: fpvals}
		currentfluvial := FrequencyData{fluvial: true, year: 2020, Values: cfvals}
		currentpluvial := FrequencyData{fluvial: false, year: 2020, Values: cpvals}
		if hasNonZeroValues(ffvals, fpvals, cfvals, cpvals) {
			//if hasValidData(fd_id, ffvals, fpvals, cfvals, cpvals) {
			r := Record{Fd_id: fd_id, FutureFluvial: futurefluvial, FuturePluvial: futurepluvial, CurrentFluvial: currentfluvial, CurrentPluvial: currentpluvial}
			m[fd_id] = r
			count++
			//} else {
			//skipping
			//}
		} else {
			//skipping.
		}

	}
	fmt.Println(count)
	ds := DataSet{Data: m}
	return ds
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
func hasValidData(fd_id string, ffvals []float64, fpvals []float64, cfvals []float64, cpvals []float64) bool {
	//ff
	ffvalid := true
	fpvalid := true
	cfvalid := true
	cpvalid := true
	datasetvalid := true
	for i := 1; i < 5; i++ {
		if ffvals[i] < ffvals[i-1] {
			ffvalid = false
			datasetvalid = false
		}
		if fpvals[i] < fpvals[i-1] {
			fpvalid = false
			datasetvalid = false
		}
		if cfvals[i] < cfvals[i-1] {
			cfvalid = false
			datasetvalid = false
		}
		if cpvals[i] < cpvals[i-1] {
			cpvalid = false
			datasetvalid = false
		}
	}
	if !datasetvalid {
		s := fd_id + " is not valid for: "
		if !ffvalid {
			s += "future fluvial,"
		}
		if !fpvalid {
			s += "future pluvial,"
		}
		if !cfvalid {
			s += "current fluvial,"
		}
		if !cpvalid {
			s += "current pluvial,"
		}
		s = strings.Trim(s, ",")
		fmt.Println(s)
		return false
	}
	return true
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func WriteBackToDisk(ds DataSet) {
	f, err := os.Create("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv")
	check(err)
	defer f.Close()
	//write header.
	//FD_ID,fluv_2020_5yr,pluv_2020_5yr,fluv_2020_20yr,pluv_2020_20yr,fluv_2020_100yr,pluv_2020_100yr,fluv_2020_250yr,pluv_2020_250yr,fluv_2020_500yr,pluv_2020_500yr,fluv_2050_5yr,pluv_2050_5yr,fluv_2050_20yr,pluv_2050_20yr,fluv_2050_100yr,pluv_2050_100yr,fluv_2050_250yr,pluv_2050_250yr,fluv_2050_500yr,pluv_2050_500yr
	w := bufio.NewWriter(f)
	w.WriteString("FD_ID,fluv_2020_5yr,pluv_2020_5yr,fluv_2020_20yr,pluv_2020_20yr,fluv_2020_100yr,pluv_2020_100yr,fluv_2020_250yr,pluv_2020_250yr,fluv_2020_500yr,pluv_2020_500yr,fluv_2050_5yr,pluv_2050_5yr,fluv_2050_20yr,pluv_2050_20yr,fluv_2050_100yr,pluv_2050_100yr,fluv_2050_250yr,pluv_2050_250yr,fluv_2050_500yr,pluv_2050_500yr\n")
	w.Flush()
	size := len(ds.Data)
	count := 0
	for _, r := range ds.Data {
		s := r.Fd_id + ","
		for i := 0; i < 5; i++ {
			s += fmt.Sprintf("%f", r.CurrentFluvial.Values[i]) + ","
			s += fmt.Sprintf("%f", r.CurrentPluvial.Values[i]) + ","
		}
		for i := 0; i < 5; i++ {
			s += fmt.Sprintf("%f", r.FutureFluvial.Values[i]) + ","
			s += fmt.Sprintf("%f", r.FuturePluvial.Values[i]) + ","
		}
		s = strings.Trim(s, ",")
		if count <= size-1 {
			s += "\n"
		}
		count++
		w.WriteString(s)
		w.Flush()
	}

}
