import { EventEmitterBase } from "./EventEmitterBase";
import type { Socket } from "socket.io";
import type { IDataProvider } from "../interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "../interfaces/encryption/IEncryption";
import type { IEncryptionResult } from "../interfaces/encryption/IEncryptionResult";
import type { IKeyMetrics } from "../interfaces/keyMetrics/IKeyMetrics";
import { TIME_INTERVAL } from "../interfaces/config/IConfig";

/**
 * Responsible for emitting key metrics events.
 * Uses `KeyMetricsDataProvider` to retrieve key metrics information and then emits it through a specified channel.
 */
export class KeyMetricsEmitter extends EventEmitterBase {
  private readonly dataProvider: IDataProvider<IKeyMetrics>;
  private readonly encryptor: IEncryption;
  private readonly socket: Socket;

  /**
   * @param {IDataProvider<IKeyMetrics>} dataProvider - The data provider used to retrieve key metrics.
   * @param {IEncryption} encryptor - The encryption utility used to encrypt metrics data before transmission.
   * @param {Socket} socket - The Socket.IO socket through which the metrics are emitted.
   */
  constructor(
    dataProvider: IDataProvider<IKeyMetrics>,
    encryptor: IEncryption,
    socket: Socket,
  ) {
    super();
    this.dataProvider = dataProvider;
    this.encryptor = encryptor;
    this.socket = socket;
  }

  /**
   * Emits metrics events at regular intervals.
   */
  public emitEvent(): void {
    const handler: TimerHandler = (): void => {
      const itemToEmit: IKeyMetrics[] = this.dataProvider.getData();
      const payload: IEncryptionResult = this.encryptor.encrypt(
        JSON.stringify(itemToEmit),
      );
      this.socket.emit("newKeyMetrics", payload);

      console.log(
        `Sent activity to ${this.socket.id}: ${JSON.stringify(itemToEmit)}`,
      );
    };

    handler();

    const intervalId: number = setInterval(handler, TIME_INTERVAL);
    this.socket.on("disconnect", () => {
      clearInterval(intervalId);
      console.log(`A user disconnected: ${this.socket.id}`);
    });
  }
}
