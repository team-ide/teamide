// Modules to control application life and create native browser window
const { app, Menu, BrowserWindow, Tray } = require('electron')
const path = require('path')

var win
function createWindow() {
  // Create the browser window.
  win = new BrowserWindow({
    title: "Team · IDE",
    width: 1440,
    height: 900,
    autoHideMenuBar: true,
    // frame: false,
    icon: path.join(__dirname, './public/static/logo.png'),
    menuBarVisible: false, //菜单栏是否可见
    webPreferences: {
      preload: path.join(__dirname, 'preload.js')
    }
  })
  // win.menuBarVisible = false;

  Menu.setApplicationMenu(null)

  win.loadURL('https://baidu.com')

  // Open the DevTools.
  // mainWindow.webContents.openDevTools()
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  createWindow()

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', function () {
  if (process.platform !== 'darwin') app.quit()
})

// In this file you can include the rest of your app's specific main process
// code. You can also put them in separate files and require them here.

//定义全局系统图标变量
let tray = null
app.on('ready', async () => {
  tray = new Tray(path.join(__dirname, './public/static/logo.png'))
  const contextMenu = Menu.buildFromTemplate([
    {
      label: '退出',
      click: function () {
        app.quit()
      }
    }
  ])
  tray.setToolTip('Team IDE')
  //显示程序页面
  tray.on('click', () => {
    win.show()
  })
  tray.setContextMenu(contextMenu)
})