import { EventEmitterBase } from "./EventEmitterBase";
import type { Socket } from "socket.io";
import { TIME_INTERVAL } from "../interfaces/config/IConfig";
import type { IDataProvider } from "../interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "../interfaces/encryption/IEncryption";
import type { IEncryptionResult } from "../interfaces/encryption/IEncryptionResult";
import type { IJobStatus } from "../interfaces/jobStatus/IJobStatus";

/**
 * Responsible for emitting job statuses events.
 * It uses `JobStatusDataProvider` to retrieve recent statistics and then emits them through a specified channel.
 */
export class JobStatusEmitter extends EventEmitterBase {
  /**
   * @param {IDataProvider<IJobStatus>} dataProvider - The data provider used to retrieve recent statistics.
   * @param {IEncryption} encryptor - The encryption utility used to encrypt activity data before transmission.
   * @param {Socket} socket - The Socket.IO socket through which the activity data is emitted.
   */
  constructor(
    private readonly dataProvider: IDataProvider<IJobStatus>,
    private readonly encryptor: IEncryption,
    private readonly socket: Socket,
  ) {
    super();
  }

  /**
   * Emits statistics events at regular intervals.
   */
  public emitEvent(): void {
    const handler: TimerHandler = (): void => {
      const itemToEmit: IJobStatus[] = this.dataProvider.getData();
      const payload: IEncryptionResult = this.encryptor.encrypt(
        JSON.stringify(itemToEmit),
      );
      this.socket.emit("newJobStatuses", payload);

      console.log(
        `Sent activity to ${this.socket.id}: ${JSON.stringify(itemToEmit)}`,
      );
    };

    handler();

    const intervalId: number = setInterval(handler, TIME_INTERVAL);
    this.socket.on("disconnect", (): void => {
      clearInterval(intervalId);
      console.log(`A user disconnected: ${this.socket.id}`);
    });
  }
}
