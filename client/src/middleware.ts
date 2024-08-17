import { NextRequest, NextResponse } from "next/server";
import { VerifyAPI } from "@/api/auth";

export async function middleware(request: NextRequest) {
  const accessToken = request.cookies.get("access_token")?.value;

  // /login ページへのアクセスで accessToken がない場合、そのまま進める
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
    const response = await fetch(VerifyAPI(), {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: formData.toString(),
    });

    const data = await response.json();
    if (data.result) {
      // /login にアクセスしていて、ログイン済みなら /admin にリダイレクト
      if (request.nextUrl.pathname === "/login") {
        return NextResponse.redirect(new URL("/admin", request.url));
      }
      return NextResponse.next();
    } else {
      // すでに /login にいる場合、リダイレクトを避ける
      if (request.nextUrl.pathname !== "/login") {
        return NextResponse.redirect(new URL("/login", request.url));
      }
    }
  } catch (error) {
    console.error(error);
    // エラー時にも、すでに /login にいる場合はリダイレクトを避ける
    if (request.nextUrl.pathname !== "/login") {
      return NextResponse.redirect(new URL("/login", request.url));
    }
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/admin/:path*", "/login"],
};
