"use client";

import { Illustration } from "@/types/illustration";
import { UpdatedAtFormat } from "@/utils/text";
import Image from "next/image";

interface Props {
  illustration: Illustration;
}

const IllustrationCard: React.FC<Props> = ({ illustration }) => {
  const isNew = (createdAtStr: string) => {
    const now: Date = new Date();
    const creationTime: Date = new Date(createdAtStr);

    const oneDay: number = 24 * 60 * 60 * 1000;
    const diff: number = Math.floor(now.getTime() - creationTime.getTime());
    return diff < oneDay;
  };

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
        {isNew(illustration.Image.created_at) && (
          <div className="absolute top-0 right-0 w-12 h-8 py-1 px-2 bg-green-600 text-white font-bold z-10">
            new
          </div>
        )}
        <div className="absolute inset-0 m-4">
          <Image
            className="md:group-hover:scale-110 duration-200 absolute top-0 left-0 w-full h-full object-cover"
            src={illustration.Image.original_src}
            alt={illustration.Image.title}
            fill
          />
        </div>
      </div>
      <div className="break-words md:group-hover:text-green-600 duration-200 mb-1">
        {illustration.Image.title}
      </div>
      <div className="text-gray-400 text-xs flex items-center justify-center">
        <Image
          className="mr-1"
          src="/icon/time.png"
          alt="timeアイコン"
          width={12}
          height={12}
        />
        <span>{UpdatedAtFormat(illustration.Image.updated_at)}</span>
      </div>
    </a>
  );
};

export default IllustrationCard;
