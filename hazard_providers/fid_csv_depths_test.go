package hazard_providers

import "testing"

func TestOpen(t *testing.T) {
	ConvertFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths.csv")
}
func TestFeetFile(t *testing.T) {
	ReadFeetFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv")
}

func TestWrite(t *testing.T) {
	WriteBackToDisk(DataSet{})
}
func TestConvert(t *testing.T) {
	WriteBackToDisk(ConvertFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths.csv"))
}
