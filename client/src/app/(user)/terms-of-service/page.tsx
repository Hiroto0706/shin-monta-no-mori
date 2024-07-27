const TermsOfService = () => {
  return (
    <>
      <div>
        <h1 className="text-3xl font-bold mb-6">ようこそ！利用規約へ</h1>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトを訪れてくれてありがとう！こちらを利用することで、あなたは以下の利用規約に同意したことになるので、さっと目を通してみてね。
        </p>

        <h2 className="text-2xl font-bold mb-4">画像のフリー使用について</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトの画像は、あなたのプロジェクトやビジネスにぴったり合うよう、個人使用でも商用使用でもOK！ただし、画像自体を再販売したり、配布することはご遠慮くださいね。
        </p>

        <h2 className="text-2xl font-bold mb-4">画像の著作権について</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          こちらで提供している画像は、作成者または私たちに著作権があります。使用する時は、「著作権所有者名」のクレジット表示を忘れずに。
        </p>

        <h2 className="text-2xl font-bold mb-4">利用のルール</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          法律に反する使い方や、他人の権利を傷つけるような使い方はダメだよ。みんなが気持ちよく使えるサイトにしましょう。
        </p>

        <h2 className="text-2xl font-bold mb-4">これはOK！</h2>
        <ul className="my-4 bg-blue-50 py-6 px-12 rounded-lg list-disc">
          <li>ブログやウェブサイトでの使用</li>
          <li>商用プロジェクトでの使用</li>
          <li>プレゼンテーションやレポートでの使用</li>
        </ul>

        <h2 className="text-2xl font-bold mb-4">これはNG！</h2>
        <ul className="my-4 bg-red-50 py-6 px-12 rounded-lg list-disc">
          <li>画像の再販売や無断配布</li>
          <li>違法行為や他人の権利侵害に使用</li>
        </ul>
      </div>
    </>
  );
};

export default TermsOfService;
