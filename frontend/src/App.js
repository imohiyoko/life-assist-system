import React, { useState, useEffect, useRef } from 'react';
import Keycloak from 'keycloak-js';

// Keycloakクライアントの初期化設定
const keycloak = new Keycloak({
  url: 'http://localhost:9080',    // docker-composeで設定したKeycloakのポート
  realm: 'lifesystem-realm',    // Keycloakで作成するレルム名
  clientId: 'frontend-client',  // Keycloakで作成するクライアントID
});

function App() {
  const [authenticated, setAuthenticated] = useState(false);
  const [userInfo, setUserInfo] = useState(null);
  
  // --- この部分を追加 ---
  // 初期化処理が一度だけ実行されるようにするためのフラグ
  const isRun = useRef(false);

  useEffect(() => {
    // 既に実行済みの場合は何もしない
    if (isRun.current) return;
    // 実行済みフラグを立てる
    isRun.current = true;

    console.log("Keycloakの初期化を開始します...");
    
    // Keycloakを初期化し、認証状態を確認
    keycloak.init({ onLoad: 'check-sso' })
      .then(auth => {
        // 初期化成功時の処理
        console.log("Keycloakの初期化完了。認証状態:", auth);
        setAuthenticated(auth);
        if (auth) {
          // 認証されていれば、ユーザー情報を取得
          keycloak.loadUserInfo().then(info => setUserInfo(info));
        }
      })
      .catch(error => {
        console.error("Keycloakの初期化に失敗しました:", error);
        alert("Keycloakの初期化に失敗しました。Keycloakサーバーが正しく起動し、設定されているか確認してください。ブラウザのコンソールで詳細なエラーを確認できます。");
      });
  }, []);

  // ログイン処理
  const login = () => {
    keycloak.login();
  };

  // ログアウト処理
  const logout = () => {
    keycloak.logout();
  };
  
  // Goサーバーの保護されたAPIを叩くテスト関数
  const callApi = () => {
    console.log('APIを呼び出します...');
    console.log('認証トークン:', keycloak.token);
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'sans-serif' }}>
      <header>
        <h1>生活支援システム</h1>
        {authenticated ? (
          <div>
            <p>ようこそ、<strong>{userInfo?.preferred_username || 'ユーザー'}</strong>さん！</p>
            <p>あなたのロール: {JSON.stringify(keycloak.tokenParsed?.realm_access?.roles)}</p>
            <button onClick={callApi}>API呼び出しテスト</button>
            <button onClick={logout} style={{ marginLeft: '10px' }}>ログアウト</button>
          </div>
        ) : (
          <div>
            <p>ログインしていません。</p>
            <button onClick={login}>ログイン</button>
          </div>
        )}
      </header>
    </div>
  );
}

export default App;
