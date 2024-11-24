export * from "./openapi.gen";

// secretQRRegex is a regular expression that matches the secret QR code format.
export const secretQRRegex = /^e2clicker:secret-v1:(.*)$/;

// secretQRData returns the secret QR code data for a given secret.
export function secretQRData(secret: string): string {
  return `e2clicker:secret-v1:${secret}`;
}

// isStatus returns true if err is an API error with the given status.
export function isStatus(err: any, ...status: number[]) {
  return "status" in err && status.includes(err.status);
}
