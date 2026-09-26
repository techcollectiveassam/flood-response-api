package commandcenter

type CreateCommandCenterRequest struct {
	DisasterID    int32  `json:"disaster_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	Type          string `json:"type" binding:"required,oneof=government ngo group other"`
	Description   string `json:"description"`
	ContactPerson string `json:"contact_person"`
	ContactMobile string `json:"contact_mobile"`
	ContactEmail  string `json:"contact_email"`
}

type CommandCenterResponse struct {
	ID            int32  `json:"id"`
	DisasterID    int32  `json:"disaster_id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Description   string `json:"description,omitempty"`
	ContactPerson string `json:"contact_person,omitempty"`
	ContactMobile string `json:"contact_mobile,omitempty"`
	ContactEmail  string `json:"contact_email,omitempty"`
}

func toCommandCenterResponse(cc *CommandCenter) CommandCenterResponse {
	return CommandCenterResponse{
		ID:            cc.ID,
		DisasterID:    cc.DisasterID,
		Name:          cc.Name,
		Type:          cc.Type,
		Description:   cc.Description,
		ContactPerson: cc.ContactPerson,
		ContactMobile: cc.ContactMobile,
		ContactEmail:  cc.ContactEmail,
	}
}
