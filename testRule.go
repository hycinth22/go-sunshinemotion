package ssmt

type TestRule struct {
	ID           int64   `json:"id"`
	SchoolID     int64   `json:"schoolID"`
	ManTime      int64   `json:"manTime"`      // seconds
	ManDistance  float64 `json:"manDistance"`  // meters
	GirlTime     int64   `json:"girlTime"`     // seconds
	GirlDistance float64 `json:"girlDistance"` // meters
}
