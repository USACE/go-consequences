package store

type fathom_result struct {
	Fd_id                 string  `db:"fd_id"`
	Hazard_Year           string  `db:"hazard_year"` //2020, 2050
	Hazard_Type           string  `db:"hazard_type"` //fluvial, pluvial
	Frequency             string  `db:"frequency"`   // 5, 20, 100, 250, 500
	Structure_Consequence float64 `db:"structure_consequence"`
	Content_Consequence   float64 `db:"content_consequence"`
}
