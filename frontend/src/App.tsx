import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { HomePage } from './pages/HomePage';
import { ScheduleDetailPage } from './pages/ScheduleDetailPage';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/schedules/:id" element={<ScheduleDetailPage />} />
      </Routes>
    </Router>
  );
}

export default App;
