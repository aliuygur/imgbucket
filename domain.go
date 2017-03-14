package imgbucket

type (
	// Image entity
	Image struct {
		ID        int64  `bson:"id,omitempty" json:"id,omitempty"`
		Name      string `bson:"name,omitempty" json:"name,omitempty"`
		Width     int    `bson:"width,omitempty" json:"width,omitempty"`
		Height    int    `bson:"height,omitempty" json:"height,omitempty"`
		Format    string `bson:"format,omitempty" json:"format,omitempty"`
		BucketID  int64  `bson:"bucket_id,omitempty" json:"bucket_id,omitempty"`
		Size      int    `bson:"size,omitempty" json:"size,omitempty"`
		IsDefault bool   `bson:"is_default,omitempty" json:"is_default,omitempty"`
	}

	// Bucket entity
	Bucket struct {
		ID     int64  `bson:"id,omitempty" json:"id,omitempty"`
		UserID int64  `bson:"user_id" json:"user_id,omitempty"`
		Name   string `bson:"name" json:"name,omitempty"`
	}
)
