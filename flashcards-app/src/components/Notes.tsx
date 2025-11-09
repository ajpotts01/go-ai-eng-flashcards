import React, { useState, useEffect } from 'react';
import {
  Typography,
  Button,
  List,
  ListItem,
  ListItemText,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  TextField,
  DialogActions,
  CircularProgress,
  Box,
  Paper,
} from '@mui/material';
import { Edit, Delete, Add } from '@mui/icons-material';
import { getNotes, createNote, updateNote, deleteNote } from '../services/api';
import type { Note } from '../types';

const Notes: React.FC = () => {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(true);
  const [open, setOpen] = useState(false);
  const [currentNote, setCurrentNote] = useState<Note | null>(null);
  const [newContent, setNewContent] = useState('');

  useEffect(() => {
    fetchNotes();
  }, []);

  const fetchNotes = async () => {
    try {
      setLoading(true);
      const response = await getNotes();
      setNotes(response.data);
    } catch (error) {
      console.error('Failed to fetch notes', error);
    } finally {
      setLoading(false);
    }
  };

  const handleOpen = (note: Note | null = null) => {
    setCurrentNote(note);
    setNewContent(note ? note.content : '');
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setCurrentNote(null);
    setNewContent('');
  };

  const handleSave = async () => {
    try {
      if (currentNote) {
        await updateNote(currentNote.id, newContent);
      } else {
        await createNote(newContent);
      }
      fetchNotes();
      handleClose();
    } catch (error) {
      console.error('Failed to save note', error);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteNote(id);
      fetchNotes();
    } catch (error) {
      console.error('Failed to delete note', error);
    }
  };

  return (
    <Paper elevation={3} sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h4">Notes</Typography>
        <Button variant="contained" startIcon={<Add />} onClick={() => handleOpen()}>
          Add Note
        </Button>
      </Box>
      {loading ? (
        <CircularProgress />
      ) : (
        <List>
          {notes.map((note) => (
            <ListItem
              key={note.id}
              secondaryAction={
                <>
                  <IconButton edge="end" aria-label="edit" onClick={() => handleOpen(note)}>
                    <Edit />
                  </IconButton>
                  <IconButton edge="end" aria-label="delete" onClick={() => handleDelete(note.id)}>
                    <Delete />
                  </IconButton>
                </>
              }
            >
              <ListItemText primary={note.content} secondary={`Last updated: ${new Date(note.updated_at).toLocaleString()}`} />
            </ListItem>
          ))}
        </List>
      )}
      <Dialog open={open} onClose={handleClose} fullWidth maxWidth="sm">
        <DialogTitle>{currentNote ? 'Edit Note' : 'Add Note'}</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Content"
            type="text"
            fullWidth
            variant="outlined"
            value={newContent}
            onChange={(e) => setNewContent(e.target.value)}
            multiline
            rows={4}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button onClick={handleSave} variant="contained">
            Save
          </Button>
        </DialogActions>
      </Dialog>
    </Paper>
  );
};

export default Notes;
