import { BookOpen } from 'lucide-react';

export default function DecksPage() {
  return (
    <div className="min-h-screen bg-gray-50 safe-top safe-bottom">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold text-gray-900">My Decks</h1>
            <button className="btn-primary flex items-center gap-2">
              <BookOpen className="h-4 w-4" />
              New Deck
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 py-8">
        <div className="text-center py-16">
          <BookOpen className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-gray-900 mb-2">No decks yet</h2>
          <p className="text-gray-600 mb-6">
            Create your first deck to start studying with flashcards
          </p>
          <button className="btn-primary">
            Create Your First Deck
          </button>
        </div>
      </div>
    </div>
  );
}