package model

// Theme represents presentation styling.
type Theme struct {
	Name       string            `json:"name,omitempty"`
	Primary    string            `json:"primary,omitempty"`    // Primary color (hex)
	Secondary  string            `json:"secondary,omitempty"`  // Secondary color (hex)
	Background string            `json:"background,omitempty"` // Background color (hex)
	Font       string            `json:"font,omitempty"`       // Font family
	Custom     map[string]string `json:"custom,omitempty"`     // Backend-specific settings
}

// DefaultTheme returns a default theme.
func DefaultTheme() *Theme {
	return &Theme{
		Name:       "default",
		Primary:    "#2196F3",
		Secondary:  "#FFC107",
		Background: "#FFFFFF",
		Font:       "sans-serif",
	}
}

// DarkTheme returns a dark theme.
func DarkTheme() *Theme {
	return &Theme{
		Name:       "dark",
		Primary:    "#90CAF9",
		Secondary:  "#FFE082",
		Background: "#121212",
		Font:       "sans-serif",
	}
}

// GetCustom returns a custom theme value with a default fallback.
func (t *Theme) GetCustom(key, defaultValue string) string {
	if t == nil || t.Custom == nil {
		return defaultValue
	}
	if v, ok := t.Custom[key]; ok {
		return v
	}
	return defaultValue
}

// SetCustom sets a custom theme value.
func (t *Theme) SetCustom(key, value string) {
	if t.Custom == nil {
		t.Custom = make(map[string]string)
	}
	t.Custom[key] = value
}
