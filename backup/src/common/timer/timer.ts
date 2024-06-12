import { hrtime } from 'process';

export class Timer {
  #startTime: number | null = null;

  private constructor() {
    // Left empty intentionally
  }

  /**
   * Creates and starts a new timer
   * @returns  The timer instance
   */
  public static start(): Timer {
    const timer = new Timer();
    [, timer.#startTime] = hrtime();
    return timer;
  }

  /**
   * Gets the elapsed time in milliseconds
   */
  public get elapsed(): number {
    const [, now] = hrtime();
    if (this.#startTime == null) throw new Error('Timer not started');
    return Math.floor((now - this.#startTime) / 1e6);
  }
}
