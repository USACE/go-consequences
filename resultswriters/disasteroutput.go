package resultswriters

import (
	"fmt"
	"io"
	"os"

	"github.com/USACE/go-consequences/consequences"
	"github.com/USACE/go-consequences/hazards"
	"github.com/USACE/go-consequences/indirecteconomics"
	"github.com/USACE/go-consequences/structureprovider"
)

type disasterOuput struct {
	filepath string
	w        io.Writer
	fipsmap  map[string]countyRecord
	sp       structureprovider.StructureProvider
}
type countyRecord struct {
	statefips      string
	countyfips     string
	resDamCount    int
	nonresdamcount int
	resTotDam      float64
	indirectLosses float64 //in millions.
	workingPop     int32
	nonResTotDam   float64
	totalValue     float64
	totalDamages   float64
	byAssetType    map[string]typeRecord
}
type typeRecord struct {
	totalInCounty        int
	damageCategorization map[string]int
	thresholds           []float64
	thresholdnames       []string
	totalDamages         float64
}

func InitDisasterOutput(filepath string, sp structureprovider.StructureProvider) *disasterOuput {
	w, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		panic(err)
	}
	//make fipsmap
	fm := make(map[string]countyRecord)
	return &disasterOuput{filepath: filepath, w: w, fipsmap: fm, sp: sp}
}
func (srw *disasterOuput) Write(r consequences.Result) {
	f, ferr := r.Fetch("cbfips")

	if ferr == nil {
		fips := f.(string)
		ssccc := fips[0:5] //first five ditgits of fips are state (2) and county (3)
		cr, crok := srw.fipsmap[ssccc]
		if crok {
			//previously existed, update and reassign.
			cr.Update(r)
			srw.fipsmap[ssccc] = cr
		} else {
			am := make(map[string]typeRecord)
			cr = countyRecord{statefips: fips[0:2], countyfips: fips[2:5], byAssetType: am}
			cr.Update(r)
			srw.fipsmap[ssccc] = cr
		}
		srw.fipsmap[ssccc] = cr
	}
}
func (cr *countyRecord) Update(r consequences.Result) {
	d, derr := r.Fetch("damage category")
	if derr == nil {
		damcat := d.(string)
		if damcat == "RES" {
			cr.resDamCount += 1
			v, _ := r.Fetch("structure damage") //unsafe skipping error.
			damage := v.(float64)
			cr.resTotDam += damage
			cr.totalDamages += damage
			v, _ = r.Fetch("content damage")
			damage = v.(float64)
			cr.resTotDam += damage
			cr.totalDamages += damage
			pop, _ := r.Fetch("pop2amu65")
			cr.workingPop += pop.(int32)
		} else {
			cr.nonresdamcount += 1
			v, _ := r.Fetch("structure damage") //unsafe skipping error.
			damage := v.(float64)
			cr.nonResTotDam += damage
			cr.totalDamages += damage
			v, _ = r.Fetch("content damage")
			damage = v.(float64)
			cr.nonResTotDam += damage
			cr.totalDamages += damage
		}
	}
	o, oerr := r.Fetch("occupancy type")
	assetType := ""
	if oerr == nil {
		occtype := o.(string)
		switch occtype {
		case "REL1":
			assetType = "Assembly"
		case "AGR1":
			assetType = "Agriculture"
		case "EDU1", "EDU2":
			assetType = "Education"
		case "GOV1", "GOV2", "COM6": //why are hospitals encoded to governement?
			assetType = "Government"
		case "IND1", "IND2", "IND3", "IND4", "IND5", "IND6":
			assetType = "Industrial"
		case "COM1", "COM2", "COM3", "COM4", "COM5", "COM7", "COM8", "COM9", "COM10":
			assetType = "Commercial"
		default:
			assetType = "Residential"
		}
		at, atok := cr.byAssetType[assetType]
		if atok {
			at.Update(r)
		} else {
			//create a new one.
			dc := make(map[string]int)
			thresholds := []float64{0.0, 2.0, 4.0, 6.0, 9998.0}
			headers := []string{"No Damage (0 ft)", "Affected (<=2 ft)", "Minor Damage (2 - 4 ft)", "Major Damage (4 - 6 ft)", "Destroyed (6+ ft)"}
			at = typeRecord{damageCategorization: dc, thresholds: thresholds, thresholdnames: headers}
			at.Update(r)
		}
		cr.byAssetType[assetType] = at
	}

}
func (at *typeRecord) Update(r consequences.Result) {
	h, hok := r.Fetch("hazard")
	if hok == nil {
		he := h.(hazards.DepthEvent)
		depth := he.Depth()
		for idx, v := range at.thresholds {
			if depth <= v {
				count, cok := at.damageCategorization[at.thresholdnames[idx]]
				if cok {
					count += 1
				} else {
					count = 1
				}
				at.damageCategorization[at.thresholdnames[idx]] = count
				break
			}
		}
		v, _ := r.Fetch("structure damage") //unsafe skipping error.
		damage := v.(float64)
		at.totalDamages += damage
		v, _ = r.Fetch("content damage")
		damage = v.(float64)
		at.totalDamages += damage
	}

}
func (at typeRecord) write() string {
	return fmt.Sprintf(",%v,%v,%v,%v,%v,%v", at.totalInCounty, at.damageCategorization["Destroyed (6+ ft)"], at.damageCategorization["Major Damage (4 - 6 ft)"], at.damageCategorization["Minor Damage (2 - 4 ft)"], at.damageCategorization["Affected (<=2 ft)"], at.totalDamages)
}
func (cr countyRecord) write() string {
	summary := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v", cr.statefips, cr.countyfips, cr.resDamCount, cr.nonresdamcount, cr.resTotDam, cr.nonResTotDam, "", cr.totalValue, cr.totalDamages, cr.indirectLosses, "")
	summary += cr.byAssetType["Residential"].write()
	summary += cr.byAssetType["Agriculture"].write()
	summary += cr.byAssetType["Assembly"].write()
	summary += cr.byAssetType["Commercial"].write()
	summary += cr.byAssetType["Education"].write()
	summary += cr.byAssetType["Government"].write()
	summary += cr.byAssetType["Industrial"].write()
	//others arent mapped based on feedback from FEMA.
	return summary
}

