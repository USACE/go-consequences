package crops

type CropType interface{
	GetCropID() byte
	GetCropCategory() string
}

