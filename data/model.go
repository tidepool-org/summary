package data

//TypeOnly is used to pull out only the type field
type TypeOnly struct {
	Type string `json:"type,omitempty" bson:"type,omitempty"`
}

//Base is a subset of the fields common to all datums
type Base struct {
	Active   bool    `json:"-" bson:"_active"` // if false, this object has been effectively deleted
	DeviceID *string `json:"deviceId,omitempty" bson:"deviceId,omitempty"`
	ID       *string `json:"id,omitempty" bson:"id,omitempty"`
	Time     *string `json:"time,omitempty" bson:"time,omitempty"`
	Type     string  `json:"type,omitempty" bson:"type,omitempty"`
	UploadID *string `json:"uploadId,omitempty" bson:"uploadId,omitempty"`
	UserID   *string `json:"-" bson:"_userId,omitempty"`
}

// Blood is the type of a blood value
type Blood struct {
	Base `bson:",inline"`

	Units *string  `json:"units,omitempty" bson:"units,omitempty"`
	Value *float64 `json:"value,omitempty" bson:"value,omitempty"`
}

//Client software used to upload data
type Client struct {
	Name    *string `json:"name,omitempty" bson:"name,omitempty"`
	Version *string `json:"version,omitempty" bson:"version,omitempty"`
}

//Upload describes the upload device and client software used
type Upload struct {
	Base `bson:",inline"`

	Client              *Client   `json:"client,omitempty" bson:"client,omitempty"`
	DeviceManufacturers *[]string `json:"deviceManufacturers,omitempty" bson:"deviceManufacturers,omitempty"`
	DeviceModel         *string   `json:"deviceModel,omitempty" bson:"deviceModel,omitempty"`
	DeviceSerialNumber  *string   `json:"deviceSerialNumber,omitempty" bson:"deviceSerialNumber,omitempty"`
	DeviceTags          *[]string `json:"deviceTags,omitempty" bson:"deviceTags,omitempty"`
}
