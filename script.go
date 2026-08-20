package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sync"

	"github.com/BurntSushi/toml"
	"golang.org/x/sync/errgroup"
)

//
//
//   CONSTS
//
//

const (
	BUILD_DIR   = "build"
	TEMP_DIR    = "build/temp"
	CONFIG_PATH = "mods.toml"
)

var (
	config   *Config
	configMu sync.RWMutex
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
	configMu.RLock()
	gameVersion := config.Meta.Version
	modLoader := config.Meta.Loader
	configMu.RUnlock()

	switch m.Source {
	case "modrinth":
		versions, err := getModrinthVersions(m.Name)
		if err != nil {
			return err
		}

		found := false
		for _, v := range versions {
			if v.Version == m.Version && slices.Contains(v.GameVersions, gameVersion) && slices.Contains(v.Loaders, modLoader) {
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
	// Создание директорий

	if err := os.MkdirAll(BUILD_DIR, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: Не удалось создать директорию сборки по пути %q\n", BUILD_DIR)
		os.Exit(1)
	}

	if err := os.MkdirAll(TEMP_DIR, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: Не удалось создать директорию временных файлов по пути %q\n", TEMP_DIR)
		os.Exit(1)
	}

	// Инициализация конфига

	cfg, err := loadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка загрузки конфига: %v\n", err)
		os.Exit(1)
	}

	configMu.Lock()
	config = cfg
	configMu.Unlock()

	// Инициализация модов

	fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ О СЕРВЕРНЫХ МОДАХ...")
	for i := range config.List.Server {
		if err := config.List.Server[i].initMissingFields(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("ИНФОРМАЦИЯ О СЕРВЕРНЫХ МОДАХ ПОЛУЧЕНА")

	fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ ОБ ОБЯЗАТЕЛЬНЫХ МОДАХ...")
	for i := range config.List.Required {
		if err := config.List.Required[i].initMissingFields(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("ИНФОРМАЦИЯ ОБ ОБЯЗАТЕЛЬНЫХ МОДАХ ПОЛУЧЕНА")

	fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ О РЕКОМЕНДУЕМЫХ МОДАХ...")
	for i := range config.List.Recommended {
		if err := config.List.Recommended[i].initMissingFields(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("ИНФОРМАЦИЯ О РЕКОМЕНДУЕМЫХ МОДАХ ПОЛУЧЕНА")

	fmt.Println("ПОЛУЧЕНИЕ ИНФОРМАЦИИ ОБ ОПЦИОНАЛЬНЫХ МОДАХ...")
	for i := range config.List.Optional {
		if err := config.List.Optional[i].initMissingFields(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Println("ИНФОРМАЦИЯ ОБ ОПЦИОНАЛЬНЫХ МОДАХ ПОЛУЧЕНА")

	// Очистка временной директории

	entries, err := os.ReadDir(TEMP_DIR)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка: Не удалось очистить директорию временных файлов")
		os.Exit(1)
	}
	for _, entry := range entries {
		path := filepath.Join(TEMP_DIR, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: Не удалось очистить директорию временных файлов")
			os.Exit(1)
		}
	}
}
