/* eslint global-require: off, no-console: off, promise/always-return: off */

/**
 * This module executes inside of electron's main process. You can start
 * electron renderer process from here and communicate with the other processes
 * through IPC.
 *
 * When running `npm run build` or `npm run build:main`, this file is compiled to
 * `./src/main.js` using webpack. This gives us some performance wins.
 */
import path from 'path';
import { app, BrowserWindow, shell, ipcMain, Menu, Tray, screen } from 'electron';
import { autoUpdater } from 'electron-updater';
import log from 'electron-log';
import MenuBuilder from './menu';
import { resolveHtmlPath, source } from './util';



var fs = require("fs")
const child_process = require('child_process');

process.env['ELECTRON_DISABLE_SECURITY_WARNINGS'] = 'true'

export default class AppUpdater {
  constructor() {
    log.transports.file.level = 'info';
    autoUpdater.logger = log;
    autoUpdater.checkForUpdatesAndNotify();
  }
}

let mainWindow: BrowserWindow | null = null;

ipcMain.on('ipc-example', async (event, arg) => {
  const msgTemplate = (pingPong: string) => `IPC test: ${pingPong}`;
  console.log(msgTemplate(arg));
  event.reply('ipc-example', msgTemplate('pong'));
});
ipcMain.on('open-new-window', async (_event: any, config: any) => {
  config = config || {};
  console.log("open-new-window:config:", config);
  source.addBrowserView(config)
});

if (process.env.NODE_ENV === 'production') {
  const sourceMapSupport = require('source-map-support');
  sourceMapSupport.install();
}

const isDev = process.env.NODE_ENV === 'development';
const isDebug =
  process.env.NODE_ENV === 'development' || process.env.DEBUG_PROD === 'true';

if (isDebug) {
  require('electron-debug')();
}

const installExtensions = async () => {
  const installer = require('electron-devtools-installer');
  const forceDownload = !!process.env.UPGRADE_EXTENSIONS;
  const extensions = ['REACT_DEVELOPER_TOOLS'];

  return installer
    .default(
      extensions.map((name) => installer[name]),
      forceDownload
    )
    .catch(console.log);
};
let serverProcess: any;


const RESOURCES_PATH = app.isPackaged
  ? path.join(process.resourcesPath, 'assets')
  : path.join(__dirname, '../../assets');

const getAssetPath = (...paths: string[]): string => {
  return path.join(RESOURCES_PATH, ...paths);
};


const ROOT_PATH = app.isPackaged
  ? path.join(process.resourcesPath, '')
  : path.join(__dirname, '../../');

export const getRootPath = (...paths: string[]): string => {
  return path.join(ROOT_PATH, ...paths);
};
export const iconPath = getAssetPath('icon.png');
source.iconPath = iconPath;
export const icon16Path = getAssetPath('icon-16.png');
source.icon16Path = icon16Path;
export const icon32Path = getAssetPath('icon-32.png');
source.icon32Path = icon32Path;

let serverUrl = resolveHtmlPath('index.html')

