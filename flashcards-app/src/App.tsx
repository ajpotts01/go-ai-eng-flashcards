import React, { useState } from 'react';
import { AppBar, Toolbar, Typography, Container, Tabs, Tab, Box } from '@mui/material';
import Notes from './components/Notes';
import Quiz from './components/Quiz';
import './App.css';

function App() {
  const [selectedTab, setSelectedTab] = useState(0);

  const handleChange = (_event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };

  return (
    <>
      <AppBar position="static">
        <Toolbar>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Flashcards AI
          </Typography>
        </Toolbar>
      </AppBar>
      <Container maxWidth="lg" sx={{ mt: 4 }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs value={selectedTab} onChange={handleChange} aria-label="basic tabs example">
            <Tab label="Notes" />
            <Tab label="Quiz" />
          </Tabs>
        </Box>
        <Box sx={{ p: 3 }}>
          {selectedTab === 0 && <Notes />}
          {selectedTab === 1 && <Quiz />}
        </Box>
      </Container>
    </>
  );
}

export default App;
