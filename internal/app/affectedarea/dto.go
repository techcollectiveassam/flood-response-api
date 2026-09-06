package affectedarea

type CreateAffectedAreaRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	DisasterID  int32  `json:"disaster_id" binding:"required"`
	Location    string `json:"location"`
	Geometry    string `json:"geometry"`
	Severity    string `json:"severity" binding:"required,oneof=low medium high critical"`
}

type AffectedAreaResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DisasterID  int32  `json:"disaster_id"`
	Location    string `json:"location,omitempty"`
	Geometry    string `json:"geometry,omitempty"`
	Severity    string `json:"severity"`
}
