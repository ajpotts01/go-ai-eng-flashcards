export interface Note {
  id: number;
  content: string;
  created_at: string;
  updated_at: string;
}

export type Message = {
  role: string;
  content: string;
};