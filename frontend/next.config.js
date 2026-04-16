
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
  typescript: {
    ignoreBuildErrors: true,
  },
  images: {
    unoptimized: true,
  },
}

module.exports = nextConfig


// /** @type {import('next').NextConfig} */
// const nextConfig = {
//   output: 'standalone', // <--- Adicione isso,
//   typescript: {
//     ignoreBuildErrors: true,
//   },
//   // O correto para ESLint no Next.js é assim:
//   //eslint: {
//   //  ignoreDuringBuilds: true,
//   //},
//   //output: 'export',
// }


// module.exports = nextConfig
