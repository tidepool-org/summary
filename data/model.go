package data

//Base is a subset of the fields common to all datums
type Base struct {
	Active   bool    `json:"-" bson:"_active"` // if false, this object has been effectively deleted
	DeviceID *string `json:"deviceId,omitempty" bson:"deviceId,omitempty"`
	ID       *string `json:"id,omitempty" bson:"id,omitempty"`
	Source   *string `json:"source,omitempty" bson:"source,omitempty"`
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
