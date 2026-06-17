package core

import (
	"testing"
)

func TestDesignNameFromFileName(t *testing.T) {
	tests := []struct {
		name         string
		fileName     string
		expectedName string
	}{
		{
			name:         "given regular yaml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yaml",
			expectedName: "my-deployment",
		},
		{
			name:         "given regular yml file when DesignNameFromFileName then return name",
			fileName:     "my-deployment.yml",
			expectedName: "my-deployment",
		},
		{
			name:         "given tar.gz file when DesignNameFromFileName then strip compound extension",
			fileName:     "my-chart.tar.gz",
			expectedName: "my-chart",
		},
		{
			name:         "given json file when DesignNameFromFileName then return name",
			fileName:     "config.json",
			expectedName: "config",
		},
		{
			name:         "given filename with multiple dots when DesignNameFromFileName then strip only last extension",
			fileName:     "my.k8s.deployment.yaml",
			expectedName: "my.k8s.deployment",
		},
		{
			name:         "given empty filename when DesignNameFromFileName then return empty",
			fileName:     "",
			expectedName: "",
		},
		{
			name:         "given filename without extension when DesignNameFromFileName then return as is",
			fileName:     "mydesign",
			expectedName: "mydesign",
		},
		{
			name:         "given tgz file when DesignNameFromFileName then strip extension",
			fileName:     "helm-chart.tgz",
			expectedName: "helm-chart",
		},
		{
			name:         "given extension only when DesignNameFromFileName then return empty",
			fileName:     ".yaml",
			expectedName: "",
		},
		{
			name:         "given tar.gz extension only when DesignNameFromFileName then return empty",
			fileName:     ".tar.gz",
			expectedName: "",
		},
		{
			name:         "given unsupported xls extension when DesignNameFromFileName then still strip extension",
			fileName:     "spreadsheet.xls",
			expectedName: "spreadsheet",
		},
		{
			name:         "given unsupported zip extension when DesignNameFromFileName then still strip extension",
			fileName:     "archive.zip",
			expectedName: "archive",
		},
		{
			name:         "given unsupported exe extension when DesignNameFromFileName then still strip extension",
			fileName:     "installer.exe",
			expectedName: "installer",
		},
		{
			name:         "given unsupported tar extension when DesignNameFromFileName then strip last suffix only",
			fileName:     "archive.tar",
			expectedName: "archive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DesignNameFromFileName(tt.fileName)
			if result != tt.expectedName {
				t.Errorf("expected %q, got %q", tt.expectedName, result)
			}
		})
	}
}

func TestNoPrettifyDictionaryCollisions(t *testing.T) {
	seen := map[string]string{}
	for raw, display := range prettifyDictionary {
		if prev, ok := seen[display]; ok {
			t.Fatalf("collision: display %q maps from both raw %q and raw %q", display, prev, raw)
		}
		seen[display] = raw
	}
}

func TestDictionaryRoundTrip(t *testing.T) {
	raw := "apiVersion"
	m := map[string]interface{}{raw: "value"}
	p := prettifier(true)
	prettied := p.Prettify(m, true)
	displayKey := prettifyDictionary[raw]
	if _, ok := prettied[displayKey]; !ok {
		t.Fatalf("expected prettified key %q present, got %v", displayKey, prettied)
	}
	deprettied := p.DePrettify(prettied, true)
	if _, ok := deprettied[raw]; !ok {
		t.Fatalf("expected deprettified key %q present, got %v", raw, deprettied)
	}

	res := ConvertMapInterfaceMapString(raw, true, true)
	s, ok := res.(string)
	if !ok {
		t.Fatalf("expected string result from ConvertMapInterfaceMapString, got %T", res)
	}
	if s != displayKey {
		t.Fatalf("expected prettified string %q, got %q", displayKey, s)
	}

	res2 := ConvertMapInterfaceMapString(s, false, true)
	s2, ok := res2.(string)
	if !ok {
		t.Fatalf("expected string result from ConvertMapInterfaceMapString, got %T", res2)
	}
	if s2 != raw {
		t.Fatalf("expected deprettified string %q, got %q", raw, s2)
	}
}