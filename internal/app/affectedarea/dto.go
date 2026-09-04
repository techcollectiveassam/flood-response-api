package affectedarea

type CreateAffectedAreaRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	DisasterID  string `json:"disaster_id" binding:"required"`
	Location    string `json:"location"`
	Severity    string `json:"severity" binding:"required"`
	Source      string `json:"source"`
}

type AffectedAreaResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	DisasterID  string `json:"disaster_id"`
	Location    string `json:"location,omitempty"`
	Severity    string `json:"severity"`
	Source      string `json:"source,omitempty"`
}
