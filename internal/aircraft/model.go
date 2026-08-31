package aircraft

import "time"

type Aircraft struct {
	ID         int       `json:"id"`
	WINNum     string    `json:"win_num"`
	Model      string    `json:"model"`
	Capacity   int       `json:"capacity"`
	CarryingKg int       `json:"carrying_kg"`
	FlightTime int       `json:"flight_time"`
	ReleasedAt time.Time `json:"released_at"`
	LicensedAt time.Time `json:"licensed_at"`
}
