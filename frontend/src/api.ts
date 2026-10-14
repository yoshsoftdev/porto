const API = (import.meta.env.VITE_API_URL || "http://localhost:8080").replace(/\/$/, "");

export function asset(path?: string | null) {
  if (!path) return "";
  if (path.startsWith("http")) return path;
  return API + path;
}

export async function getJSON<T>(path: string, token?: string): Promise<T> {
  const res = await fetch(API + path, {
    headers: token ? { Authorization: `Bearer ${token}` } : undefined,
  });
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

export async function sendJSON<T>(
  path: string,
  method: string,
  body?: unknown,
  token?: string,
): Promise<T> {
  const res = await fetch(API + path, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

export async function uploadImage(file: File, token: string) {
  const fd = new FormData();
  fd.append("file", file);
  const res = await fetch(API + "/api/admin/images", {
    method: "POST",
    headers: { Authorization: `Bearer ${token}` },
    body: fd,
  });
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json() as Promise<{ id: number; url: string }>;
}

async function errorMessage(res: Response) {
  try {
    const data = await res.json();
    return data.error || res.statusText;
  } catch {
    return res.statusText;
  }
}
