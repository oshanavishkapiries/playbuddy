// main.js
const { app, BrowserWindow, session, ipcMain } = require('electron');
const path = require('path');

let mainWindow;
let targetWindow;
let playerWindow;

function createMainWindow() {
  mainWindow = new BrowserWindow({
    width: 1000,
    height: 800,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    }
  });

  mainWindow.loadFile('index.html');
}

function createTargetWindow(url) {
  targetWindow = new BrowserWindow({
    width: 800,
    height: 600,
    show: false, // Hide the window
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true
    }
  });

  targetWindow.loadURL(url);

  // Monitor network requests
  session.defaultSession.webRequest.onCompleted({ urls: ['<all_urls>'] }, (details) => {
    const url = details.url;
    // Check for common video formats and streaming protocols
    if (url.match(/\.(mp4|m3u8|m3u|ts|mpd)(\?|$)/i) ||
      url.includes('video') ||
      url.includes('stream') ||
      url.includes('playlist')) {
      mainWindow.webContents.send('stream-detected', url);
    }
    mainWindow.webContents.send('network-request', {
      url: url,
      method: details.method,
      statusCode: details.statusCode
    });
  });

  // Close the hidden window after 30 seconds if no stream is found
  setTimeout(() => {
    if (targetWindow && !targetWindow.isDestroyed()) {
      targetWindow.close();
    }
  }, 30000);
}

function createPlayerWindow(url) {
  if (playerWindow) {
    playerWindow.close();
  }

  playerWindow = new BrowserWindow({
    width: 1280,
    height: 720,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false
    }
  });

  playerWindow.loadFile('player.html');

  playerWindow.webContents.on('did-finish-load', () => {
    playerWindow.webContents.send('play-video', url);
  });

  playerWindow.on('closed', () => {
    playerWindow = null;
  });
}

app.whenReady().then(createMainWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createMainWindow();
  }
});

ipcMain.on('load-url', (event, url) => {
  if (targetWindow && !targetWindow.isDestroyed()) {
    targetWindow.close();
  }
  createTargetWindow(url);
});

ipcMain.on('open-player', (event, url) => {
  createPlayerWindow(url);
});
