import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,
  output: 'export',
  images: { unoptimized: true },
  allowedDevOrigins: ['test.antipratik.com'],
};

export default nextConfig;
