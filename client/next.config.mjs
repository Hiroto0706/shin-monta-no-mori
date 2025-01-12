/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['storage.googleapis.com'],
  },
  staticPageGenerationTimeout: 960,
};

export default nextConfig;
