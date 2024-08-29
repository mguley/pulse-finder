import type { Tunnel } from "localtunnel";
import localtunnel from "localtunnel";

/**
 * Responsible for managing a local tunnel that exposes the server publicly.
 */
export class TunnelManager {
  private tunnel?: Tunnel;

  /**
   * @param {number} port - The port on which the local server is running and which should be exposed publicly.
   * @param {string} subdomain - The desired subdomain for the public tunnel.
   */
  constructor(
    private readonly port: number,
    private readonly subdomain: string,
  ) {}

  /**
   * Sets up a local tunnel to expose the server publicly.
   *
   * @returns {Promise<void>} A promise that resolves when the tunnel is successfully set up.
   */
  public async setupTunnel(): Promise<void> {
    try {
      this.tunnel = await localtunnel({
        port: this.port,
        subdomain: this.subdomain,
      });
      console.log(`Server is publicly accessible via ${this.tunnel.url}`);

      this.tunnel.on("close", (): void => {
        console.log(`Tunnel closed.`);
      });
    } catch (e) {
      console.log(`Error setting up tunnel: ${e}`);
    }
  }

  /**
   * Closes the local tunnel, making the server no longer accessible publicly.
   *
   * @returns {Promise<void>} A promise that resolves when the tunnel is successfully closed.
   */
  public async closeTunnel(): Promise<void> {
    if (this.tunnel) {
      this.tunnel.close();
    }
  }
}
