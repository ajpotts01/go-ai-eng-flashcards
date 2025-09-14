ALTER TABLE flashcards.flashcards RENAME TO notes;
ALTER INDEX flashcards.idx_flashcards_created_at RENAME TO idx_notes_created_at;
