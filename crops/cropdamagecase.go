package crops

type CropDamageCase byte

const (
	Unassigned              CropDamageCase = 0
	Impacted                CropDamageCase = 1
	NotImpactedDuringSeason CropDamageCase = 2
	PlantingDelayed         CropDamageCase = 4
	NotPlanted              CropDamageCase = 8
	SubstituteCrop          CropDamageCase = 16
)