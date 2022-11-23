import { MemoryRouter as Router, Routes, Route } from 'react-router-dom';
import './App.css';

const AppData = { isServerError: false }
const Index = () => {
  return (
    <div>
      {
        AppData.isServerError ?
          <div className="app-system-message-info">
            Team · IDE 服务启动异常，点击重试
          </div>
          :
          <div className="app-system-message-info">
            Team · IDE 服务启动中，请稍后~~~
          </div>
      }
    </div>
  );
};

// calling IPC exposed from preload script
window.electron.ipcRenderer.once('ipc-example', (arg) => {
  // eslint-disable-next-line no-console
  console.log("window.electron.ipcRenderer.once ipc-example:", arg);
  if (arg == "pong") {
    //
    console.log("pong");
    window.electron.ipcRenderer.sendMessage('ipc-example', ['info']);
  } else if (arg == "serverStarted") {
    //
    console.log("serverStarted");
    AppData.isServerError = false
  } else {
    let data = JSON.parse("" + arg)
    if (data.isServerError) {
      AppData.isServerError = true
    }
  }
});
window.electron.ipcRenderer.sendMessage('ipc-example', ['ping']);

export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Index />} />
      </Routes>
    </Router>
  );
}