const viewWindowList: BrowserWindow[] = []
const createWindow = async () => {

  if (mainWindow != null && !mainWindow.isDestroyed()) {
    mainWindow.show()
    return
  }

  if (isDebug) {
    await installExtensions();
  }

  const { width, height } = screen.getPrimaryDisplay().workAreaSize;//获取到屏幕的宽度和高度

  source.windowWidth = 1440;
  source.windowHeight = 900;
  if (source.windowWidth > width) {
    source.windowWidth = (width - 40);
  }
  if (source.windowHeight > height) {
    source.windowHeight = (height - 40);
  }

  mainWindow = new BrowserWindow({
    title: "Team · IDE",
    show: false,
    width: source.windowWidth,
    height: source.windowHeight,
    icon: iconPath,
    autoHideMenuBar: true,
    webPreferences: {
      preload: app.isPackaged
        ? path.join(__dirname, 'preload.js')
        : path.join(__dirname, '../../.erb/dll/preload.js'),
    },
  });

  source.addBrowserView = (config: any) => {
    let viewWindow = new BrowserWindow({
      title: config.title || "Team · IDE",
      show: false,
      width: source.windowWidth,
      height: source.windowHeight,
      icon: iconPath,
      autoHideMenuBar: true,
      skipTaskbar: true,
      webPreferences: {
        preload: app.isPackaged
          ? path.join(__dirname, 'preload.js')
          : path.join(__dirname, '../../.erb/dll/preload.js'),
      },
    });
    viewWindow.loadURL(config.url);
    viewWindow.show()
    viewWindowList.push(viewWindow)

    viewWindow.on('closed', () => {
      let index = viewWindowList.indexOf(viewWindow)
      if (index >= 0) {
        viewWindowList.splice(index, 1)
      }
      if (mainWindow != null && !mainWindow.isDestroyed()) {
        mainWindow.webContents.send('close-open-window', config);
      }
    });
  }

  source.mainWindow = mainWindow;

  mainWindow.loadURL(serverUrl);

  mainWindow.on('ready-to-show', () => {
    if (!mainWindow) {
      throw new Error('"mainWindow" is not defined');
    }
    if (process.env.START_MINIMIZED) {
      mainWindow.minimize();
    } else {
      mainWindow.show();


      if (serverProcess == null) {
        if (mainWindow !== null && serverProcess == null) {

          // 打开 Team IDE 服务

          if (isDev) {
            const rootPath = getRootPath("../")
            log.info("root path:", rootPath)
            serverProcess = child_process.spawn(
              "go",
              ["run", ".", "--isDev", "--isElectron"],
              {
                cwd: getRootPath("../"),
              },
            );
          } else {
            let exePath = getRootPath('teamide-windows-x64.exe')
            try {
              fs.statSync(exePath);
            } catch (error) {

              try {
                exePath = getRootPath('teamide-darwin-x64')
                fs.statSync(exePath);
              } catch (error) {
                try {
                  exePath = getRootPath('teamide-linux-x64')
                  fs.statSync(exePath);
                } catch (error) {
                  exePath = "";
                }
              }
            }
            if (exePath == "") {
              // alert("Team IDE Server not found.")
              log.error("Team IDE Server not found.")
              if (app != null) {
                app.quit();
              }
              return
            }
            log.info("exePath:" + exePath)
            serverProcess = child_process.spawn(
              exePath,
              ["--isElectron"],
              {
                cwd: getRootPath(""),
              },
            );
          }
          serverProcess.stdout.on('data', (data: any) => {
            if (data == null) {
              return
            }
            let msg = data.toString()
            if (msg.startsWith("TeamIDE:event:serverUrl:")) {
              serverUrl = msg.substring("TeamIDE:event:serverUrl:".length)
              if (mainWindow != null) {
                mainWindow.loadURL(serverUrl);
              }
              return
            }
            log.info("msg:", msg);
          });
          serverProcess.stderr.on('data', (data: any) => {
            console.log('错误输出:');
            console.log(data.toString());
            console.log('--------------------');
          });
          serverProcess.on('close', (code: any) => {
            serverProcess = null;
            log.info(`server process close: ${code}`);
            destroyAll()
          });
        }
      }
    }
  });
  mainWindow.on('close', (e) => {
    e.preventDefault();
    // mainWindow = null;
    allWindowHide()
  });

  const menuBuilder = new MenuBuilder(mainWindow);
  menuBuilder.buildMenu();

  // Open urls in the user's browser
  mainWindow.webContents.setWindowOpenHandler((edata) => {
    shell.openExternal(edata.url);
    return { action: 'deny' };
  });

  // Remove this if your app does not use auto updates
  // eslint-disable-next-line
  new AppUpdater();
};

/**
 * Add event listeners...
 */

app.on('window-all-closed', () => {
  // Respect the OSX convention of having the application in memory even
  // after all windows have been closed
  if (process.platform !== 'darwin') {
    // app.quit();

    // log.info(`window all closed`);
    // if (serverProcess != null) {
    //   serverProcess.kill();
    // }
  }
});

let isAllWindowHide = false;
let allWindowShow = async () => {
  isAllWindowHide = false;
  if (mainWindow != null && !mainWindow.isDestroyed()) {
    mainWindow.show();
  }
  viewWindowList.forEach((one: BrowserWindow) => {
    if (!one.isDestroyed()) {
      one.show();
    }
  })
};
let allWindowHide = () => {
  isAllWindowHide = true;
  if (mainWindow != null && !mainWindow.isDestroyed()) {
    mainWindow.hide();
  }
  viewWindowList.forEach((one: BrowserWindow) => {
    if (!one.isDestroyed()) {
      one.hide();
    }
  })

};
let allWindowDestroy = () => {
  if (mainWindow != null && !mainWindow.isDestroyed()) {
    mainWindow.destroy();
    mainWindow = null;
  }
  viewWindowList.forEach((one: BrowserWindow) => {
    if (!one.isDestroyed()) {
      one.destroy();
    }
  })
  viewWindowList.splice(0, viewWindowList.length)
};

let destroyAll = () => {
  try {
    allWindowDestroy()
  } catch (error) {

  }
  try {
    if (serverProcess != null) {
      serverProcess.kill();
    }
  } catch (error) {

  }
  try {
    if (tray != null) {
      tray.destroy()
    }
  } catch (error) {

  }
  try {
    if (app != null) {
      app.quit()
    }
  } catch (error) {

  }
}
app
  .whenReady()
  .then(() => {
    createWindow();
    app.on('activate', () => {
      // On macOS it's common to re-create a window in the app when the
      // dock icon is clicked and there are no other windows open.
      if (mainWindow === null) createWindow();
    });
  })
  .catch(console.log);

let tray: Tray | null = null;
app.on('ready', async () => {
  tray = new Tray(icon16Path)
  const contextMenu = Menu.buildFromTemplate([
    {
      label: '退出',
      click: function () {
        destroyAll()
      }
    }
  ])
  tray.setToolTip('Team · IDE')
  //显示程序页面
  tray.on('click', () => {
    if (isAllWindowHide) {
      allWindowShow();
    } else {
      allWindowHide();
    }
  })
  tray.setContextMenu(contextMenu)
})