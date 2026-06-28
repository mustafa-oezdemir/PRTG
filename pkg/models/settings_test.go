package models

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestLoadPluginSettingsUsesConfiguredValues(t *testing.T) {
	settings, err := LoadPluginSettings(backend.DataSourceInstanceSettings{
		JSONData: []byte(`{"path":"prtg.example.local","cacheTime":120,"timeZone":"UTC"}`),
		DecryptedSecureJSONData: map[string]string{
			"apiKey": "secret-token",
		},
	})
	if err != nil {
		t.Fatalf("LoadPluginSettings returned error: %v", err)
	}

	if settings.Path != "prtg.example.local" {
		t.Fatalf("expected path to be loaded, got %q", settings.Path)
	}
	if settings.CacheTime != 120*time.Nanosecond {
		t.Fatalf("expected raw cacheTime duration value, got %s", settings.CacheTime)
	}
	if settings.Timezone != "UTC" {
		t.Fatalf("expected configured timezone UTC, got %q", settings.Timezone)
	}
	if settings.Secrets == nil || settings.Secrets.ApiKey != "secret-token" {
		t.Fatalf("expected apiKey secret to be loaded, got %#v", settings.Secrets)
	}
}

func TestLoadPluginSettingsDefaultsTimezone(t *testing.T) {
	settings, err := LoadPluginSettings(backend.DataSourceInstanceSettings{
		JSONData: []byte(`{"path":"prtg.example.local"}`),
	})
	if err != nil {
		t.Fatalf("LoadPluginSettings returned error: %v", err)
	}

	if settings.Timezone != "Europe/Berlin" {
		t.Fatalf("expected default timezone Europe/Berlin, got %q", settings.Timezone)
	}
}

func TestLoadPluginSettingsFallsBackToUTCForInvalidTimezone(t *testing.T) {
	settings, err := LoadPluginSettings(backend.DataSourceInstanceSettings{
		JSONData: []byte(`{"timeZone":"Not/A_Timezone"}`),
	})
	if err != nil {
		t.Fatalf("LoadPluginSettings returned error: %v", err)
	}

	if settings.Timezone != "UTC" {
		t.Fatalf("expected invalid timezone to fall back to UTC, got %q", settings.Timezone)
	}
}

func TestLoadPluginSettingsReturnsJSONError(t *testing.T) {
	_, err := LoadPluginSettings(backend.DataSourceInstanceSettings{
		JSONData: []byte(`{"path":`),
	})
	if err == nil {
		t.Fatal("expected invalid JSON to return an error")
	}
}
