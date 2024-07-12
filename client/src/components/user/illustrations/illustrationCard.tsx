"use client";

import { Illustration } from "@/types/illustration";
import Image from "next/image";

interface Props {
  illustration: Illustration;
}

const IllustrationCard: React.FC<Props> = ({ illustration }) => {
  return (
    <a
      key={illustration.Image.id}
      href={`/illustrations/${illustration.Image.id}`}
      className="group cursor-pointer"
    >
      <div
        className="mb-2 border-2 border-gray-200 rounded-xl bg-white relative w-full overflow-hidden"
        style={{ paddingTop: "100%" }}
      >
        <div className="absolute inset-0 m-4">
          <Image
            className="group-hover:scale-110 duration-200 absolute top-0 left-0 w-full h-full object-cover"
            src={illustration.Image.original_src}
            alt={illustration.Image.title}
            fill
          />
        </div>
      </div>
      <span className="break-words group-hover:text-green-600 duration-200">
        {illustration.Image.title}
      </span>
    </a>
  );
};

export default IllustrationCard;
