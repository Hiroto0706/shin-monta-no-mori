"use client";

import Image from "next/image";
import { Illustration } from "@/types/illustration";

interface IllustrationsProps {
  illustrations: Illustration[];
}

export default function Illustrations({ illustrations }: IllustrationsProps) {
  if (!illustrations || illustrations.length === 0) {
    return <div>No illustrations available</div>;
  }

  return (
    <div>
      <h1 className="text-2xl font-bold">イラスト一覧</h1>

      <div className="my-12">
        <ul className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
          {illustrations.map((illustration, index) => (
            <li key={index} className="p-2 border-2 border-gray-200 rounded-xl">
              <div className="flex items-center">
                <span className="ml-4 font-bold text-2xl">
                  {illustration.Image.title}
                </span>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
