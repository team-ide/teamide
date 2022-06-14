
import { BrowserView } from 'electron';
import { source } from './util'

const windowList: any = []

export const createWindow = (config: any) => {

    var view = new BrowserView()   //new出对象
    source.mainWindow.setBrowserView(view)   // 在主窗口中设置view可用
    view.setBounds({ x: 0, y: 100, width: 1200, height: 800 })  //定义view的具体样式和位置
    view.webContents.loadURL(config.url)  //wiew载入的页面

    windowList.push(view)
};