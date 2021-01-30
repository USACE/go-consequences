package structures

import (
	"github.com/USACE/go-consequences/paireddata"
)

type CoastalOccupancyType struct {
	Name            string
	Structuredamfun paireddata.ValueSampler
}

func CoastalOccupancyTypeMap() map[string]CoastalOccupancyType {
	m := make(map[string]CoastalOccupancyType)
	m["SFH-1SDFFI"] = sfh1sdffi()
	m["SFH-1SDFSI"] = sfh1sdfsi()
	m["SFH-1SDFMW"] = sfh1sdfmw()
	m["SFH-1SDFHW"] = sfh1sdfhw()
	m["SFH-2SDFFI"] = sfh2sdffi()
	m["SFH-2SDFSI"] = sfh2sdfsi()
	m["SFH-2SDFMW"] = sfh2sdfmw()
	m["SFH-2SDFHW"] = sfh2sdfhw()
	m["SFH-1SSFFI"] = sfh1ssffi()
	m["SFH-1SSFSI"] = sfh1ssfsi()
	m["SFH-1SSFMW"] = sfh1ssfmw()
	m["SFH-1SSFHW"] = sfh1ssfhw()
	m["SFH-2SSFFI"] = sfh2ssffi()
	m["SFH-2SSFSI"] = sfh2ssfsi()
	m["SFH-2SSFMW"] = sfh2ssfmw()
	m["SFH-2SSFHW"] = sfh2ssfhw()
	m["SFH-1SBFIU"] = sfh1sbfiu()
	m["SFH-1SBSIU"] = sfh1sbsiu()
	m["SFH-1SBMWU"] = sfh1sbmwu()
	m["SFH-1SBHWU"] = sfh1sbhwu()
	m["SFH-1SBFIF"] = sfh1sbfif()
	m["SFH-1SBSIF"] = sfh1sbsif()
	m["SFH-1SBMWF"] = sfh1sbmwf()
	m["SFH-1SBHWF"] = sfh1sbhwf()
	m["SFH-2SBFIU"] = sfh2sbfiu()
	m["SFH-2SBSIU"] = sfh2sbsiu()
	m["SFH-2SBMWU"] = sfh2sbmwu()
	m["SFH-2SBHWU"] = sfh2sbhwu()
	m["SFH-2SBFIF"] = sfh2sbfif()
	m["SFH-2SBSIF"] = sfh2sbsif()
	m["SFH-2SBMWF"] = sfh2sbmwf()
	m["SFH-2SBHWF"] = sfh2sbhwf()

	return m
}

func sfh1sdffi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 9, 22, 30, 34, 39, 43, 48, 51, 54, 57, 59, 61, 63, 64, 66, 68, 69}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SDFFI", Structuredamfun: structuredamagefunction}
}

func sfh1sdfsi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 11, 29, 38, 44, 51, 56, 63, 66, 71, 75, 77, 79, 81, 84, 86, 88, 89}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SDFSI", Structuredamfun: structuredamagefunction}
}

func sfh1sdfmw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{2, 2, 3, 4, 14, 32, 42, 48, 56, 61, 68, 72, 77, 81, 84, 86, 89, 91, 94, 96, 97}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SDFMW", Structuredamfun: structuredamagefunction}
}

func sfh1sdfhw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{5, 8, 10, 12, 20, 38, 50, 58, 66, 73, 82, 86, 92, 97, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SDFHW", Structuredamfun: structuredamagefunction}
}

func sfh2sdffi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 6, 17, 22, 26, 29, 32, 36, 38, 41, 43, 44, 46, 47, 48, 50, 51, 52}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SDFFI", Structuredamfun: structuredamagefunction}
}

func sfh2sdfsi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 8, 22, 29, 33, 38, 42, 47, 50, 53, 56, 58, 60, 61, 63, 65, 66, 67}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SDFSI", Structuredamfun: structuredamagefunction}
}

func sfh2sdfmw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 1, 10, 23, 33, 41, 50, 57, 61, 65, 69, 72, 75, 77, 80, 82, 84, 85, 87}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SDFMW", Structuredamfun: structuredamagefunction}
}

func sfh2sdfhw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 2, 2, 13, 27, 43, 60, 76, 90, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SDFHW", Structuredamfun: structuredamagefunction}
}

func sfh1ssffi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 9, 22, 30, 34, 39, 43, 48, 51, 54, 57, 59, 61, 63, 64, 66, 68, 69}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SSFFI", Structuredamfun: structuredamagefunction}
}

func sfh1ssfsi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 11, 29, 38, 44, 51, 56, 63, 66, 71, 75, 77, 79, 81, 84, 86, 88, 89}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SSFSI", Structuredamfun: structuredamagefunction}
}

func sfh1ssfmw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{3, 4, 5, 8, 22, 37, 46, 53, 60, 66, 72, 77, 81, 85, 87, 90, 92, 95, 97, 99, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SSFMW", Structuredamfun: structuredamagefunction}
}

