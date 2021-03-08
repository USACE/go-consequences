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
