"use client";

import Image from "next/image";
import Link from "next/link";
import { useState } from "react";

const AbsolutePromotion = () => {
  const [isImageVisible, setImageVisible] = useState(true);

  return (
    <>
      {isImageVisible && (
        <div className="fixed bottom-4 w-full max-w-[400px] md:max-w-[500px] aspect-[4/1] z-50 left-1/2 transform -translate-x-1/2 md:right-4 md:left-auto md:transform-none">
          <div className="relative w-full h-full">
            <>
              <Link
                href="https://store.line.me/emojishop/product/66ece5ac61b93a07d864a1a5/ja"
                target="_blank"
              >
                <Image
                  src="/promotions/20240924_montanomori_emoji_ab.svg"
                  alt="もんたの森のLINE絵文字プロモーション"
                  fill
                  style={{ objectFit: "cover" }}
                />
              </Link>
            </>

            <button
              onClick={() => setImageVisible(false)}
              className="absolute -top-6 -right-4 m-2 text-white bg-black bg-opacity-50 rounded-full w-6 h-6 flex items-center justify-center z-10"
            >
              ×
            </button>
          </div>
        </div>
      )}
    </>
  );
};

export default AbsolutePromotion;
