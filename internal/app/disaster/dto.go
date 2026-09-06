package disaster

import "time"

type CreateDisasterRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	Type        string     `json:"type" binding:"required,oneof=flood earthquake landslide"`
	StartsAt    *time.Time `json:"starts_at"`
	EndsAt      *time.Time `json:"ends_at"`
}

type DisasterResponse struct {
	ID          int32   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	StartsAt    *string `json:"starts_at,omitempty"`
	EndsAt      *string `json:"ends_at,omitempty"`
}

func toDisasterResponse(d *Disaster) DisasterResponse {
	resp := DisasterResponse{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Type:        d.Type,
		Status:      d.Status,
	}
	if d.StartsAt != nil {
		s := d.StartsAt.Format(time.RFC3339)
		resp.StartsAt = &s
	}
	if d.EndsAt != nil {
		s := d.EndsAt.Format(time.RFC3339)
		resp.EndsAt = &s
	}
	return resp
}
