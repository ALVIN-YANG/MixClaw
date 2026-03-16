package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

//go:embed web/*
var webFS embed.FS

// Config represents the agent configuration
type Config struct {
	Provider  string `json:"provider"`
	Endpoint  string `json:"endpoint"`
	Model     string `json:"model"`
	APIKey    string `json:"apiKey,omitempty"`
	Workspace string `json:"workspace"`
	Restrict  bool   `json:"restrict"`
}

var (
	config     Config
	configLock sync.RWMutex
	configPath string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	configDir := filepath.Join(home, ".mixclaw")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Warning: failed to create config directory: %v", err)
	}
	configPath = filepath.Join(configDir, "config.json")
	loadConfig()
}

func loadConfig() {
	configLock.Lock()
	defer configLock.Unlock()

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Set defaults if config doesn't exist
		config = Config{
			Provider:  "volcengine",
			Workspace: "./workspace",
			Restrict:  true,
		}
		return
	}
	if err := json.Unmarshal(data, &config); err != nil {
		log.Printf("Warning: failed to parse config file: %v", err)
	}
}

func saveConfig() error {
	configLock.RLock()
	data, err := json.MarshalIndent(config, "", "  ")
	configLock.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func main() {
	mux := http.NewServeManager()
	
	// Serve embedded web UI
	staticFS, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API Endpoints
	mux.HandleFunc("/api/health", handleHealth)
	mux.HandleFunc("/api/config", handleConfig)
	mux.HandleFunc("/api/chat", handleChat)

	port := "18790"
	fmt.Printf("🦞 MixClaw Gateway starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func http.NewServeManager() *http.ServeMux {
	return http.NewServeMux()
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	configLock.RLock()
	defer configLock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"model":  config.Model,
	})
}

func handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		configLock.RLock()
		defer configLock.RUnlock()

		cfgCopy := config
		if cfgCopy.APIKey != "" {
			cfgCopy.APIKey = "********" // Mask API key for GET requests
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfgCopy)
		return
	}

	if r.Method == http.MethodPost {
		var newConfig Config
		if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		configLock.Lock()
		// Only update API key if it's not the masked string and not empty
		if newConfig.APIKey != "" && newConfig.APIKey != "********" {
			config.APIKey = newConfig.APIKey
		}
		config.Provider = newConfig.Provider
		config.Endpoint = newConfig.Endpoint
		config.Model = newConfig.Model
		config.Workspace = newConfig.Workspace
		config.Restrict = newConfig.Restrict
		configLock.Unlock()

		if err := saveConfig(); err != nil {
			http.Error(w, "Failed to save config", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Message   string `json:"message"`
		SessionID string `json:"sessionId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	configLock.RLock()
	model := config.Model
	configLock.RUnlock()

	if model == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "请先在设置页面配置大语言模型",
		})
		return
	}

	// TODO: Integrate actual LLM call here using the configured provider
	// For MVP Phase 1, we return a mock response
	mockResponse := fmt.Sprintf("🦞 (MixClaw 后端已收到) 这是针对 '%s' 的回复。目前我还在开发阶段，大模型接入层即将完成！当前使用的模型是: %s", req.Message, model)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"response":  mockResponse,
		"sessionId": req.SessionID,
	})
}
