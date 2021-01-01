package crops

//CropType describes the relationship of the byte value in a NASSCDL to the Crop name (crop category in their terminology)
type CropType interface {
	GetCropID() byte
	GetCropCategory() string
}
