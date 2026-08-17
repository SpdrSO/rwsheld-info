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

	switch mod.Source {
	case "modrinth":
		if err := mod.initFromModrinth(); err != nil {
			return nil, fmt.Errorf("Инициализация с Modrinth: %w", err)
		}
	default:
		return nil, fmt.Errorf("Неизвестный источник модов: %q", source)
	}

	return mod, nil
}

func (m *Mod) initFromModrinth() error {
	resp, err := http.Get("https://api.modrinth.com/v2/project/" + m.Name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var modrinthProject struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modrinthProject); err != nil {
		return err
	}

	resp, err = http.Get("https://api.modrinth.com/v2/project/" + m.Name + "/version")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	m.PrettyName = modrinthProject.Title

	var modrinthVersions []struct {
		Version string `json:"version_number"`
		Files   []struct {
			URL     string `json:"url"`
			Primary bool   `json:"primary"`
		} `json:"files"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&modrinthVersions); err != nil {
		return err
	}

	target := m.Version + "+26.2"
	found := false

	for _, v := range modrinthVersions {
		if v.Version != target {
			continue
		}

		found = true

		for _, f := range v.Files {
			if f.Primary {
				m.DownloadLink = f.URL
				break
			}
		}
	}

	if !found {
		return fmt.Errorf("Версия %q для мода %q не найдена", target, m.Name)
	}

	return nil
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
