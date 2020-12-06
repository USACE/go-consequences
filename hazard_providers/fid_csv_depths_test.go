package hazard_providers

import (
	"fmt"
	"testing"

	"github.com/USACE/go-consequences/hazards"
)

func TestOpen(t *testing.T) {
	ConvertFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths.csv")
}
func TestFeetFile(t *testing.T) {
	ReadFeetFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv")
}
func TestConvertToSqlite(t *testing.T) {
	ReadFeetFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths_Filtered_Feet.csv").WriteToSqlite()
}
func TestReadSqliteEvent(t *testing.T) {
	db := OpenSQLDepthDataSet()
	fe := FathomEvent{Year: 2050, Frequency: 500, Fluvial: true}
	fq := FathomQuery{Fd_id: "9856109", FathomEvent: fe}
	h, _ := db.ProvideHazard(fq)
	depthevent, _ := h.(hazards.DepthEvent)
	fmt.Println(depthevent.Depth)
}
func TestWrite(t *testing.T) {
	WriteBackToDisk(DataSet{})
}
func TestConvert(t *testing.T) {
	WriteBackToDisk(ConvertFile("C:\\Users\\Q0HECWPL\\Documents\\NSI\\NSI_Fathom_depths\\NSI_Fathom_depths.csv"))
}
