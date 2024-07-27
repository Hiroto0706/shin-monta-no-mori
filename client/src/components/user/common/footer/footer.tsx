"use client";

import Image from "next/image";
import Link from "next/link";
import React from "react";

type Props = {
  sns: {
    name: string;
    src: string;
    link: string;
  }[];
};

const Footer: React.FC<Props> = ({ sns }) => {
  return (
    <>
      <footer>
        <div className="max-w-[1100px] m-auto px-4 md:px-12">
          <div className="border-t border-gray-200 py-12 flex flex-wrap">
            {sns.map((s, index) => (
              <Link
                key={index}
                href={s.link}
                target="_blank"
                className="flex items-center my-2 cursor-pointer block w-28 md:w-40 group"
              >
                <div className="w-6 md:w-8 h-6 md:h-8 relative">
                  <Image
                    className="absolute w-full h-full mx-1 group-hover:opacity-70"
                    alt={s.name}
                    src={s.src}
                    fill
                  />
                </div>
                <span className="text-sm md:text-lg border-b border-gray-600 group-hover:opacity-70 duration-200 ml-2">
                  {s.name}
                </span>
              </Link>
            ))}
          </div>

          <div className="py-12 border-t border-gray-200 flex text-sm md:text-md">
            <span>&copy;もんたの森 2024</span>
            <Link
              href="/terms-of-service"
              className="ml-4 md:ml-8 underline cursor-pointer"
            >
              ご利用規約
            </Link>
            <Link
              href="/privacy-policy"
              className="ml-4 md:ml-8 underline cursor-pointer"
            >
              プライバシーポリシー
            </Link>
          </div>
        </div>
      </footer>
    </>
  );
};

export default Footer;
