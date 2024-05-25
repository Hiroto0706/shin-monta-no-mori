"use client";

import Image from "next/image";

export default function Header() {
  return (
    <header className="bg-green-600 text-white h-16 flex items-center">
      <nav className="w-full ml-4 mr-8 flex justify-between">
        <a href="/" className="flex items-end">
          <Image
            src="/monta-no-mori-logo.svg"
            alt="もんたの森のロゴ"
            height={110} // 必須項目なのでとりあえず設定してるだけ
            width={110}
            style={{ height: "auto", objectFit: "contain" }}
          />
          <div className="ml-2 relative w-20 h-6">
            <span className="absolute -bottom-1.5">for ADMIN</span>
          </div>
        </a>
        <ul className="flex space-x-8">
          <li>
            <a href="/">Home</a>
          </li>
          <li>
            <a href="/admin">Admin</a>
          </li>
        </ul>

        <button>button</button>
      </nav>
    </header>
  );
}
