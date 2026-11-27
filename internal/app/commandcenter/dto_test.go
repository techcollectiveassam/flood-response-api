package commandcenter

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var uuidZero = uuid.UUID{}

func TestCommandCenterResponseOmitsEmptyOptionalFields(t *testing.T) {
	c := &CommandCenter{
		ID:            1,
		DisasterRefID: uuidZero,
		Name:          "Minimal Center",
		Type:          TypeOther,
	}

	resp := toCommandCenterResponse(c)

	jsonBytes, err := json.Marshal(resp)
	assert.NoError(t, err)

	var decoded map[string]interface{}
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), decoded["id"])
	assert.Equal(t, uuidZero.String(), decoded["disaster_ref_id"])
	assert.Equal(t, "Minimal Center", decoded["name"])
	assert.Equal(t, "other", decoded["type"])
	assert.NotContains(t, decoded, "description")
	assert.NotContains(t, decoded, "contact_person")
	assert.NotContains(t, decoded, "contact_mobile")
	assert.NotContains(t, decoded, "contact_email")
}
