package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slices"

	"github.com/BurntSushi/toml"
	"golang.org/x/sync/errgroup"
)

//
//
//   CONSTS
//
//

const (
	BUILD_DIR    = "build"
	CONFIG_PATH  = "mods.toml"
	GAME_VERSION = "26.2"
	MOD_LOADER   = "fabric"
)

//
//
//   TYPES
//
//

type Config struct {
	Meta struct {
		Version string `toml:"version"`
		Loader  string `toml:"loader"`
	} `toml:"meta"`
	List struct {
		Server      []Mod `toml:"server"`
		Required    []Mod `toml:"required"`
		Recommended []Mod `toml:"recommended"`
		Optional    []Mod `toml:"optional"`
	} `toml:"list"`
}

func loadConfig() (*Config, error) {
	var config Config

	if _, err := toml.DecodeFile(CONFIG_PATH, &config); err != nil {
		return nil, fmt.Errorf("Ошибка при декодировании файла: %w", err)
	}

	return &config, nil
}

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

	fmt.Printf("Информация о моде %q получена\n", m.Name)

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

		found := false
		for _, v := range versions {
			if v.Version == m.Version && slices.Contains(v.GameVersions, GAME_VERSION) && slices.Contains(v.Loaders, MOD_LOADER) {
				found = true
				for _, f := range v.Files {
					if f.Primary {
						m.DownloadLink = f.URL
						break
					}
				}
				break
			}
		}

		if !found {
			return fmt.Errorf("Версия %q для мода %q не найдена", m.Version, m.Name)
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
	Version      string   `json:"version_number"`
	GameVersions []string `json:"game_versions"`
	Loaders      []string `json:"loaders"`
	Files        []struct {
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
	fmt.Println("ЗАГРУЗКА КОНФИГА...")
	config, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при загрузке конфига: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("КОНФИГ ЗАГРУЖЕН")

	// fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ О СЕРВЕРНЫХ МОДАХ...")
	// for i := range config.List.Server {
	// 	mod := &config.List.Server[i]
	// 	if err := mod.initMissingFields(); err != nil {
	// 		fmt.Fprintf(os.Stderr, "Ошибка при получении информации о моде: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }
	// fmt.Println("ИНФОРМАЦИЯ О СЕРВЕРНЫХ МОДАХ ПОЛУЧЕНА")

	// fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ ОБ ОБЯЗАТЕЛЬНЫХ МОДАХ...")
	// for i := range config.List.Required {
	// 	mod := &config.List.Required[i]
	// 	if err := mod.initMissingFields(); err != nil {
	// 		fmt.Fprintf(os.Stderr, "Ошибка при получении информации о моде: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }
	// fmt.Println("ИНФОРМАЦИЯ ОБ ОБЯЗАТЕЛЬНЫХ МОДАХ ПОЛУЧЕНА")

	// fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ О РЕКОМЕНДУЕМЫХ МОДАХ...")
	// for i := range config.List.Recommended {
	// 	mod := &config.List.Recommended[i]
	// 	if err := mod.initMissingFields(); err != nil {
	// 		fmt.Fprintf(os.Stderr, "Ошибка при получении информации о моде %q: %v\n", mod.Name, err)
	// 		os.Exit(1)
	// 	}
	// }
	// fmt.Println("ИНФОРМАЦИЯ О РЕКОМЕНДУЕМЫХ МОДАХ ПОЛУЧЕНА")

	// fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ ОБ ОПЦИОНАЛЬНЫХ МОДАХ...")
	// for i := range config.List.Optional {
	// 	mod := &config.List.Optional[i]
	// 	if err := mod.initMissingFields(); err != nil {
	// 		fmt.Fprintf(os.Stderr, "Ошибка при получении информации о моде: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }
	// fmt.Println("ИНФОРМАЦИЯ ОБ ОПЦИОНАЛЬНЫХ МОДАХ ПОЛУЧЕНА")

	for _, mod := range config.List.Recommended {
		fmt.Println("Name:", mod.Name)
		fmt.Println("PrettyName:", mod.PrettyName)
		fmt.Println("Source:", mod.Source)
		fmt.Println("Version:", mod.Version)
		fmt.Println("DownloadLink:", mod.DownloadLink)
		fmt.Println("ModType:", mod.ModType)
		fmt.Println("Description:", mod.Description)
		fmt.Println()
	}
}
