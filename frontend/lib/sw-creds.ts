const DB_NAME = "e2clicker-sw";
const STORE_NAME = "creds";
const CREDS_KEY = "auth";

export type SWCredentials = {
  token: string;
  deviceID: string;
};

function openCredsDB(): Promise<IDBDatabase> {
  return new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, 1);
    req.onupgradeneeded = () => {
      if (!req.result.objectStoreNames.contains(STORE_NAME)) {
        req.result.createObjectStore(STORE_NAME);
      }
    };
    req.onsuccess = () => resolve(req.result);
    req.onerror = () => reject(req.error);
  });
}

export async function saveSWCredentials(creds: SWCredentials): Promise<void> {
  if (typeof indexedDB === "undefined") return;
  const db = await openCredsDB();
  try {
    await new Promise<void>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, "readwrite");
      tx.objectStore(STORE_NAME).put(creds, CREDS_KEY);
      tx.oncomplete = () => resolve();
      tx.onerror = () => reject(tx.error);
    });
  } finally {
    db.close();
  }
}

export async function loadSWCredentials(): Promise<SWCredentials | null> {
  if (typeof indexedDB === "undefined") return null;
  const db = await openCredsDB();
  try {
    return await new Promise<SWCredentials | null>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, "readonly");
      const req = tx.objectStore(STORE_NAME).get(CREDS_KEY);
      req.onsuccess = () => resolve((req.result as SWCredentials | undefined) ?? null);
      req.onerror = () => reject(req.error);
    });
  } finally {
    db.close();
  }
}

export async function clearSWCredentials(): Promise<void> {
  if (typeof indexedDB === "undefined") return;
  const db = await openCredsDB();
  try {
    await new Promise<void>((resolve, reject) => {
      const tx = db.transaction(STORE_NAME, "readwrite");
      tx.objectStore(STORE_NAME).delete(CREDS_KEY);
      tx.oncomplete = () => resolve();
      tx.onerror = () => reject(tx.error);
    });
  } finally {
    db.close();
  }
}
