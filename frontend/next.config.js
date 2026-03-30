/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone', // <--- Adicione isso,
  typescript: {
    ignoreBuildErrors: true,
  },
  // O correto para ESLint no Next.js é assim:
  //eslint: {
  //  ignoreDuringBuilds: true,
  //},
  //output: 'export',
}


module.exports = nextConfig
