import Header from "@/components/user/common/handleHeader";
import UserSidebar from "@/components/user/common/sidebar/sidebar";
import Image from "next/image";

const NotFoundPage = () => {
  return (
    <html lang="jp" className="text-gray-600">
      <body>
        <Header />
        <UserSidebar />
        <div className="pl-0 md:pl-[calc(4rem+14rem)] pt-16 duration-200">
          <div className="p-4 md:p-12 flex justify-center items-center flex-col">
            <Image
              className="my-4"
              src="/top-sumasenn.png"
              alt="not foundイメージ"
              width={100}
              height={100}
            />
            <span className="flex justify-center text-lg mb-4">
              お探しのページは見つかりませんでした
            </span>
            <a
              href="/"
              className="flex justify-center underline border-blue-600 text-blue-600 cursor-pointer hover:text-blue-700 duration-200"
            >
              ホームに戻る
            </a>
          </div>
        </div>
      </body>
    </html>
  );
};

export default NotFoundPage;
