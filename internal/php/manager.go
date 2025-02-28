package php

type PHPVersion struct {
	Version string
	Path    string
	Active  bool
}

type Manager struct {
	versions []PHPVersion
}

func NewManager() *Manager {
	// mock data for now
	return &Manager{
		versions: []PHPVersion{
			{Version: "8.3.6", Path: "/usr/bin/php8.3", Active: true},
			{Version: "8.2.27", Path: "/usr/bin/php8.2", Active: false},
			{Version: "8.1.31", Path: "/usr/bin/php8.1", Active: false},
		},
	}
}

func (m *Manager) GetVersions() []PHPVersion {
	return m.versions
}
