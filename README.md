# Life-Assist System (生活支援システム)
日々の生活から学習し、あなたの家庭だけの「執事」を育てるオープンソース・生活支援システム。

## プロジェクトの目的
生活上の会話内容や映像情報、ユーザーのスマートフォンの位置情報、レシートの画像情報、メールに届いた支払情報などから家庭内の情報を解析します。
解析した情報をもとにスケジュール管理や買い物リストの作成、食の好みの解析などを行い、生活の最適化提案のための通知を行うことを目指します。

## 主な機能 (v1.0目標)
音声によるスケジュール登録: ユーザーの家庭内での会話内容から予定やタスクを推察し、カレンダーに登録。


レシートOCRと家計簿: レシートを撮影するだけで、家計簿入力を半自動化。

シンプルなWebダッシュボード: 予定、タスク、支出を一覧できる閲覧専用の画面。

## アーキテクチャ
システムは、データ収集を行うクライアントと、分析処理を行うサーバー群からなるマイクロサービスアーキテクチャを採用しています。

## 運用・開発を行う前に必要なこと
このシステムをローカル環境で起動するためには、いくつかの事前準備が必要です。

### 1. Hugging Faceの準備
本システムの音声分析サービスは、pyannote.audioの事前学習済みモデルを利用します。

    1. アカウント作成: Hugging Faceでアカウントを作成・ログインしてください。

    2. 利用規約への同意: 以下の2つのモデルページにアクセスし、それぞれの利用規約をよく読み、同意してください。

    3. [pyannote/speaker-diarization-3.1][def-1]

    4. [pyannote/segmentation-3.0][def-2]

    5. アクセストークンの発行: [アクセストークンのページ][def-3] にアクセスし、read権限を持つ新しいトークンを発行します。

### 2. 設定ファイルの準備
    1.  プロジェクトのルートにある `config.yml.example` をコピーし、`config.yml` という名前で新しいファイルを作成します。
    2.  作成した `config.yml` を開き、`models.hugging_face_token` の値を、あなたがHugging Faceで発行したアクセストークンに書き換えてください。
    3.  GPUを利用したい場合は、`compute.device` の値を `cpu` から `gpu` に変更してください。


## 技術スタック
バックエンドAPIサーバー: Go

AI / MLサービス: Python (FastAPI, Whisper, PyTorch, pyannote.audio)

フロントエンド / ダッシュボード: React

データベース: PostgreSQL

認証基盤: Keycloak

デプロイ環境: Docker / Docker Compose

## コミュニティ
開発に関する議論や雑談は、公式Discordサーバーで行っています。お気軽にご参加ください！

Discord: https://discord.gg/Py7Nx38tNh

## 貢献方法
このプロジェクトはまだ初期段階であり、あらゆる形の貢献を歓迎します！バグ報告や機能提案など、気軽にIssueを作成してください。
(ここにCONTRIBUTING.mdへのリンクを将来的に追加)

## ライセンス
このプロジェクトはMITライセンスの下で公開されています。詳細はLICENSEファイルをご覧ください。

[def-2]: https://huggingface.co/pyannote/speaker-diarization-3.0
[def-1]: https://huggingface.co/pyannote/speaker-diarization-3.1
[def-3]: https://huggingface.co/settings/tokens