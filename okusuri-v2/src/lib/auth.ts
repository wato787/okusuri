import { apiClient } from './apiClient';

export interface User {
  id: string;
  name: string;
  email: string;
  emailVerified: boolean;
  image?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Session {
  id: string;
  token: string;
  expiresAt: string;
}

export interface AuthResponse {
  user: User;
  session: Session;
}

// Google OAuth認証開始
export const signInWithGoogle = async (): Promise<void> => {
  const data = await apiClient.get<{ url: string }>('/api/auth/google');

  if (data.url) {
    window.location.href = data.url;
  } else {
    throw new Error('OAuth URLの取得に失敗しました');
  }
};

// セッション情報取得
export const getSession = async (
  token?: string
): Promise<AuthResponse | null> => {
  try {
    const storedToken = token || getAuthToken();
    if (!storedToken) return null;

    // トークンが文字列であることを確認
    if (typeof storedToken !== 'string') {
      console.error('無効なトークン形式:', storedToken);
      localStorage.removeItem('authToken');
      return null;
    }

    // 一時的にトークンを設定
    const originalToken = getAuthToken();
    setAuthToken(storedToken);

    try {
      const response = await apiClient.get<AuthResponse>('/api/auth/session');
      return response;
    } catch (error) {
      if (error instanceof Error && error.message.includes('401')) {
        localStorage.removeItem('authToken');
        return null;
      }
      throw error;
    } finally {
      // 元のトークンを復元
      if (originalToken) {
        setAuthToken(originalToken);
      } else {
        localStorage.removeItem('authToken');
      }
    }
  } catch (error) {
    console.error('セッション取得エラー:', error);
    return null;
  }
};

// サインアウト
export const signOut = async (): Promise<void> => {
  try {
    const token = getAuthToken();
    if (token) {
      await apiClient.post('/api/auth/signout');
    }
  } catch (error) {
    console.error('サインアウトエラー:', error);
  } finally {
    localStorage.removeItem('authToken');
  }
};

// トークンを保存
export const setAuthToken = (token: string): void => {
  localStorage.setItem('authToken', token);
};

// トークンを取得
export const getAuthToken = (): string | null => {
  return localStorage.getItem('authToken');
};
