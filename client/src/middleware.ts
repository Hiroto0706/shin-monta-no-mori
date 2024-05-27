import { NextRequest, NextResponse } from "next/server";
import { Verify } from "@/api/auth";

export async function middleware(request: NextRequest) {
  const accessToken = request.cookies.get("access_token")?.value;

  if (!accessToken) {
    if (request.nextUrl.pathname === "/login") {
      return NextResponse.next();
    } else {
      return NextResponse.redirect(new URL("/login", request.url));
    }
  }

  const formData = new URLSearchParams();
  formData.append("access_token", accessToken);

  try {
    const response = await fetch(Verify(), {
      method: "POST",
      body: formData,
    });

    const data = await response.json();
    if (data.result) {
      if (request.nextUrl.pathname === "/login") {
        return NextResponse.redirect(new URL("/admin", request.url));
      }
      return NextResponse.next();
    } else {
      return NextResponse.redirect(new URL("/login", request.url));
    }
  } catch (error) {
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

export const config = {
  matcher: ["/admin/:path*", "/login"],
};
