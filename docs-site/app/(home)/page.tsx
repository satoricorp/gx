import { redirect } from 'next/navigation';

// The docs are the whole site — land visitors directly on them.
export default function HomePage() {
  redirect('/docs');
}
