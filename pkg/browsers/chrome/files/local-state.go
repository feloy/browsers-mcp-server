package files

import (
	"encoding/json"
	"path/filepath"

	"github.com/feloy/mcp-server/pkg/system"
)

type LocalState struct {
	Profile LocalStateProfile `json:"profile"`
}

type LocalStateProfile struct {
	ProfilesOrder []string `json:"profiles_order"`
}

func ReadLocalState() (*LocalState, error) {
	path := filepath.Join(getUserDataDirecory(), "Local State")
	jsonData, err := system.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var localState LocalState
	err = json.Unmarshal(jsonData, &localState)
	if err != nil {
		return nil, err
	}
	return &localState, nil
}
