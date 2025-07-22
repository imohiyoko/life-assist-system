package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/DaikiTanaka-learner/life-assist-system/api-server/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/yaml.v3"
)

// config.ymlの構造に対応するGoの構造体
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Services struct {
		VoiceAnalysisURL string `yaml:"voice_analysis_url"`
	} `yaml:"services"`
	ServerPorts struct {
		APIServer int `yaml:"api_server"`
	} `yaml:"server_ports"`
}

// App構造体でリポジトリと設定を保持
type App struct {
	repo   *repository.Repository
	config *Config
}

func (a *App) analyzeSpeechHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("audio_file")
	if err != nil {
		http.Error(w, "Could not get audio file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio_file", header.Filename)
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", a.config.Services.VoiceAnalysisURL, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call analysis service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	analysisRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read analysis service response", http.StatusInternalServerError)
		return
	}

	if err := a.repo.SaveEvent("SPEECH_ANALYSIS_COMPLETED", string(analysisRespBody)); err != nil {
		log.Printf("Failed to save event to DB: %v", err)
	} else {
		log.Println("✅ Successfully saved analysis result to DB!")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(analysisRespBody)
}

func main() {
	// --- この部分を修正 ---
	// 1. 設定ファイルを絶対パスで読み込む
	configFile, err := os.ReadFile("/app/config.yml")
	if err != nil {
		log.Fatalf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(configFile, &cfg); err != nil {
		log.Fatalf("設定ファイルの解析に失敗しました: %v", err)
	}
	log.Println("✅ 設定ファイルを正常に読み込みました。")

	// 2. データベース接続URLを構築
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	// 3. データベースに接続
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("データベースへの接続に失敗しました: %v\n", err)
	}
	defer dbpool.Close()
	log.Println("✅ データベースに正常に接続しました！")

	repo := repository.NewRepository(dbpool)

	if err := repo.InitTable(); err != nil {
		log.Fatalf("テーブルの初期化に失敗しました: %v\n", err)
	}

	app := &App{
		repo:   repo,
		config: &cfg,
	}

	// 4. ハンドラを登録
	http.HandleFunc("/v1/analyze-speech", app.analyzeSpeechHandler)

	// 5. サーバーを起動
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPorts.APIServer)
	log.Printf("Go APIサーバーをポート %s で起動します...", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
