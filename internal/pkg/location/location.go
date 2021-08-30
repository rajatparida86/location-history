package location

type Location struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
	createdAt int64
}
