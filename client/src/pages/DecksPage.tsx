import { useState } from 'react';
import { Link } from 'react-router-dom';
import { BookOpen, Plus, Trash2, Clock, Layers } from 'lucide-react';
import { useDecks, useCreateDeck, useDeleteDeck } from '../hooks/useDecks';
import { useFlashcards } from '../hooks/useFlashcards';
import type { CreateDeckRequest } from '../types/api';

export default function DecksPage() {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [formData, setFormData] = useState<CreateDeckRequest>({
    name: '',
    description: '',
  });

  const { data: decks, isLoading } = useDecks();
  const createDeckMutation = useCreateDeck();
  const deleteDeckMutation = useDeleteDeck();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name.trim()) return;

    try {
      await createDeckMutation.mutateAsync(formData);
      setFormData({ name: '', description: '' });
      setShowCreateForm(false);
    } catch (error) {
      console.error('Failed to create deck:', error);
    }
  };

  const handleDelete = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this deck? This action cannot be undone.')) {
      try {
        await deleteDeckMutation.mutateAsync(id);
      } catch (error) {
        console.error('Failed to delete deck:', error);
      }
    }
  };
  const DeckCard = ({ deck }: { deck: { id: string; name: string; description?: string; created_at: string } }) => {
    const { data: cards } = useFlashcards(deck.id);
    const cardCount = cards?.length || 0;

    return (
      <div className="card hover:shadow-md transition-shadow duration-200">
        <div className="flex items-start justify-between">
          <div className="flex-1">
            <h3 className="font-semibold text-gray-900 mb-2">{deck.name}</h3>
            {deck.description && (
              <p className="text-sm text-gray-600 mb-3">{deck.description}</p>
            )}
            <div className="flex items-center gap-4 text-sm text-gray-500">
              <div className="flex items-center gap-1">
                <Layers className="h-4 w-4" />
                <span>{cardCount} cards</span>
              </div>
              <div className="flex items-center gap-1">
                <Clock className="h-4 w-4" />
                <span>Created {new Date(deck.created_at).toLocaleDateString()}</span>
              </div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Link
              to={`/study?deck=${deck.id}`}
              className="btn-secondary btn-sm"
            >
              Study
            </Link>
            <button
              onClick={() => handleDelete(deck.id)}
              className="text-red-600 hover:text-red-700 p-2"
              disabled={deleteDeckMutation.isPending}
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50 safe-top safe-bottom">
      {/* Header */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-4xl mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-xl font-bold text-gray-900">My Decks</h1>
            <button
              onClick={() => setShowCreateForm(true)}
              className="btn-primary flex items-center gap-2"
            >
              <Plus className="h-4 w-4" />
              New Deck
            </button>
          </div>
        </div>
      </div>

      {/* Create Form Modal */}
      {showCreateForm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Create New Deck</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label htmlFor="name" className="block text-sm font-medium text-gray-700 mb-1">
                  Deck Name *
                </label>
                <input
                  type="text"
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  className="input-field"
                  placeholder="e.g., Spanish Vocabulary"
                  required
                  autoFocus
                />
              </div>
              <div>
                <label htmlFor="description" className="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                  className="input-field"
                  rows={3}
                  placeholder="Optional description for your deck"
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="button"
                  onClick={() => setShowCreateForm(false)}
                  className="btn-secondary flex-1"
                >
                  Cancel
                </button>
                <button
                  type="submit"
                  className="btn-primary flex-1"
                  disabled={createDeckMutation.isPending || !formData.name.trim()}
                >
                  {createDeckMutation.isPending ? 'Creating...' : 'Create Deck'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 py-8">
        {isLoading ? (
          <div className="text-center py-16">
            <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <p className="text-gray-600 mt-4">Loading your decks...</p>
          </div>
        ) : decks && decks.length > 0 ? (
          <div className="space-y-4">
            <p className="text-gray-600">
              You have {decks.length} deck{decks.length !== 1 ? 's' : ''}
            </p>
            {decks.map((deck) => (
              <DeckCard key={deck.id} deck={deck} />
            ))}
          </div>
        ) : (
          <div className="text-center py-16">
            <BookOpen className="h-16 w-16 text-gray-400 mx-auto mb-4" />
            <h2 className="text-xl font-semibold text-gray-900 mb-2">No decks yet</h2>
            <p className="text-gray-600 mb-6">
              Create your first deck to start studying with flashcards
            </p>
            <button
              onClick={() => setShowCreateForm(true)}
              className="btn-primary"
            >
              Create Your First Deck
            </button>
          </div>
        )}
      </div>
    </div>
  );
}