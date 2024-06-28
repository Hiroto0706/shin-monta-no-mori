import { NextRequest, NextResponse } from "next/server";
import { VerifyAPI } from "@/api/auth";

export async function middleware(request: NextRequest) {
  const accessToken = request.cookies.get("access_token")?.value;

  console.log("そもそもこことおってる？");
  console.log(accessToken);

  if (!accessToken) {
    console.log("access tokenがない場合はここ")
    return NextResponse.redirect(new URL("/login", request.url));
  }

  const formData = new URLSearchParams();
  formData.append("access_token", accessToken);

  try {
    const response = await fetch(VerifyAPI(), {
      method: "POST",
      body: formData,
    });

    const data = await response.json();
    if (data.result) {
      console.log("ここまできてたらうまくいく")
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
