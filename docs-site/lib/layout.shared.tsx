import type { BaseLayoutProps } from 'fumadocs-ui/layouts/shared';
import { basementGrotesque } from './fonts';
import { appName, gitConfig } from './shared';

export function baseOptions(): BaseLayoutProps {
  return {
    nav: {
      // JSX supported
      title: (
        <span className={`${basementGrotesque.className} uppercase tracking-wide text-xl`}>
          {appName}
        </span>
      ),
    },
    githubUrl: `https://github.com/${gitConfig.user}/${gitConfig.repo}`,
  };
}
