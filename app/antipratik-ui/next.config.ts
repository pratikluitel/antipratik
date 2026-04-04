import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactCompiler: true,
  // Disable static export during isolated UI build checks so dynamic routes
  // don't require backend-driven params.
  output: process.env.NEXT_DISABLE_EXPORT ? undefined : 'export',
  images: { unoptimized: true },
  allowedDevOrigins: ['test.antipratik.com'],
};

export default nextConfig;
