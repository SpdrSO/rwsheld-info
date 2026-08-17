package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CONSTS

const (
	BUILD_DIR string = "build"
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
	mod := &Mod{
		Name:        name,
		Source:      source,
		Version:     version,
		Description: description,
	}

	resp, err := http.Get("https://api.modrinth.com/v2/project/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var modrinthProject struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modrinthProject); err != nil {
		return nil, err
	}

	resp, err = http.Get("https://api.modrinth.com/v2/project/" + name + "/version")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	mod.PrettyName = modrinthProject.Title

	var modrinthVersions []struct {
		Version string `json:"version_number"`
		Files   []struct {
			URL     string `json:"url"`
			Primary bool   `json:"primary"`
		} `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modrinthVersions); err != nil {
		return nil, err
	}

	target := version + "+26.2"

	for _, v := range modrinthVersions {
		if v.Version != target {
			continue
		}

		for _, f := range v.Files {
			if f.Primary {
				mod.DownloadLink = f.URL
				break
			}
		}
	}

	return mod, nil
}

// FUNCS

func main() {
	mod, err := newMod("fabric-api", "modrinth", "0.155.2", "Эт Фабрик ЭйПиАй")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Name:", mod.Name)
	fmt.Println("Source:", mod.Source)
	fmt.Println("Version:", mod.Version)
	fmt.Println("Description:", mod.Description)
	fmt.Println("PrettyName:", mod.PrettyName)
	fmt.Println("DownloadLink:", mod.DownloadLink)
}
