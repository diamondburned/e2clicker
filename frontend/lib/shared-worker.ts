export type WorkerMessage = {
  type: "error";
  data: WorkerError;
};

export class WorkerError extends Error {
  constructor(message: string, options?: ErrorOptions) {
    super("internal worker error: " + message, options);
  }
}
