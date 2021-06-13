package crops

//CropDamageCase provides a structured response for the possible outcomes from CropSchedule
type CropDamageCase byte

//CropDamageCases include Unassigned, Impacted, NotImpactedDuringSeason, PlantingDelayed, NotPlanted, or SubstituteCrop
const (
	Unassigned              CropDamageCase = 0
	Impacted                CropDamageCase = 1
	NotImpactedDuringSeason CropDamageCase = 2
	PlantingDelayed         CropDamageCase = 4
	NotPlanted              CropDamageCase = 8
	SubstituteCrop          CropDamageCase = 16
)

func (cdc CropDamageCase) String() string {
	switch cdc {
	case Unassigned:
		return "Unassigned"
	case Impacted:
		return "Impacted"
	case NotImpactedDuringSeason:
		return "Not Impacted During Season"
	case NotPlanted:
		return "Not Planted"
	case SubstituteCrop:
		return "Substitute Crop"
	default:
		return "Unassigned"
	}
}
