import localFont from 'next/font/local';

// Basement Grotesque Black — self-hosted wordmark face for the "Totality" logo.
// Source: github.com/basementstudio/basement-grotesque (free license).
export const basementGrotesque = localFont({
  src: '../app/fonts/BasementGrotesque-Black.woff2',
  weight: '800',
  display: 'swap',
});
