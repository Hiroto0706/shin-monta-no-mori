/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['storage.googleapis.com'],
  },
  staticPageGenerationTimeout: 480,
};

export default nextConfig;