func sfh1ssfhw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{8, 10, 12, 20, 38, 50, 58, 66, 73, 82, 86, 92, 97, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SSFHW", Structuredamfun: structuredamagefunction}
}

func sfh2ssffi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 6, 17, 22, 26, 29, 32, 36, 38, 41, 43, 44, 46, 47, 48, 50, 51, 52}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SSFFI", Structuredamfun: structuredamagefunction}
}

func sfh2ssfsi() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 0, 0, 8, 22, 29, 33, 38, 42, 47, 50, 53, 56, 58, 60, 61, 63, 65, 66, 67}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SSFSI", Structuredamfun: structuredamagefunction}
}

func sfh2ssfmw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 1, 1, 6, 16, 30, 41, 50, 59, 65, 69, 72, 76, 80, 82, 85, 87, 89, 91, 92, 94}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SSFMW", Structuredamfun: structuredamagefunction}
}

func sfh2ssfhw() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 2, 2, 13, 27, 43, 60, 76, 90, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SSFHW", Structuredamfun: structuredamagefunction}
}

func sfh1sbfiu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 3, 5, 14, 27, 35, 39, 44, 48, 53, 56, 59, 62, 64, 66, 68, 69, 71, 73, 74}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBFIU", Structuredamfun: structuredamagefunction}
}

func sfh1sbsiu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 5, 7, 16, 34, 43, 49, 56, 61, 68, 71, 76, 80, 82, 84, 86, 89, 91, 93, 94}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBSIU", Structuredamfun: structuredamagefunction}
}

func sfh1sbmwu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{8, 9, 13, 16, 29, 44, 53, 60, 67, 74, 80, 84, 89, 92, 93, 95, 96, 97, 98, 99, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBMWU", Structuredamfun: structuredamagefunction}
}

func sfh1sbhwu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{16, 18, 20, 25, 43, 55, 63, 71, 78, 87, 91, 97, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBHWU", Structuredamfun: structuredamagefunction}
}

func sfh1sbfif() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 13, 15, 24, 37, 45, 49, 54, 58, 63, 66, 69, 72, 74, 76, 78, 79, 81, 83, 84}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBFIF", Structuredamfun: structuredamagefunction}
}

func sfh1sbsif() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 15, 17, 26, 44, 53, 59, 66, 71, 78, 81, 86, 90, 92, 94, 96, 99, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBSIF", Structuredamfun: structuredamagefunction}
}

func sfh1sbmwf() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{8, 9, 18, 21, 34, 49, 58, 65, 72, 79, 85, 89, 94, 97, 98, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBMWF", Structuredamfun: structuredamagefunction}
}

func sfh1sbhwf() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{16, 18, 20, 25, 43, 55, 63, 71, 78, 87, 91, 97, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-1SBHWF", Structuredamfun: structuredamagefunction}
}

func sfh2sbfiu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 2, 3, 9, 20, 25, 29, 32, 35, 39, 41, 44, 46, 47, 49, 50, 51, 53, 54, 55}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBFIU", Structuredamfun: structuredamagefunction}
}

func sfh2sbsiu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 5, 7, 11, 25, 32, 36, 41, 45, 50, 53, 56, 59, 61, 63, 64, 66, 68, 69, 70}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBSIU", Structuredamfun: structuredamagefunction}
}

func sfh2sbmwu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{6, 7, 11, 14, 23, 38, 50, 60, 70, 77, 79, 81, 82, 84, 84, 85, 86, 87, 88, 89, 89}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBMWU", Structuredamfun: structuredamagefunction}
}

func sfh2sbhwu() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{12, 14, 16, 21, 35, 51, 68, 84, 98, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBHWU", Structuredamfun: structuredamagefunction}
}

func sfh2sbfif() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 12, 13, 19, 30, 35, 39, 42, 45, 49, 51, 54, 56, 57, 59, 60, 61, 63, 64, 65}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBFIF", Structuredamfun: structuredamagefunction}
}

func sfh2sbsif() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{0, 0, 15, 17, 21, 35, 42, 46, 51, 55, 60, 63, 66, 69, 71, 73, 74, 76, 78, 79, 80}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBSIF", Structuredamfun: structuredamagefunction}
}

func sfh2sbmwf() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{6, 7, 16, 19, 28, 43, 55, 65, 75, 82, 84, 86, 87, 89, 89, 90, 91, 92, 93, 94, 94}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBMWF", Structuredamfun: structuredamagefunction}
}

func sfh2sbhwf() CoastalOccupancyType {
	structurexs := []float64{-4.0, -3.0, -2.0, -1.0, 0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0}
	structureys := []float64{12, 14, 16, 21, 35, 51, 68, 84, 98, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100, 100}
	var structuredamagefunction = paireddata.PairedData{Xvals: structurexs, Yvals: structureys}
	return CoastalOccupancyType{Name: "SFH-2SBHWF", Structuredamfun: structuredamagefunction}
}
