import axiosInstance from "../axios/axios";

/**
 * A singleton service to manage JSON Web Tokens (JWT).
 * Handles fetching, storing, refreshing and clearing tokens.
 */
class JwtService {
  private static instance: JwtService;
  private token: string | null = null;
  private tokenExpiry: Date | null = null;
  private refreshTimeout: ReturnType<typeof setTimeout> | null = null;

  /**
   * Private constructor to enforce singleton pattern.
   */
  private constructor() {}

  /**
   * Returns the single instance of the JwtService class.
   * If no instance exists, it creates one.
   *
   * @returns {JwtService} The singleton instance of JwtService.
   */
  public static getInstance(): JwtService {
    if (!JwtService.instance) {
      JwtService.instance = new JwtService();
    }
    return JwtService.instance;
  }

  /**
   * Fetches a new JWT token from the server and schedules its refresh.
   *
   * @throws {Error} If the token cannot be fetched.
   * @private
   */
  private async fetchToken(): Promise<void> {
    try {
      const { data } = await axiosInstance.get<{ token: string }>("v1/jwt");
      this.setToken(data.token);
    } catch (error) {
      console.error("Failed to fetch JWT token:", error);
      throw new Error("Could not fetch JWT token.");
    }
  }

  /**
   * Sets the JWT token, calculates its expiry, and schedules a refresh.
   *
   * @param {string} token - The JWT token string.
   * @private
   */
  private setToken(token: string): void {
    this.token = token;
    this.tokenExpiry = new Date(Date.now() + 24 * 60 * 60 * 1000); // Token valid for 24 hours
    this.scheduleTokenRefresh();
  }

  /**
   * Schedules the token refresh 1 hour before expiry.
   * Clears any existing refresh timeout before scheduling.
   *
   * @private
   */
  private scheduleTokenRefresh(): void {
    if (this.refreshTimeout) {
      clearTimeout(this.refreshTimeout);
    }

    if (this.tokenExpiry) {
      const timeUntilRefresh =
        this.tokenExpiry.getTime() - Date.now() - 60 * 60 * 1000; // Refresh 1 hour early
      if (timeUntilRefresh > 0) {
        this.refreshTimeout = setTimeout(
          () => this.fetchToken(),
          timeUntilRefresh,
        );
      }
    }
  }

  /**
   * Retrieves the current valid token.
   * If the token has expired or is not available, fetches a new one.
   *
   * @returns {Promise<string | null>} A promise that resolves to the current valid token or null.
   */
  public async getToken(): Promise<string | null> {
    if (this.token && this.tokenExpiry && this.tokenExpiry > new Date()) {
      return this.token;
    }
    await this.fetchToken();
    return this.token;
  }

  /**
   * Clears the current token and any scheduled refresh timeout.
   * Useful for logout scenarios or resetting the service.
   */
  public clearToken(): void {
    this.token = null;
    this.tokenExpiry = null;
    if (this.refreshTimeout) {
      clearTimeout(this.refreshTimeout);
      this.refreshTimeout = null;
    }
  }
}

export default JwtService;
