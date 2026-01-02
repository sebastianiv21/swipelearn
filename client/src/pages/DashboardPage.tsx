import { useSession, signOut } from '../lib/auth-client';
import { Link, useNavigate } from 'react-router-dom';
import { BookOpen, BarChart3, LogOut, Clock } from 'lucide-react';
import { useDecks } from '../hooks/useDecks';
import { useDueFlashcards } from '../hooks/useFlashcards';

export default function DashboardPage() {
  const session = useSession();
  const navigate = useNavigate();
  
  // Fetch real data for statistics
  const { data: decks, isLoading: isLoadingDecks } = useDecks();
  const { data: dueCards, isLoading: isLoadingDueCards } = useDueFlashcards();

  const handleLogout = async () => {
    try {
      await signOut();
      navigate('/login');
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };

  if (!session.user) {
    return null; // Should be handled by ProtectedRoute
  }

  return (
    <div className="min-h-screen bg-gray-50 safe-top safe-bottom">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold text-gray-900">
              Welcome back, {session.user ? session.user.name.split(' ')[0] : 'Guest'}!
            </h1>
            <div className="flex items-center space-x-4">
              <Link
                to="/profile"
                className="text-sm text-gray-600 hover:text-gray-900"
              >
                Profile
              </Link>
              <button
                onClick={handleLogout}
                className="flex items-center space-x-1 text-sm text-red-600 hover:text-red-700"
              >
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="max-w-4xl mx-auto px-4 py-8">
        {/* Quick Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          <div className="card text-center">
            <BookOpen className="h-8 w-8 text-blue-600 mx-auto mb-2" />
            <h3 className="font-semibold text-gray-900">Decks</h3>
            <p className="text-2xl font-bold text-gray-900">
              {isLoadingDecks ? '...' : decks?.length || 0}
            </p>
          </div>
          <div className="card text-center">
            <Clock className="h-8 w-8 text-green-600 mx-auto mb-2" />
            <h3 className="font-semibold text-gray-900">Due Today</h3>
            <p className="text-2xl font-bold text-gray-900">
              {isLoadingDueCards ? '...' : dueCards?.length || 0}
            </p>
          </div>
          <div className="card text-center">
            <BarChart3 className="h-8 w-8 text-purple-600 mx-auto mb-2" />
            <h3 className="font-semibold text-gray-900">Study Streak</h3>
            <p className="text-2xl font-bold text-gray-900">0 days</p>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="space-y-4">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <Link
              to="/decks"
              className="card hover:shadow-md transition-shadow duration-200"
            >
              <h3 className="font-semibold text-gray-900 mb-2">Manage Decks</h3>
              <p className="text-sm text-gray-600">Create and edit your flashcard decks</p>
            </Link>
            
            <Link
              to="/study"
              className="card hover:shadow-md transition-shadow duration-200"
            >
              <h3 className="font-semibold text-gray-900 mb-2">Start Studying</h3>
              <p className="text-sm text-gray-600">
                Review cards using spaced repetition
                {dueCards && dueCards.length > 0 && (
                  <span className="block text-green-600 font-semibold mt-1">
                    {dueCards.length} cards due today!
                  </span>
                )}
              </p>
            </Link>
          </div>
        </div>

        {/* Recent Activity */}
        <div className="mt-8">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h2>
          <div className="card">
            <p className="text-gray-600 text-center py-8">
              No recent activity. Start studying to see your progress here!
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}