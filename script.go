package script

import (
	"net/http"
)

// TYPES

type Mod struct {
	Name         string
	Source       string
	Version      string
	Description  string
	PrettyName   string
	DownloadLink string
}

func newMod(name, source, version, description string) (*Mod, error) {
	resp, err := http.Get("https://api.modrinth.com/v2/project/" + source)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var modrinthProject struct {
		Title string `json:"title"`
	}

	var modrinthVersion struct {
	}

	mod := &Mod{
		Name:        name,
		Source:      source,
		Version:     version,
		Description: description,
	}

	return mod, nil
}
