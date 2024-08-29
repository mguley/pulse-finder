import type {
  Server as HTTPServer,
  IncomingMessage,
  ServerResponse,
} from "http";
import { createServer } from "http";

/**
 * Responsible for managing the lifecycle of an HTTP server.
 */
export class HTTPServerManager {
  private readonly httpServer: HTTPServer;

  constructor() {
    this.httpServer = this.createHTTPServer();
  }

  /**
   * Starts the HTTP server.
   */
  public start(): void {
    this.httpServer.listen(Number(process.env.PORT) || 4000, (): void => {
      console.log(
        `HTTP Server is running on port ${Number(process.env.PORT) || 4000}`,
      );
    });
  }

  /**
   * Stops the HTTP server.
   */
  public stop(): void {
    this.httpServer.close((): void => {
      console.log(`HTTP Server is closed`);
    });
  }

  /**
   * Provides access to the underlying HTTP server instance.
   *
   * @returns {HTTPServer} The HTTP server instance.
   */
  public getServer(): HTTPServer {
    return this.httpServer;
  }

  /**
   * Creates and configures the HTTP server.
   *
   * @returns {HTTPServer} The configured HTTP server instance.
   */
  private createHTTPServer(): HTTPServer {
    const requestListener = (
      request: IncomingMessage,
      response: ServerResponse,
    ) => {
      if (request.method === "OPTIONS") {
        this.handlePreflightRequest(response);
      }
    };
    return createServer(requestListener);
  }

  /**
   * Handles CORS preflight requests.
   *
   * @param {ServerResponse} response - The server response object used to send the CORS headers.
   */
  private handlePreflightRequest(response: ServerResponse): void {
    response.writeHead(200, {
      "Access-Control-Allow-Origin": "*",
      "Access-Control-Allow-Methods": "GET, POST, OPTIONS",
      "Access-Control-Allow-Headers": "Content-Type, Authorization",
    });
    response.end();
  }
}
