import { NextRequest, NextResponse } from "next/server";
import { VerifyAPI } from "@/api/auth";

export async function middleware(request: NextRequest) {
  const accessToken = request.cookies.get("access_token")?.value;

  console.log("そもそもこことおってる？");
  console.log(accessToken);

  if (!accessToken) {
    console.log("access tokenがない場合はここ");
    if (request.nextUrl.pathname === "/login") {
      return NextResponse.next();
    } else {
      return NextResponse.redirect(new URL("/login", request.url));
    }
  }

  const formData = new URLSearchParams();
  formData.append("access_token", accessToken);

  console.log(formData);

  try {
    const response = await fetch(VerifyAPI(), {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body: formData.toString(),
    });

    const data = await response.json();
    console.log(data);
    if (data.result) {
      console.log("ここまできてたらうまくいく");
      return NextResponse.next();
    } else {
      console.log("verifyAPIのエラー");
      return NextResponse.redirect(new URL("/login", request.url));
    }
  } catch (error) {
    console.log("verify失敗");
    console.error(error);
    return NextResponse.redirect(new URL("/login", request.url));
  }
}

export const config = {
  matcher: ["/admin/:path*", "/login"],
};
