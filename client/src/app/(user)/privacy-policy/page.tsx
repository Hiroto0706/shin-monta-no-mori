import Link from "next/link";

const PrivacyPolicy = () => {
  return (
    <>
      <div>
        <h1 className="text-3xl font-bold mb-6">プライバシーポリシー</h1>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトをご利用いただき、ありがとうございます。
          <br />
          このプライバシーポリシーは、お客様の個人情報の収集、使用、保護について説明するものです。
        </p>

        <h2 className="text-2xl font-bold mb-4">1. 個人情報の収集について</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトでは、ユーザーが画像を投稿する際やお問い合わせをする際に、名前、メールアドレスなどの個人情報を収集することがあります。
          <br />
          また、アクセス解析ツールを使用して、ユーザーの利用状況を収集することがあります。
        </p>

        <h2 className="text-2xl font-bold mb-4">2. 個人情報の利用について</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          収集した個人情報は、以下の目的で使用します：
          <ul className="list-disc ml-6">
            <li>サービスの提供および運営</li>
            <li>ユーザーサポートの提供</li>
            <li>新機能や更新情報の通知</li>
            <li>利用状況の分析およびサービス向上</li>
          </ul>
        </p>

        <h2 className="text-2xl font-bold mb-4">3. 個人情報の保護について</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトは、お客様の個人情報を適切に管理し、不正アクセス、紛失、破壊、改ざん、漏洩などの防止に努めます。また、個人情報の保護に関する法令を遵守します。
        </p>

        <h2 className="text-2xl font-bold mb-4">4. 第三者提供について</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトは、法令に基づく場合を除き、事前の同意なく個人情報を第三者に提供することはありません。ただし、以下の場合には個人情報を提供することがあります：
          <ul className="list-disc ml-6">
            <li>ユーザーの同意がある場合</li>
            <li>法令に基づく場合</li>
            <li>サービス提供のために必要な場合</li>
          </ul>
        </p>

        <h2 className="text-2xl font-bold mb-4">5. クッキーについて</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトでは、ユーザーの利便性向上のためにクッキーを使用することがあります。
          <br />
          クッキーは、ユーザーのブラウザに保存される小さなデータファイルです。
          <br />
          ブラウザの設定により、クッキーの受け入れを拒否することができますが、その場合、サイトの一部機能が利用できないことがあります。
        </p>

        <h2 className="text-2xl font-bold mb-4">
          6. プライバシーポリシーの変更について
        </h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          当サイトは、必要に応じてプライバシーポリシーを変更することがあります。
          <br />
          変更後のプライバシーポリシーは、当ページに掲載された時点から効力を生じます。
        </p>

        <h2 className="text-2xl font-bold mb-4">7. お問い合わせ</h2>
        <p className="my-4 bg-gray-100 p-6 rounded-lg">
          本ポリシーに関するお問い合わせは、以下のフォームからお願いいたします：
          <Link
            href="https://forms.gle/THqHAigzTZa7J9D28"
            className="text-blue-600 underline"
          >
            お問い合わせフォーム
          </Link>
        </p>
      </div>
    </>
  );
};

export default PrivacyPolicy;
