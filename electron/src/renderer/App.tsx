import { MemoryRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';

const Index = () => {
  return (
    <div>
      <div className="app-system-message-info">
        Team · IDE 服务启动中，请稍后~~~
      </div>
    </div>
  );
};

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Index />} />
      </Routes>
    </Router>
  );
}
