import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import { Greet, GetVersion, CheckForUpdate, ApplyUpdate, GetPlatformInfo } from '../wailsjs/go/main/App';

// 应用 HTML 结构
document.querySelector('#app').innerHTML = `
    <div class="container">
        <!-- 主标题区域 -->
        <header class="header">
            <img id="logo" class="logo" alt="Wails Logo">
            <h1 class="title">Hello World</h1>
            <p class="subtitle">欢迎使用 Wails 演示应用</p>
        </header>

        <!-- 问候交互区域 -->
        <section class="greet-section">
            <div class="result" id="result">请在下方输入您的姓名 👇</div>
            <div class="input-box">
                <input class="input" id="name" type="text" placeholder="输入您的姓名..." autocomplete="off" />
                <button class="btn btn-primary" id="greet-btn">打招呼</button>
            </div>
        </section>

        <!-- 更新区域 -->
        <section class="update-section">
            <div class="update-status" id="update-status">
                <span class="status-text" id="status-text">正在检查更新...</span>
            </div>
            <div class="update-actions">
                <button class="btn btn-secondary" id="check-update-btn">检查更新</button>
                <button class="btn btn-success hidden" id="apply-update-btn">下载更新</button>
            </div>
            <div class="progress-container hidden" id="progress-container">
                <div class="progress-bar" id="progress-bar"></div>
                <span class="progress-text" id="progress-text">0%</span>
            </div>
        </section>

        <!-- 页脚版本信息 -->
        <footer class="footer">
            <span id="version-info">版本: 加载中...</span>
            <span class="separator">|</span>
            <span id="platform-info">平台: 加载中...</span>
        </footer>
    </div>
`;

// 设置 Logo
document.getElementById('logo').src = logo;

// 获取 DOM 元素
const nameElement = document.getElementById('name');
const resultElement = document.getElementById('result');
const greetBtn = document.getElementById('greet-btn');
const versionInfo = document.getElementById('version-info');
const platformInfo = document.getElementById('platform-info');
const statusText = document.getElementById('status-text');
const checkUpdateBtn = document.getElementById('check-update-btn');
const applyUpdateBtn = document.getElementById('apply-update-btn');
const progressContainer = document.getElementById('progress-container');
const progressBar = document.getElementById('progress-bar');
const progressText = document.getElementById('progress-text');

// 存储最新版本信息
let latestUpdateInfo = null;

// 初始化
async function init() {
    nameElement.focus();
    
    // 加载版本信息
    try {
        const version = await GetVersion();
        versionInfo.textContent = `版本: ${version}`;
    } catch (err) {
        console.error('获取版本失败:', err);
        versionInfo.textContent = '版本: 未知';
    }
    
    // 加载平台信息
    try {
        const platform = await GetPlatformInfo();
        platformInfo.textContent = `平台: ${platform}`;
    } catch (err) {
        console.error('获取平台信息失败:', err);
        platformInfo.textContent = '平台: 未知';
    }
    
    // 自动检查更新
    await checkUpdate();
}

// 问候功能
async function greet() {
    const name = nameElement.value.trim();
    if (!name) {
        resultElement.textContent = '请输入您的姓名！';
        resultElement.classList.add('error');
        setTimeout(() => resultElement.classList.remove('error'), 2000);
        return;
    }
    
    try {
        greetBtn.disabled = true;
        greetBtn.textContent = '处理中...';
        const result = await Greet(name);
        resultElement.textContent = result;
        resultElement.classList.add('success');
        setTimeout(() => resultElement.classList.remove('success'), 2000);
    } catch (err) {
        console.error('问候失败:', err);
        resultElement.textContent = '出错了，请稍后重试';
        resultElement.classList.add('error');
    } finally {
        greetBtn.disabled = false;
        greetBtn.textContent = '打招呼';
    }
}

// 检查更新
async function checkUpdate() {
    checkUpdateBtn.disabled = true;
    checkUpdateBtn.textContent = '检查中...';
    statusText.textContent = '正在检查更新...';
    statusText.className = 'status-text checking';
    
    try {
        const info = await CheckForUpdate();
        latestUpdateInfo = info;
        
        if (info.available) {
            statusText.textContent = `发现新版本: ${info.latestVersion}`;
            statusText.className = 'status-text available';
            applyUpdateBtn.classList.remove('hidden');
        } else {
            statusText.textContent = `已是最新版本 (${info.currentVersion})`;
            statusText.className = 'status-text latest';
            applyUpdateBtn.classList.add('hidden');
        }
    } catch (err) {
        console.error('检查更新失败:', err);
        statusText.textContent = '检查更新失败，请稍后重试';
        statusText.className = 'status-text error';
    } finally {
        checkUpdateBtn.disabled = false;
        checkUpdateBtn.textContent = '检查更新';
    }
}

// 应用更新
async function applyUpdate() {
    if (!latestUpdateInfo || !latestUpdateInfo.available) {
        statusText.textContent = '没有可用更新';
        return;
    }
    
    applyUpdateBtn.disabled = true;
    applyUpdateBtn.textContent = '下载中...';
    progressContainer.classList.remove('hidden');
    progressBar.style.width = '0%';
    progressText.textContent = '0%';
    statusText.textContent = '正在下载更新...';
    statusText.className = 'status-text downloading';
    
    try {
        // 模拟进度更新
        let progress = 0;
        const progressInterval = setInterval(() => {
            progress += Math.random() * 15;
            if (progress > 90) progress = 90;
            progressBar.style.width = `${progress}%`;
            progressText.textContent = `${Math.round(progress)}%`;
        }, 300);
        
        const result = await ApplyUpdate();
        
        clearInterval(progressInterval);
        progressBar.style.width = '100%';
        progressText.textContent = '100%';
        
        if (result.needRestart) {
            statusText.textContent = result.message;
            statusText.className = 'status-text ready';
            applyUpdateBtn.textContent = '重启应用';
            applyUpdateBtn.disabled = false;
            applyUpdateBtn.onclick = () => {
                // 提示用户手动重启
                alert('请关闭应用后重新打开以完成更新');
            };
        }
    } catch (err) {
        console.error('应用更新失败:', err);
        statusText.textContent = '更新失败，请稍后重试';
        statusText.className = 'status-text error';
        applyUpdateBtn.disabled = false;
        applyUpdateBtn.textContent = '重试下载';
        progressContainer.classList.add('hidden');
    }
}

// 事件绑定
greetBtn.addEventListener('click', greet);
nameElement.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') greet();
});
checkUpdateBtn.addEventListener('click', checkUpdate);
applyUpdateBtn.addEventListener('click', applyUpdate);

// 启动初始化
init();
