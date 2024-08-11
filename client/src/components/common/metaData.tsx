import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "もんたの森｜ゆるーくてかわいい無料イラスト",
  description:
    "もんたの森は『もんた』が書いたイラストたちを無料で保存、コピーして使用することができる無料イラストサイトです。ゆるくてかわいいもの大好きな『もんた』が趣味で描いてる絵を暇な時にアップロードしています！",
  icons: {
    icon: [
      {
        url: "/src/app/favicon.ico",
        href: "/src/app/favicon.ico",
      },
    ],
  },
  openGraph: {
    images: [
      {
        url: "https://storage.googleapis.com/shin-monta-no-mori/montanomori-top-image.png",
        alt: "もんたの森のイメージ画像",
      },
    ],
  },
};
