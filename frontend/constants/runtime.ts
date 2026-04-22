export function isDesktopRuntime(): boolean {
  if (typeof window === 'undefined') return false;

  const envRuntime = String(process.env.NEXT_PUBLIC_APP_RUNTIME ?? '').trim().toLowerCase();
  if (envRuntime === 'desktop' || envRuntime === 'binary') return true;

  return window.location.port === '9000';
}

export function isWebRuntime(): boolean {
  return !isDesktopRuntime();
}

