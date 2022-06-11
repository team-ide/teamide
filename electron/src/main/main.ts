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
import { app, BrowserWindow, shell, ipcMain, Menu, Tray } from 'electron';
import { autoUpdater } from 'electron-updater';
import log from 'electron-log';
import MenuBuilder from './menu';
import { resolveHtmlPath } from './util';

var fs = require("fs")
const child_process = require('child_process');

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

if (process.env.NODE_ENV === 'production') {
  const sourceMapSupport = require('source-map-support');
  sourceMapSupport.install();
}

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

let iconPath = getAssetPath('icon.png');

const createWindow = async () => {
  if (isDebug) {
    await installExtensions();
  }


  mainWindow = new BrowserWindow({
    title: "Team · IDE",
    show: false,
    width: 1440,
    height: 900,
    icon: iconPath,
    autoHideMenuBar: true,
    webPreferences: {
      preload: app.isPackaged
        ? path.join(__dirname, 'preload.js')
        : path.join(__dirname, '../../.erb/dll/preload.js'),
    },
  });

  mainWindow.loadURL(resolveHtmlPath('index.html'));

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

          let exePath = getAssetPath('teamide-windows-x64.exe')
          try {
            fs.statSync(exePath);
          } catch (error) {

            try {
              exePath = getAssetPath('teamide-darwin-x64')
              fs.statSync(exePath);
            } catch (error) {
              try {
                exePath = getAssetPath('teamide-linux-x64')
                fs.statSync(exePath);
              } catch (error) {
                exePath = "";
              }
            }
            if (exePath != "") {
              try {
                child_process.spawn(
                  "chmod",
                  ["+x", exePath],
                  {
                    cwd: getAssetPath("")
                  },
                );
              } catch (error) {

              }
            }
          }
          if (exePath == "") {
            log.error("Team IDE Server not found.")
            if (app != null) {
              app.quit();
            }
            return
          }
          log.info("exePath:" + exePath)
          const path = exePath
          serverProcess = child_process.spawn(
            path,
            ["--isElectron"],
            {
              cwd: getAssetPath("")
            },
          );
          serverProcess.stdout.on('data', (data: any) => {
            if (data == null) {
              return
            }
            let msg = data.toString()
            if (msg.startsWith("TeamIDE:event:serverUrl:")) {
              let serverUrl = msg.substring("TeamIDE:event:serverUrl:".length)
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
            log.info(`server process close: ${code}`);
            if (app != null) {
              app.quit();
            }
          });
        }
      }
    }
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
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
    app.quit();
  }
  log.info(`window all closed`);
  if (serverProcess != null) {
    serverProcess.kill();
  }
});

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
let tray = null

app.on('ready', async () => {
  tray = new Tray(iconPath)
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
    if (mainWindow != null) {
      mainWindow.show();
    }
  })
  tray.setContextMenu(contextMenu)
})