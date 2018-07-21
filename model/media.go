package model

// Media is a comment's media
type Media struct {
	MIME string `json:"mime" bson:"MIME"`
	Name string `json:"name" bson:"name"`
	Size int64  `json:"size" bson:"size"`

	Storage            string `json:"-" bson:"storage"`
	StorageID          string `json:"-" bson:"storageId"`
	ThumbnailStorageID string `json:"-" bson:"thumbStorageId"`

	URL          string `json:"url" bson:"url"`
	ThumbnailURL string `json:"thumbnailUrl" bson:"thumbnailUrl"`
}
