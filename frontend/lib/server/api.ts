import type { Cookies } from "@sveltejs/kit";

export function isAuthorized(request: { cookies: Cookies }): boolean {
  return !!request.cookies.get("token");
}
