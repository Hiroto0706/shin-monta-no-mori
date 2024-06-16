import { Illustration } from "@/types/illustration";
import Image from "next/image";
import React from "react";

interface Props {
  illustrations: Illustration[];
}

const ListIllustrations: React.FC<Props> = ({ illustrations }) => {
  return (
    <>
      <div className="mt-8">
        <div className="grid grid-cols-2 gap-x-4 gap-y-8 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 2xl:grid-cols-6">
          {illustrations.map((illustration) => (
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
              <span className="group-hover:text-green-600 group-hover:font-bold duration-200">
                {illustration.Image.title}
              </span>
            </a>
          ))}
        </div>
      </div>
    </>
  );
};

export default ListIllustrations;
