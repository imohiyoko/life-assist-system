# frontend/Dockerfile
FROM node:18-alpine

WORKDIR /app

# 依存関係をインストールするためにpackage.jsonとlockファイルを先にコピー
# これにより、コードを変更しても毎回npm installが走るのを防ぐ
COPY package*.json ./
RUN npm install

# アプリケーションのソースコードをコピー
COPY . .

# アプリケーションを起動
CMD ["npm", "start"]
