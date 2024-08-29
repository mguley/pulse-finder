import type { Socket } from "socket.io";
import type { IDataProvider } from "../interfaces/dataProvider/IDataProvider";
import type { IEncryption } from "../interfaces/encryption/IEncryption";
import type { IEventEmitter } from "../interfaces/event/IEventEmitter";
import type { IRecentActivity } from "../interfaces/recentActivity/IRecentActivity";
import type { IKeyMetrics } from "../interfaces/keyMetrics/IKeyMetrics";
import { RecentActivityEmitter } from "./RecentActivityEmitter";
import { KeyMetricsEmitter } from "./KeyMetricsEmitter";
import { JobStatusEmitter } from "./JobStatusEmitter";
import type { IJobStatus } from "../interfaces/jobStatus/IJobStatus";

/**
 * Manages the process of creating emitters and emitting events to connected WebSocket clients.
 */
export class EventEmitterManager {
  /**
   * @param {IDataProvider<IRecentActivity>} recentActivityDataProvider - The data provider for recent activity data.
   * @param {IDataProvider<IKeyMetrics>} keyMetricsDataProvider - The data provider for key metrics data.
   * @param {IDataProvider<IJobStatus>} jobStatusDataProvider - The data provider for job status data.
   * @param {IEncryption} encryptor - The encryption service used to secure data before emitting.
   * @param {Socket} socket - The WebSocket connection through which events are emitted.
   */
  constructor(
    private readonly recentActivityDataProvider: IDataProvider<IRecentActivity>,
    private readonly keyMetricsDataProvider: IDataProvider<IKeyMetrics>,
    private readonly jobStatusDataProvider: IDataProvider<IJobStatus>,
    private readonly encryptor: IEncryption,
    private readonly socket: Socket,
  ) {}

  /**
   * Creates and initializes event emitters for recent activity, key metrics, and job status data.
   */
  public emit(): void {
    try {
      const recentActivityEmitter: IEventEmitter = new RecentActivityEmitter(
        this.recentActivityDataProvider,
        this.encryptor,
        this.socket,
      );
      const keyMetricsEmitter: IEventEmitter = new KeyMetricsEmitter(
        this.keyMetricsDataProvider,
        this.encryptor,
        this.socket,
      );
      const jobStatusEmitter: IEventEmitter = new JobStatusEmitter(
        this.jobStatusDataProvider,
        this.encryptor,
        this.socket,
      );

      recentActivityEmitter.emitEvent();
      keyMetricsEmitter.emitEvent();
      jobStatusEmitter.emitEvent();
    } catch (error) {
      console.log(`Error initializing emitters: ${error}`);
    }
  }
}
