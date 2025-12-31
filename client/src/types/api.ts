/// <reference types="vite/client" />

// API Response Types based on Go backend
export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Deck {
  id: string;
  user_id: string;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Flashcard {
  id: string;
  user_id: string;
  front: string;
  back: string;
  deck_id: string;
  difficulty: number;
  interval: number;
  ease_factor: number;
  review_count: number;
  last_review: string | null;
  next_review: string | null;
  created_at: string;
  updated_at: string;
}

// Request DTOs
export interface CreateDeckRequest {
  name: string;
  description?: string;
}

export interface UpdateDeckRequest {
  name?: string;
  description?: string;
}

export interface CreateFlashcardRequest {
  front: string;
  back: string;
  deck_id: string;
}

export interface UpdateFlashcardRequest {
  front?: string;
  back?: string;
  difficulty?: number;
  interval?: number;
  ease_factor?: number;
  review_count?: number;
  last_review?: string | null;
  next_review?: string | null;
}

export interface ReviewFlashcardRequest {
  quality: number; // 0-5 rating
}

// API Response Types
export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  count: number;
  filters?: Record<string, unknown>;
}