import './App.css';
var React = require('react');



class App extends React.Component {
  constructor() {
    super();
    //定义数据
    this.state = { isServerError: false, serverError: "" };
  }
  //实现该生命周期的方法，react底层会自动在对应周期中调用该钩子方法
  componentDidMount() {
    this.init()
    //直接赋值，不会重新渲染html。必须调用setState方法才会监听到html是否变化，然后react才再重新渲染。
    //并非直接的双向数据绑定
    //this.state.text = "byebye world";
  }
  changeData = (data: any) => {
    console.log("changeData:", data)
    this.setState({
      isServerError: data.isServerError,
      serverError: data.serverError
    });
  }
  callServerStart = () => {
    window.electron.ipcRenderer.once('ipc-example', () => {
      this.changeData({ isServerError: false, serverError: "" })
    });
    window.electron.ipcRenderer.sendMessage('ipc-example', ['startServer']);
  };
  init = () => {
    // calling IPC exposed from preload script
    window.electron.ipcRenderer.once('ipc-example', () => {
      //
      window.electron.ipcRenderer.once('ipc-example', (arg: any) => {
        let data = JSON.parse("" + arg)
        this.changeData(data)
      });
      window.electron.ipcRenderer.sendMessage('ipc-example', ['info']);
    });
    window.electron.ipcRenderer.sendMessage('ipc-example', ['ping']);
  }
  render() {
    return (
      <div>
        {
          this.state.isServerError ?
            <div className="app-system-message-info">
              Team · IDE 服务启动异常 {this.state.serverError}，<a onClick={this.callServerStart} className="reset-btn">点击重试</a>
            </div>
            :
            <div className="app-system-message-info">
              Team · IDE 服务启动中，请稍后~~~
            </div>
        }
      </div>
    );
  }
}


export default App;