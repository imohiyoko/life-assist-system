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
)

// voice-analysis-serviceのコンテナ内アドレス
const voiceAnalysisServiceURL = "http://voice-analysis-service:8000/v1/analyze"

// App構造体でリポジトリを保持
type App struct {
	repo *repository.Repository
}

// 音声分析とDB保存を行うハンドラ
func (a *App) analyzeSpeechHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for speech analysis")

	// 1. クライアントからのファイルを取得
	file, header, err := r.FormFile("audio_file")
	if err != nil {
		http.Error(w, "Could not get audio file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. voice-analysis-serviceへ転送するための新しいリクエストボディを作成
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("audio_file", header.Filename)
	if err != nil {
		http.Error(w, "Failed to create form file for forwarding", http.StatusInternalServerError)
		return
	}
	if _, err = io.Copy(part, file); err != nil {
		http.Error(w, "Failed to copy file content", http.StatusInternalServerError)
		return
	}
	writer.Close()

	// 3. voice-analysis-serviceへの新しいPOSTリクエストを作成
	req, err := http.NewRequest("POST", voiceAnalysisServiceURL, body)
	if err != nil {
		http.Error(w, "Failed to create new request to analysis service", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// 4. リクエストを実行し、応答を取得
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to call analysis service", http.StatusInternalServerError)
		log.Printf("Error calling analysis service: %v", err)
		return
	}
	defer resp.Body.Close()

	// 5. 応答ボディを一度だけ読み込む
	analysisRespBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read analysis service response", http.StatusInternalServerError)
		return
	}

	// 6. 応答をDBに保存
	// JSONを文字列としてそのまま'details'カラムに保存する
	if err := a.repo.SaveEvent("SPEECH_ANALYSIS_COMPLETED", string(analysisRespBody)); err != nil {
		log.Printf("Failed to save event to DB: %v", err)
		// ここではエラーを返さず、ログに残すだけにする
	} else {
		log.Println("✅ Successfully saved analysis result to DB!")
	}

	// 7. 元の応答をそのままクライアントに返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(analysisRespBody)
}

func main() {
	// 環境変数からデータベース接続URLを取得
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// データベースに接続
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()
	log.Println("Successfully connected to the database!")
	
	repo := repository.NewRepository(dbpool)
	
	// テーブルを初期化
	if err := repo.InitTable(); err != nil {
		log.Fatalf("Failed to initialize table: %v\n", err)
	}
	
	app := &App{repo: repo}

	// 新しいエンドポイント `/v1/analyze-speech` を登録
	http.HandleFunc("/v1/analyze-speech", app.analyzeSpeechHandler)

	log.Println("Starting Go API server on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
