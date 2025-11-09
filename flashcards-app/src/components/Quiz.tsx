import React, { useState } from 'react';
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
} from '@mui/material';
import { generateQuizTurn } from '../services/api';
import type { Message } from '../types';

const Quiz: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSend = async () => {
    if (!input.trim()) return;

    const newMessages: Message[] = [...messages, { role: 'user', content: input }];
    setMessages(newMessages);
    setInput('');
    setLoading(true);

    try {
      const response = await generateQuizTurn(newMessages);
      setMessages(response.data.messages);
    } catch (error) {
      console.error('Failed to generate quiz turn', error);
      // Optionally, add an error message to the chat
      setMessages([...newMessages, { role: 'model', content: 'Sorry, something went wrong.' }]);
    } finally {
      setLoading(false);
    }
  };

  const handleStartQuiz = async () => {
    setMessages([]);
    setLoading(true);
    try {
      const response = await generateQuizTurn([]); // Start with an empty conversation
      setMessages(response.data.messages);
    } catch (error) {
      console.error('Failed to start quiz', error);
      setMessages([{ role: 'model', content: 'Sorry, something went wrong.' }]);
    } finally {
      setLoading(false);
    }
  };


  return (
    <Paper elevation={3} sx={{ p: 3, display: 'flex', flexDirection: 'column', height: '70vh' }}>
      <Typography variant="h4" sx={{ mb: 2 }}>
        AI Quiz
      </Typography>
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
       {messages.length === 0 && !loading && (
        <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
          <Button variant="contained" color="primary" onClick={handleStartQuiz}>
            Start Quiz
          </Button>
        </Box>
      )}
      {messages.length > 0 && (
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
      )}
    </Paper>
  );
};

export default Quiz;
