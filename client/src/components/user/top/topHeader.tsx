import Image from "next/image";
import SearchFormTop from "./searchForm";
import { Category } from "@/types/category";

type Props = {
  categories: Category[];
};

const TopHeader: React.FC<Props> = ({ categories }) => {
  console.log(categories);
  console.log("ここ通ってる？")
  return (
    <div className="bg-green-600 text-white h-80 z-40">
      <nav className="w-full h-16 flex justify-between items-center py-2 px-4">
        <a href="/" className="flex items-end">
          <Image
            src="/monta-no-mori-logo.svg"
            alt="もんたの森のロゴ"
            height={110} // 必須項目なのでとりあえず設定してるだけ
            width={110}
            style={{ height: "auto", objectFit: "contain" }}
          />
        </a>

        <div className="cursor-pointer w-12 h-12 rounded-full flex flex-col items-center justify-center hover:bg-white hover:bg-opacity-20 duration-200">
          <span className="w-8 h-1 bg-white block rounded-full mb-2"></span>
          <span className="w-8 h-1 bg-white block rounded-full mb-2"></span>
          <span className="w-8 h-1 bg-white block rounded-full"></span>
        </div>
      </nav>

      <div className="w-full h-64">
        <div className="h-full flex items-center justify-center flex-col">
          <SearchFormTop categories={categories}/>
        </div>
      </div>
    </div>
  );
};

export default TopHeader;
