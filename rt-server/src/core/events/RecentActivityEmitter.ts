import { EventEmitterBase } from "./EventEmitterBase";
import type { Socket } from "socket.io";
import { TIME_INTERVAL } from "../interfaces/config/IConfig";
import type { IDataProvider } from "../interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "../interfaces/encryption/IEncryption";
import type { IEncryptionResult } from "../interfaces/encryption/IEncryptionResult";
import type { IRecentActivity } from "../interfaces/recentActivity/IRecentActivity";

/**
 * Responsible for emitting recent activity events.
 * Uses `RecentActivityDataProvider` to retrieve recent activities and then emits them through a specified channel.
 */
export class RecentActivityEmitter extends EventEmitterBase {
  private readonly dataProvider: IDataProvider<IRecentActivity>;
  private readonly encryptor: IEncryption;
  private readonly socket: Socket;

  /**
   * @param {IDataProvider<IRecentActivity>} dataProvider - The data provider used to retrieve recent activities.
   * @param {IEncryption} encryptor - The encryption utility used to encrypt activity data before transmission.
   * @param {Socket} socket - The Socket.IO socket through which the activity data is emitted.
   */
  constructor(
    dataProvider: IDataProvider<IRecentActivity>,
    encryptor: IEncryption,
    socket: Socket,
  ) {
    super();
    this.dataProvider = dataProvider;
    this.encryptor = encryptor;
    this.socket = socket;
  }

  /**
   * Emits recent activity events at regular intervals.
   */
  public emitEvent(): void {
    const items: IRecentActivity[] = this.dataProvider.getData();

    const handler: TimerHandler = (): void => {
      const itemToEmit: IRecentActivity =
        items[Math.floor(Math.random() * items.length)];
      const payload: IEncryptionResult = this.encryptor.encrypt(
        JSON.stringify(itemToEmit),
      );
      this.socket.emit("newActivity", payload);

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
