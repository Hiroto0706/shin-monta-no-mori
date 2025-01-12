/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['storage.googleapis.com'],
  },
  staticPageGenerationTimeout: 240,
};

export default nextConfig;
