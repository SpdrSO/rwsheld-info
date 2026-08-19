package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

//
//
//   CONSTS
//
//

const (
	BUILD_DIR    = "build"
	GAME_VERSION = "26.2"
	MOD_LOADER   = "fabric"
)

//
//
//   TYPES
//
//

type Mod struct {
	Name         string `toml:"name"`
	PrettyName   string `toml:"pretty_name"`
	Source       string `toml:"source"`
	Version      string `toml:"version"`
	DownloadLink string `toml:"download_link"`
	ModType      string `toml:"type"`
	Description  string `toml:"description"`
}

func (m *Mod) initMissingFields() error {
	if m.Name == "" {
		return fmt.Errorf("Для мода не указано обязательное поле \"name\"")
	}

	if m.ModType != "general" && m.ModType != "optimization" && m.ModType != "dependency" {
		return fmt.Errorf("Неизвестный тип мода: %q", m.ModType)
	}

	var g errgroup.Group

	if m.PrettyName == "" {
		g.Go(func() error {
			if err := m.initPrettyName(); err != nil {
				m.PrettyName = m.Name
			}

			return nil
		})
	}

	if m.DownloadLink == "" {
		if m.Source == "" || m.Version == "" {
			return fmt.Errorf("У мода %q отсутствует поле \"download_link\" или поля \"source\" + \"version\"", m.Name)
		}

		g.Go(func() error {
			return m.initDownloadLink()
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (m *Mod) initPrettyName() error {
	switch m.Source {
	case "modrinth":
		project, err := getModrinthProject(m.Name)
		if err != nil {
			return err
		}

		m.PrettyName = project.Title

		return nil
	case "curseforge":
		return fmt.Errorf("Источник CurseForge пока что не поддерживается")
	default:
		return fmt.Errorf("Неизвестный источник: %q", m.Source)
	}
}

func (m *Mod) initDownloadLink() error {
	switch m.Source {
	case "modrinth":
		versions, err := getModrinthVersions(m.Name)
		if err != nil {
			return err
		}

		target := m.Version + "+" + GAME_VERSION
		found := false
		for _, version := range versions {
			if version.Version != target {
				continue
			}

			loaderSupported := false
			for _, loader := range version.Loaders {
				if loader == MOD_LOADER {
					loaderSupported = true
					break
				}
			}
			if !loaderSupported {
				continue
			}

			found = true

			for _, file := range version.Files {
				if file.Primary {
					m.DownloadLink = file.URL
					break
				}
			}
		}

		if !found {
			return fmt.Errorf("Версия %q для мода %q не найдена", target, m.Name)
		}

		return nil
	case "curseforge":
		return fmt.Errorf("Источник CurseForge пока что не поддерживается")
	default:
		return fmt.Errorf("Неизвестный источник: %q", m.Source)
	}
}

//
//
//   GET
//
//

type ModrinthProject struct {
	Title string `json:"title"`
}

func getModrinthProject(name string) (*ModrinthProject, error) {
	resp, err := http.Get("https://api.modrinth.com/v2/project/" + name)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var project ModrinthProject

	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

type ModrinthVersions struct {
	Version string   `json:"version_number"`
	Loaders []string `json:"loaders"`
	Files   []struct {
		URL     string `json:"url"`
		Primary bool   `json:"primary"`
	} `json:"files"`
}

func getModrinthVersions(name string) ([]ModrinthVersions, error) {
	resp, err := http.Get("https://api.modrinth.com/v2/project/" + name + "/version")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var versions []ModrinthVersions

	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		return nil, err
	}

	return versions, nil
}

//
//
//   FUNCS
//
//

func main() {
	mod := Mod{
		Name:        "fabric-api",
		Source:      "modrinth",
		Version:     "0.155.2",
		ModType:     "general",
		Description: "Эт Фабрик ЭйПиАй",
	}

	if err := mod.initMissingFields(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при инициализации мода %q: %v", mod.Name, err)
	} else {
		fmt.Println("Name:", mod.Name)
		fmt.Println("PrettyName:", mod.PrettyName)
		fmt.Println("Source:", mod.Source)
		fmt.Println("Version:", mod.Version)
		fmt.Println("DownloadLink:", mod.DownloadLink)
		fmt.Println("ModType:", mod.ModType)
		fmt.Println("Description:", mod.Description)
	}
}
