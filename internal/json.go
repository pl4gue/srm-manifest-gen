package srm

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"path/filepath"
)

type manifest struct {
	Title                  string `json:"title"`
	Target                 string `json:"target"`
	StartIn                string `json:"startIn"`
	LaunchOptions          string `json:"launchOptions"`
	AppendArgsToExecutable bool   `json:"appendArgsToExecutable"`
}

// Maps a game and it's executable
type Games map[string]string

func (self Games) JSONify() ([]byte, error) {
	var manifests []manifest

	for game, exec := range self {
		manifests = append(manifests, manifest{
			Title:                  game,
			Target:                 exec,
			StartIn:                filepath.Dir(exec),
			LaunchOptions:          "",
			AppendArgsToExecutable: false,
		})
	}

	json_string, err := json.MarshalIndent(manifests, "", "\t")
	if err != nil {
    log.With("err", err).Debugf("Encountered an error while parsing to JSON.")
	}

	return json_string, err
}
