import React, { useState, useEffect } from 'react';
import {
  Typography,
  Button,
  TextField,
  Paper,
  Box,
  List,
  ListItem,
  ListItemText,
  CircularProgress,
  Card,
  CardContent,
} from '@mui/material';
import { getNotes, generateQuizTurn } from '../services/api';
import type { Message, Note } from '../types';

const Quiz: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const [notes, setNotes] = useState<Note[]>([]);
  const [selectedNoteIDs, setSelectedNoteIDs] = useState<Set<number>>(new Set());
  const [quizNoteIDs, setQuizNoteIDs] = useState<number[]>([]);

  useEffect(() => {
    const fetchNotes = async () => {
      try {
        const response = await getNotes();
        setNotes(response.data);
      } catch (error) {
        console.error('Failed to fetch notes', error);
      }
    };
    fetchNotes();
  }, []);

  const handleNoteSelection = (id: number) => {
    setSelectedNoteIDs((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  const handleSend = async () => {
    if (!input.trim()) return;

    const newMessages: Message[] = [...messages, { role: 'user', content: input }];
    setMessages(newMessages);
    setInput('');
    setLoading(true);

    try {
      const response = await generateQuizTurn(newMessages, quizNoteIDs);
      setMessages(response.data.messages);
      // The backend should return the note_ids used for the quiz
      if (response.data.note_ids) {
        setQuizNoteIDs(response.data.note_ids);
      }
    } catch (error) {
      console.error('Failed to generate quiz turn', error);
      setMessages([...newMessages, { role: 'model', content: 'Sorry, something went wrong.' }]);
    } finally {
      setLoading(false);
    }
  };

  const handleStartQuiz = async () => {
    if (selectedNoteIDs.size === 0) {
      alert('Please select at least one note to start the quiz.');
      return;
    }
    const noteIDs = Array.from(selectedNoteIDs);
    setQuizNoteIDs(noteIDs);
    setMessages([]);
    setLoading(true);
    try {
      const response = await generateQuizTurn([], noteIDs);
      setMessages(response.data.messages);
      if (response.data.note_ids) {
        setQuizNoteIDs(response.data.note_ids);
      }
    } catch (error) {
      console.error('Failed to start quiz', error);
      setMessages([{ role: 'model', content: 'Sorry, something went wrong.' }]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Paper elevation={3} sx={{ p: 3, display: 'flex', flexDirection: 'column', height: '80vh' }}>
      <Typography variant="h4" sx={{ mb: 2 }}>
        AI Quiz
      </Typography>
      {messages.length === 0 && !loading ? (
        <>
          <Typography variant="h6" sx={{ mb: 2 }}>
            Select Notes for the Quiz
          </Typography>
          <Box
            sx={{
              flexGrow: 1,
              overflowY: 'auto',
              mb: 2,
              display: 'grid',
              gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))',
              gap: 2,
            }}
          >
            {notes.map((note) => (
              <Card
                key={note.id}
                onClick={() => handleNoteSelection(note.id)}
                sx={{
                  cursor: 'pointer',
                  border: selectedNoteIDs.has(note.id) ? '2px solid' : '2px solid transparent',
                  borderColor: selectedNoteIDs.has(note.id) ? 'primary.main' : 'transparent',
                }}
              >
                <CardContent>
                  <Typography variant="body2" sx={{
                    display: '-webkit-box',
                    WebkitLineClamp: 4,
                    WebkitBoxOrient: 'vertical',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                  }}>
                    {note.content}
                  </Typography>
                </CardContent>
              </Card>
            ))}
          </Box>
          <Box sx={{ display: 'flex', justifyContent: 'center', mt: 2 }}>
            <Button variant="contained" color="primary" onClick={handleStartQuiz}>
              Start Quiz
            </Button>
          </Box>
        </>
      ) : (
        <>
          <Box sx={{ flexGrow: 1, overflowY: 'auto', mb: 2, p: 1, border: '1px solid #ddd', borderRadius: '4px' }}>
            <List>
              {messages.map((msg, index) => (
                <ListItem key={index}>
                  <ListItemText
                    primary={msg.content}
                    secondary={msg.role === 'user' ? 'You' : 'AI'}
                    sx={{ textAlign: msg.role === 'user' ? 'right' : 'left' }}
                  />
                </ListItem>
              ))}
              {loading && (
                <ListItem sx={{ justifyContent: 'center' }}>
                  <CircularProgress size={24} />
                </ListItem>
              )}
            </List>
          </Box>
          <Box sx={{ display: 'flex' }}>
            <TextField
              fullWidth
              variant="outlined"
              placeholder="Type your answer..."
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSend()}
              disabled={loading}
            />
            <Button variant="contained" onClick={handleSend} disabled={loading} sx={{ ml: 1 }}>
              Send
            </Button>
          </Box>
        </>
      )}
    </Paper>
  );
};

export default Quiz;