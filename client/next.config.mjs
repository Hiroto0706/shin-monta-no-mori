/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['storage.googleapis.com'],
  },
  staticPageGenerationTimeout: 120,
};

export default nextConfig;
