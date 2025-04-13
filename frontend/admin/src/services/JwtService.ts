const ACCESS_TOKEN_KEY = 'access_token' as string
const REFRESH_TOKEN_KEY = 'refresh_token' as string

class JwtService {
  public static getAccessToken(): string | null {
    return window.localStorage.getItem(ACCESS_TOKEN_KEY)
  }

  public static getRefreshToken(): string | null {
    return window.localStorage.getItem(REFRESH_TOKEN_KEY)
  }

  public static saveTokens(accessToken: string, refreshToken: string): void {
    window.localStorage.setItem(ACCESS_TOKEN_KEY, accessToken)
    window.localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
  }

  public static destroyTokens(): void {
    window.localStorage.removeItem(ACCESS_TOKEN_KEY)
    window.localStorage.removeItem(REFRESH_TOKEN_KEY)
  }
}

export default JwtService
