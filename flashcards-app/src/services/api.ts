import axios from 'axios';
import type { Note } from '../types';
import type { Message } from '../types';

const apiClient = axios.create({
  baseURL: 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

export const getNotes = () => apiClient.get<Note[]>('/notes');
export const getNoteById = (id: number) => apiClient.get<Note>(`/notes/${id}`);
export const createNote = (content: string) => apiClient.post<Note>('/notes', { content });
export const updateNote = (id: number, content: string) => apiClient.put<Note>(`/notes/${id}`, { content });
export const deleteNote = (id: number) => apiClient.delete(`/notes/${id}`);

export const generateQuizTurn = (messages: Message[]) => apiClient.post<{ messages: Message[] }>('/quiz', { messages });

export default apiClient;