type stateSummary struct {
	totalPopulation                  int32
	totalHomesWithMajorDamage        int32
	totalHomesDestroyed              int32
	totalInfrastructureDamge         float64
	damagesPerCapita                 float64
	totalReplacementValueOfBuildings float64
	totalPropertyDamage              float64
	totalBuisnessInteruptionLosses   float64
}

func (srw *disasterOuput) Close() {
	statelist := make(map[string]stateSummary)
	for k, v := range srw.fipsmap {

		s := structureprovider.StatsByFips(k, srw.sp)
		v.totalValue = s.TotalValue
		for ak, av := range v.byAssetType {
			av.totalInCounty = s.CountByCategory[ak]
			v.byAssetType[ak] = av
		}
		laborLossRatio := 0.0
		if s.WorkingResPop2AM != 0 {
			laborLossRatio = float64(v.workingPop) / float64(s.WorkingResPop2AM)
		}

		capitalLossRatio := 0.0
		if s.NonResidentialValue != 0.0 {
			capitalLossRatio = v.nonResTotDam / s.NonResidentialValue
		}
		//compute ecam.
		er, err := indirecteconomics.ComputeEcam(v.statefips, v.countyfips, capitalLossRatio, laborLossRatio)
		if err == nil {
			for _, pr := range er.ProductionImpacts {
				if pr.Sector == "Total" {
					v.indirectLosses = pr.Change
					fmt.Printf("State %v, County %v, Losses %v\n", v.statefips, v.countyfips, pr.Change)
					break
				}
			}
		} else {
			fmt.Println(err)
		}
		srw.fipsmap[k] = v
		state, stateok := statelist[v.statefips]
		if !stateok {
			state = stateSummary{}
		}
		state.totalHomesDestroyed += int32(v.byAssetType["Residential"].damageCategorization["Destroyed (6+ ft)"])
		state.totalHomesWithMajorDamage += int32(v.byAssetType["Residential"].damageCategorization["Major Damage (4 - 6 ft)"])
		state.totalBuisnessInteruptionLosses += v.indirectLosses
		state.totalPropertyDamage += v.totalDamages
		statelist[v.statefips] = state

	}
	//loop over states and compute population and value totals
	for k, v := range statelist {
		s := structureprovider.StatsByFips(k, srw.sp)
		v.totalPopulation = s.TotalResPop2AM //night time population is my best estimate for state population
		v.totalReplacementValueOfBuildings = s.TotalValue
		v.damagesPerCapita = v.totalPropertyDamage / float64(v.totalPopulation)
		statelist[k] = v
		//write out each state summary?
	}

	//create results.
	result := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,Residential,Residential,Residential,Residential,Residential,Residential,Agriculture,Agriculture,Agriculture,Agriculture,Agriculture,Agriculture,Assembly,Assembly,Assembly,Assembly,Assembly,Assembly,Commercial,Commercial,Commercial,Commercial,Commercial,Commercial,Education,Education,Education,Education,Education,Education,Government,Government,Government,Government,Government,Government,Industrial,Industrial,Industrial,Industrial,Industrial,Industrial\n", "", "", "", "", "", "", "", "", "", "", "")
	result += "State,County,Residential Structures Destroyed or Majorly Damaged,Non-Residential Destroyed or Majorly Damaged,Residential Damages ($),Non-Residential Damages ($),Damages to Infrastructure Per Capita,Total Building Value,Total Building Losses,Buisness Interuption Costs ($M),HomeOwnership Rate of Impacted Residential Structures,Total Structures,Destroyed,Major Damage,Minor Damages,Affected Damages,Damages ($),Total Structures,Destroyed,Major Damage,Minor Damages,Affected Damages,Damages ($),Total Structures,Destroyed,Major Damage,Minor Damages,Affected Damages,Damages ($),Total Structures,Destroyed,Major Damage,Minor Damages,Affected Damages,Damages ($),Total Structures,Destroyed,Major Damage,Minor Damages,Affected Damages,Damages ($),Total Structures,Destroyed,Major Damage,Minor Damages,Affected Damages,Damages ($),Total Structures,Destroyed,Major Damage,Minor Damages,Affected Damages,Damages ($)\n"
	for _, v := range srw.fipsmap {
		result += v.write() + "\n"
	}
	//write to file.
	srw.w.Write([]byte(result))
	w2, ok := srw.w.(io.WriteCloser)
	if ok {
		w2.Close()
	}
}
