import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import { useSession } from '../lib/auth-client';
import { Loader2 } from 'lucide-react';

// Import route components (to be created)
import LoginPage from '../pages/LoginPage';
import RegisterPage from '../pages/RegisterPage';
import DashboardPage from '../pages/DashboardPage';
import DecksPage from '../pages/DecksPage';
import StudyPage from '../pages/StudyPage';
import NotFoundPage from '../pages/NotFoundPage';

// Protected Route Component
function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const session = useSession();

  if (session.isPending) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (!session.user) {
    return <LoginPage />;
  }

  return <>{children}</>;
}

// Router configuration
const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/register',
    element: <RegisterPage />,
  },
  {
    path: '/',
    element: (
      <ProtectedRoute>
        <DashboardPage />
      </ProtectedRoute>
    ),
  },
  {
    path: '/decks',
    element: (
      <ProtectedRoute>
        <DecksPage />
      </ProtectedRoute>
    ),
  },
  {
    path: '/study',
    element: (
      <ProtectedRoute>
        <StudyPage />
      </ProtectedRoute>
    ),
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
]);

// Router Provider Component
export function Router() {
  return <RouterProvider router={router} />;
}