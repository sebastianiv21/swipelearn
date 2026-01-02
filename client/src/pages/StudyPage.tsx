import { useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { BookOpen, Eye, EyeOff, ArrowRight, Home } from 'lucide-react';
import { useDueFlashcards } from '../hooks/useFlashcards';
import { useReviewFlashcard } from '../hooks/useStudySession';
import { useDecks } from '../hooks/useDecks';

export default function StudyPage() {
  const [searchParams] = useSearchParams();
  const deckId = searchParams.get('deck');
  
  const [currentCardIndex, setCurrentCardIndex] = useState(0);
  const [showBack, setShowBack] = useState(false);
  const [sessionStats, setSessionStats] = useState({
    studied: 0,
    easy: 0,
    good: 0,
    hard: 0,
  });

  const { data: dueCards, isLoading, error } = useDueFlashcards(deckId || undefined);
  const { data: decks } = useDecks();
  const reviewMutation = useReviewFlashcard();

  const currentCard = dueCards?.[currentCardIndex];
  const isSessionComplete = dueCards && currentCardIndex >= dueCards.length;



  const handleReview = async (quality: number) => {
    if (!currentCard) return;

    try {
      await reviewMutation.mutateAsync({
        id: currentCard.id,
        quality,
      });

      // Update stats
      const rating = quality >= 3 ? 'good' : quality >= 2 ? 'hard' : 'easy';
      setSessionStats(prev => ({
        ...prev,
        studied: prev.studied + 1,
        [rating]: prev[rating as keyof typeof prev] + 1,
      }));

      // Move to next card
      setCurrentCardIndex(prev => prev + 1);
      setShowBack(false);
    } catch (error) {
      console.error('Failed to review card:', error);
    }
  };

  const getDeckName = () => {
    if (!deckId || !decks) return null;
    const deck = decks.find(d => d.id === deckId);
    return deck ? <span className="text-gray-600">({deck.name})</span> : null;
  };

  return (
    <div className="min-h-screen bg-gray-50 safe-top safe-bottom">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold text-gray-900">Study Session</h1>
            <Link to="/" className="btn-secondary flex items-center gap-2">
              <Home className="h-4 w-4" />
              End Session
            </Link>
          </div>
        </div>
      </div>

      {/* Loading State */}
      {isLoading && (
        <div className="max-w-2xl mx-auto px-4 py-16">
          <div className="text-center">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="text-gray-600 mt-4">Loading study session...</p>
          </div>
        </div>
      )}

      {/* Error State */}
      {error && (
        <div className="max-w-2xl mx-auto px-4 py-16">
          <div className="text-center">
            <div className="text-red-600 mb-4">Failed to load study session</div>
            <p className="text-gray-600">Please try again later</p>
          </div>
        </div>
      )}

      {/* No Cards State */}
      {!isLoading && !error && dueCards && dueCards.length === 0 && (
        <div className="max-w-2xl mx-auto px-4 py-8">
          <div className="bg-blue-50 border-b border-blue-200 px-4 py-3">
            <div className="max-w-4xl mx-auto flex items-center justify-between">
              <div className="text-sm text-blue-900">
                Progress: 0 / 0 cards
              </div>
              <div className="flex items-center gap-4 text-sm text-blue-900">
                <span>Studied: {sessionStats.studied}</span>
                <div className="h-4 w-px bg-blue-300"></div>
                <span>Easy: {sessionStats.easy}</span>
                <span>Good: {sessionStats.good}</span>
                <span>Hard: {sessionStats.hard}</span>
              </div>
            </div>
          </div>
          <div className="text-center py-16">
            <BookOpen className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-gray-900 mb-2">Ready to Study!</h2>
            <p className="text-gray-600 mb-6">
              You have no cards due for review right now.
            </p>
            <div className="space-y-4">
              <Link to="/decks" className="btn-primary inline-flex items-center gap-2">
                <ArrowRight className="h-4 w-4" />
                Manage Decks
              </Link>
              <div className="text-sm text-gray-600">
                or come back later when cards are due for review
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Study Session */}
      {!isLoading && !error && dueCards && dueCards.length > 0 && (
        <>
          <div className="bg-blue-50 border-b border-blue-200 px-4 py-3">
            <div className="max-w-4xl mx-auto flex items-center justify-between">
              <div className="text-sm text-blue-900">
                Progress: {currentCardIndex} / {dueCards?.length || 0} cards
              </div>
              <div className="flex items-center gap-4 text-sm text-blue-900">
                <span>Studied: {sessionStats.studied}</span>
                <div className="h-4 w-px bg-blue-300"></div>
                <span>Easy: {sessionStats.easy}</span>
                <span>Good: {sessionStats.good}</span>
                <span>Hard: {sessionStats.hard}</span>
              </div>
            </div>
          </div>
          {!isSessionComplete ? (
            <div className="max-w-2xl mx-auto px-4 py-8">
              <div className="text-center mb-8">
                <h2 className="text-lg font-semibold text-gray-600 mb-2">
                  Card {currentCardIndex + 1} of {dueCards?.length}
                </h2>
                <p className="text-sm text-gray-500">
                  {getDeckName()}
                </p>
              </div>

              <div className="card mb-8 min-h-64 flex items-center justify-center">
                <div className="w-full">
                  <div className="text-center mb-4">
                    <button
                      onClick={() => setShowBack(!showBack)}
                      className="btn-secondary btn-sm"
                    >
                      {showBack ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      {showBack ? ' Hide Answer' : ' Show Answer'}
                    </button>
                  </div>
                  
                  <div className="text-center">
                    <h3 className="text-xl font-semibold text-gray-900 mb-6">
                      {currentCard?.front}
                    </h3>
                    
                    {showBack && (
                      <div className="border-t pt-6">
                        <p className="text-lg text-gray-700">
                          {currentCard?.back}
                        </p>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {showBack && (
                <div className="space-y-4">
                  <p className="text-center text-sm text-gray-600 mb-4">
                    How well did you know this card?
                  </p>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                    <button
                      onClick={() => handleReview(0)}
                      className="btn-secondary flex flex-col items-center gap-2 p-4"
                      disabled={reviewMutation.isPending}
                    >
                      <span className="text-lg">😰</span>
                      <span className="text-sm">Again</span>
                      <span className="text-xs text-gray-500">(0)</span>
                    </button>
                    <button
                      onClick={() => handleReview(2)}
                      className="btn-secondary flex flex-col items-center gap-2 p-4"
                      disabled={reviewMutation.isPending}
                    >
                      <span className="text-lg">😕</span>
                      <span className="text-sm">Hard</span>
                      <span className="text-xs text-gray-500">(2)</span>
                    </button>
                    <button
                      onClick={() => handleReview(3)}
                      className="btn-secondary flex flex-col items-center gap-2 p-4"
                      disabled={reviewMutation.isPending}
                    >
                      <span className="text-lg">😊</span>
                      <span className="text-sm">Good</span>
                      <span className="text-xs text-gray-500">(3)</span>
                    </button>
                    <button
                      onClick={() => handleReview(5)}
                      className="btn-primary flex flex-col items-center gap-2 p-4"
                      disabled={reviewMutation.isPending}
                    >
                      <span className="text-lg">🎉</span>
                      <span className="text-sm">Easy</span>
                      <span className="text-xs text-gray-500">(5)</span>
                    </button>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <div className="max-w-2xl mx-auto px-4 py-16">
              <div className="text-center">
                <div className="text-6xl mb-4">🎉</div>
                <h2 className="text-2xl font-bold text-gray-900 mb-4">
                  Session Complete!
                </h2>
                <div className="card mb-6">
                  <h3 className="font-semibold text-gray-900 mb-4">Session Summary</h3>
                  <div className="grid grid-cols-2 gap-4 text-left">
                    <div>
                      <span className="text-gray-600">Cards Studied:</span>
                      <span className="font-semibold ml-2">{sessionStats.studied}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">Easy:</span>
                      <span className="font-semibold ml-2 text-green-600">{sessionStats.easy}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">Good:</span>
                      <span className="font-semibold ml-2 text-blue-600">{sessionStats.good}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">Hard:</span>
                      <span className="font-semibold ml-2 text-orange-600">{sessionStats.hard}</span>
                    </div>
                  </div>
                </div>
                <div className="space-y-3">
                  <Link to="/" className="btn-primary inline-flex items-center gap-2">
                    <Home className="h-4 w-4" />
                    Back to Dashboard
                  </Link>
                  <Link 
                    to="/decks" 
                    className="btn-secondary inline-flex items-center gap-2"
                  >
                    <BookOpen className="h-4 w-4" />
                    Manage Decks
                  </Link>
                </div>
              </div>
            </div>
          )}
        </>
      )}
    </div>
  );
}