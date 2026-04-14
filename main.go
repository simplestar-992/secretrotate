package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	storeFile = filepath.Join(os.Getenv("HOME"), ".config/secretrotate/store.enc")
	showOnly  bool
	addSecret, delSecret, rotSecret string
	listMode  bool
)

type Secret struct {
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Version   int       `json:"version"`
	RotatedAt time.Time `json:"rotated_at"`
}

type Store struct {
	Version   int            `json:"version"`
	Secrets   map[string]Secret `json:"secrets"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func main() {
	home := os.Getenv("HOME")
	flag.StringVar(&storeFile, "f", filepath.Join(home, ".config/secretrotate/store.enc"), "Store file")
	flag.BoolVar(&showOnly, "s", false, "Show secret value")
	flag.StringVar(&addSecret, "a", "", "Add secret: name=value")
	flag.StringVar(&delSecret, "d", "", "Delete secret by name")
	flag.StringVar(&rotSecret, "r", "", "Rotate secret by name")
	flag.BoolVar(&listMode, "l", false, "List all secrets")
	flag.Parse()

	if flag.NFlag() == 0 || listMode {
		listSecrets()
		return
	}

	if addSecret != "" {
		parts := strings.SplitN(addSecret, "=", 2)
		if len(parts) != 2 {
			fmt.Fprintln(os.Stderr, "❌ Use: -a NAME=value")
			os.Exit(1)
		}
		addSecretCmd(parts[0], parts[1])
	}

	if rotSecret != "" {
		rotateSecret(rotSecret)
	}

	if delSecret != "" {
		deleteSecret(delSecret)
	}
}

func getStore() *Store {
	data, err := os.ReadFile(storeFile)
	if err != nil {
		s := Store{Version: 1}
		s.Secrets = make(map[string]Secret)
		return &s
	}

	out, _ := decrypt(data)
	var s Store
	json.Unmarshal(out, &s)
	if s.Secrets == nil {
		s.Secrets = make(map[string]Secret)
	}
	return &s
}

func saveStore(s *Store) {
	s.UpdatedAt = time.Now()
	data, _ := json.Marshal(s)
	enc, _ := encrypt(data)
	os.MkdirAll(filepath.Dir(storeFile), 0755)
	os.WriteFile(storeFile, enc, 0600)
}

func addSecretCmd(name, value string) {
	store := getStore()
	store.Secrets[name] = Secret{Name: name, Value: value, Version: 1, RotatedAt: time.Now()}
	saveStore(store)
	fmt.Printf("✅ Added secret: %s\n", name)
}

func rotateSecret(name string) {
	store := getStore()
	if s, ok := store.Secrets[name]; ok {
		s.Version++
		s.RotatedAt = time.Now()
		store.Secrets[name] = s
		saveStore(store)
		fmt.Printf("✅ Rotated secret: %s (v%d)\n", name, s.Version)
	} else {
		fmt.Printf("❌ Secret not found: %s\n", name)
	}
}

func deleteSecret(name string) {
	store := getStore()
	delete(store.Secrets, name)
	saveStore(store)
	fmt.Printf("✅ Deleted secret: %s\n", name)
}

func listSecrets() {
	store := getStore()
	fmt.Printf("🔐 SecretRotate Store (%s)\n", storeFile)
	fmt.Println("=============================")
	if len(store.Secrets) == 0 {
		fmt.Println("No secrets stored")
		return
	}
	for _, s := range store.Secrets {
		age := time.Since(s.RotatedAt)
		fmt.Printf("  %s: v%d (%s ago)\n", s.Name, s.Version, age.Round(time.Second))
	}
}

func encrypt(data []byte) ([]byte, error) {
	key := getKey()
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt(data []byte) ([]byte, error) {
	key := getKey()
	block, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(block)
	nonce := data[:gcm.NonceSize()]
	return gcm.Open(nil, nonce, data[gcm.NonceSize():], nil)
}

func getKey() []byte {
	envKey := os.Getenv("SECRET_KEY")
	if len(envKey) == 32 {
		return []byte(envKey)
	}
	home := os.Getenv("HOME")
	if home == "" {
		home = "/root"
	}
	keySrc := home + "secretrotate-default"
	key := make([]byte, 32)
	copy(key, []byte(keySrc))
	return key
}
