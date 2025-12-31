import { BookOpen } from 'lucide-react';

export default function StudyPage() {
  return (
    <div className="min-h-screen bg-gray-50 safe-top safe-bottom">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold text-gray-900">Study Session</h1>
            <button className="btn-secondary">End Session</button>
          </div>
        </div>
      </div>

      {/* Study Area */}
      <div className="max-w-2xl mx-auto px-4 py-8">
        <div className="text-center py-16">
          <BookOpen className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h2 className="text-xl font-semibold text-gray-900 mb-2">Ready to Study!</h2>
          <p className="text-gray-600 mb-6">
            You have no cards due for review right now.
          </p>
          <div className="space-y-4">
            <button className="btn-primary">
              Create Flashcards
            </button>
            <div className="text-sm text-gray-600">
              or come back later when cards are due for review
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}