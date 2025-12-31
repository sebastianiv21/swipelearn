import { Link } from 'react-router-dom';
import { Home, BookOpen } from 'lucide-react';

export default function NotFoundPage() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 safe-top safe-bottom">
      <div className="text-center max-w-md mx-4">
        <div className="card">
          <h1 className="text-6xl font-bold text-gray-900 mb-4">404</h1>
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Page Not Found</h2>
          <p className="text-gray-600 mb-8">
            The page you're looking for doesn't exist or has been moved.
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              to="/"
              className="btn-primary flex items-center justify-center gap-2"
            >
              <Home className="h-4 w-4" />
              Dashboard
            </Link>
            
            <Link
              to="/decks"
              className="btn-secondary flex items-center justify-center gap-2"
            >
              <BookOpen className="h-4 w-4" />
              Decks
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}