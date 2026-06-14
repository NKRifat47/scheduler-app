// ==========================================================================
// CHRONOTASK CLIENT CONTROLLER
// ==========================================================================

let currentUser = null;
let tasks = [];
let eventSource = null;
let currentFilter = 'all';

// Element Cache
const screens = {
    loading: document.getElementById('loading-screen'),
    auth: document.getElementById('auth-screen'),
    dashboard: document.getElementById('dashboard-screen')
};

const authTabs = {
    login: document.getElementById('tab-login'),
    signup: document.getElementById('tab-signup'),
    loginForm: document.getElementById('login-form'),
    signupForm: document.getElementById('signup-form')
};

// ==========================================================================
// DOM INITIALIZATION
// ==========================================================================

document.addEventListener('DOMContentLoaded', () => {
    initAuthTabs();
    initForms();
    initFilterTabs();
    initModal();
    showScreen('auth');
    
    // Set default value for datetime-local to be current time + 1 minute
    const datetimeInput = document.getElementById('task-datetime');
    const now = new Date();
    now.setMinutes(now.getMinutes() + 1);
    // Format to YYYY-MM-DDThh:mm
    const tzoffset = now.getTimezoneOffset() * 60000; //offset in milliseconds
    const localISOTime = (new Date(now - tzoffset)).toISOString().slice(0, 16);
    datetimeInput.value = localISOTime;
    datetimeInput.min = (new Date(new Date() - tzoffset)).toISOString().slice(0, 16);

    // Live counter for current time hint
    setInterval(updateTimeHint, 1000);
    updateTimeHint();

    // Request Notification permission early
    if ('Notification' in window && Notification.permission === 'default') {
        Notification.requestPermission();
    }
});

function updateTimeHint() {
    const hint = document.getElementById('current-time-hint');
    if (hint) {
        const now = new Date();
        hint.textContent = `Current time: ${now.toLocaleDateString()} ${now.toLocaleTimeString()}`;
    }
}

// ==========================================================================
// SCREEN UTILITIES
// ==========================================================================

function showScreen(screenKey) {
    Object.keys(screens).forEach(key => {
        if (key === screenKey) {
            screens[key].style.display = 'flex';
            // Force reflow
            screens[key].offsetHeight;
            screens[key].classList.add('active');
        } else {
            screens[key].classList.remove('active');
            screens[key].style.display = 'none';
        }
    });
}

// ==========================================================================
// AUTHENTICATION FLOW
// ==========================================================================

function initAuthTabs() {
    authTabs.login.addEventListener('click', () => {
        authTabs.login.classList.add('active');
        authTabs.signup.classList.remove('active');
        authTabs.loginForm.classList.add('active');
        authTabs.signupForm.classList.remove('active');
    });

    authTabs.signup.addEventListener('click', () => {
        authTabs.signup.classList.add('active');
        authTabs.login.classList.remove('active');
        authTabs.signupForm.classList.add('active');
        authTabs.loginForm.classList.remove('active');
    });
}

function initForms() {
    // Login
    authTabs.loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const usernameInput = document.getElementById('login-username');
        const passwordInput = document.getElementById('login-password');
        
        try {
            const res = await fetch('/api/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    username: usernameInput.value,
                    password: passwordInput.value
                })
            });
            
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || 'Login failed');
            
            showToast('Welcome back!', 'success');
            usernameInput.value = '';
            passwordInput.value = '';
            
            await showDashboard(data.username);
        } catch (err) {
            showToast(err.message, 'error');
        }
    });

    // Signup
    authTabs.signupForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const usernameInput = document.getElementById('signup-username');
        const passwordInput = document.getElementById('signup-password');
        
        try {
            const res = await fetch('/api/signup', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    username: usernameInput.value,
                    password: passwordInput.value
                })
            });
            
            const data = await res.json();
            if (!res.ok) throw new Error(data.error || 'Registration failed');
            
            showToast('Account created successfully!', 'success');
            usernameInput.value = '';
            passwordInput.value = '';
            
            await showDashboard(data.username);
        } catch (err) {
            showToast(err.message, 'error');
        }
    });

    // Logout
    document.getElementById('logout-btn').addEventListener('click', async () => {
        try {
            await fetch('/api/logout', { method: 'POST' });
            currentUser = null;
            closeSSEConnection();
            showScreen('auth');
            showToast('Logged out successfully', 'info');
        } catch (err) {
            showToast('Failed to logout', 'error');
        }
    });

    // Create Task
    document.getElementById('task-form').addEventListener('submit', async (e) => {
        e.preventDefault();
        const titleInput = document.getElementById('task-title');
        const descInput = document.getElementById('task-desc');
        const datetimeInput = document.getElementById('task-datetime');
        
        try {
            // Parse local time to ISO representation
            const localDate = new Date(datetimeInput.value);
            const res = await fetch('/api/tasks', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    title: titleInput.value,
                    description: descInput.value,
                    scheduled_time: localDate.toISOString()
                })
            });

            const data = await res.json();
            if (!res.ok) throw new Error(data.error || 'Failed to create task');

            showToast('Task scheduled successfully', 'success');
            titleInput.value = '';
            descInput.value = '';
            
            // Set default +1m for next task
            const nextTime = new Date();
            nextTime.setMinutes(nextTime.getMinutes() + 1);
            const tzoffset = nextTime.getTimezoneOffset() * 60000;
            datetimeInput.value = (new Date(nextTime - tzoffset)).toISOString().slice(0, 16);

            await fetchTasks();
        } catch (err) {
            showToast(err.message, 'error');
        }
    });
}

async function showDashboard(username) {
    currentUser = username;
    document.getElementById('username-display').textContent = currentUser;
    showScreen('dashboard');
    await fetchTasks();
    setupSSEConnection();
}

// ==========================================================================
// TASKS OPERATIONS
// ==========================================================================

async function fetchTasks() {
    try {
        const res = await fetch('/api/tasks');
        if (!res.ok) throw new Error('Could not fetch tasks');
        tasks = await res.json();
        renderTasks();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

function initFilterTabs() {
    const filterButtons = document.querySelectorAll('.filter-btn');
    filterButtons.forEach(btn => {
        btn.addEventListener('click', (e) => {
            filterButtons.forEach(b => b.classList.remove('active'));
            e.target.classList.add('active');
            currentFilter = e.target.dataset.filter;
            renderTasks();
        });
    });
}

function renderTasks() {
    const container = document.getElementById('tasks-container');
    container.innerHTML = '';

    const filtered = tasks.filter(task => {
        if (currentFilter === 'pending') return !task.triggered;
        if (currentFilter === 'triggered') return task.triggered;
        return true;
    });

    if (filtered.length === 0) {
        let msg = "No tasks scheduled yet.";
        if (currentFilter === 'pending') msg = "No pending tasks scheduled.";
        if (currentFilter === 'triggered') msg = "No tasks have triggered yet.";
        
        container.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">📂</div>
                <p>${msg}</p>
            </div>
        `;
        return;
    }

    filtered.forEach(task => {
        const card = document.createElement('div');
        card.className = `task-card ${task.triggered ? 'triggered' : 'pending'}`;
        
        const scheduledDate = new Date(task.scheduled_time);
        const formattedTime = scheduledDate.toLocaleString();
        
        card.innerHTML = `
            <div class="task-info">
                <div class="task-title-row">
                    <span class="task-title">${escapeHTML(task.title)}</span>
                    <span class="task-badge">${task.triggered ? 'Triggered' : 'Pending'}</span>
                </div>
                ${task.description ? `<p class="task-desc-text">${escapeHTML(task.description)}</p>` : ''}
                <div class="task-meta">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
                    </svg>
                    <span>${formattedTime}</span>
                </div>
            </div>
            <button class="delete-task-btn" data-id="${task.id}" title="Delete Task">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
                </svg>
            </button>
        `;

        // Bind delete action
        card.querySelector('.delete-task-btn').addEventListener('click', async (e) => {
            e.stopPropagation();
            const taskId = e.currentTarget.dataset.id;
            if (confirm('Are you sure you want to cancel and delete this task?')) {
                await deleteTask(taskId);
            }
        });

        container.appendChild(card);
    });
}

async function deleteTask(id) {
    try {
        const res = await fetch(`/api/tasks/${id}`, { method: 'DELETE' });
        if (!res.ok) {
            const data = await res.json();
            throw new Error(data.error || 'Failed to delete task');
        }
        showToast('Task removed', 'info');
        await fetchTasks();
    } catch (err) {
        showToast(err.message, 'error');
    }
}

// ==========================================================================
// SERVER-SENT EVENTS (REAL-TIME NOTIFICATION BROADCAST)
// ==========================================================================

function setupSSEConnection() {
    closeSSEConnection(); // Close existing just in case

    const statusEl = document.getElementById('stream-status');
    const dot = statusEl.querySelector('.status-dot');
    const text = statusEl.querySelector('.status-text');

    eventSource = new EventSource('/api/events');

    eventSource.onopen = () => {
        dot.className = 'status-dot connected';
        text.textContent = 'Real-time feed active';
    };

    eventSource.onerror = (e) => {
        dot.className = 'status-dot disconnected';
        text.textContent = 'Disconnected. Retrying...';
        // SSE handles reconnection automatically, but we can monitor it
    };

    eventSource.addEventListener('task_triggered', async (event) => {
        try {
            const task = JSON.parse(event.data);
            triggerNotificationAlert(task);
            await fetchTasks(); // Update list to show marked triggered
        } catch (err) {
            console.error('Error handling task_triggered event:', err);
        }
    });
}

function closeSSEConnection() {
    if (eventSource) {
        eventSource.close();
        eventSource = null;
    }
    const statusEl = document.getElementById('stream-status');
    if (statusEl) {
        statusEl.querySelector('.status-dot').className = 'status-dot disconnected';
        statusEl.querySelector('.status-text').textContent = 'Disconnected';
    }
}

// ==========================================================================
// NOTIFICATION DISPATCH (MODAL, OS PUSH & SOUND)
// ==========================================================================

function triggerNotificationAlert(task) {
    // 1. Play Synthesized Chime sound
    playChime();

    // 2. Open Custom Glass Modal Popup
    openModal(task);

    // 3. Desktop Notification (if supported & permission granted)
    if ('Notification' in window && Notification.permission === 'granted') {
        const options = {
            body: task.description || 'ChronoTask Triggered!',
            icon: '/favicon.ico',
            tag: `chrono-task-${task.id}`,
            requireInteraction: true
        };
        new Notification(`⏰ CHRONOTASK: ${task.title}`, options);
    }
}

function playChime() {
    try {
        const AudioContextClass = window.AudioContext || window.webkitAudioContext;
        if (!AudioContextClass) return;
        
        const ctx = new AudioContextClass();
        
        // Tone generator helper
        const playTone = (freq, startTime, duration) => {
            const osc = ctx.createOscillator();
            const gain = ctx.createGain();
            
            osc.connect(gain);
            gain.connect(ctx.destination);
            
            // Futuristic sine/triangle blend
            osc.type = 'sine';
            osc.frequency.setValueAtTime(freq, startTime);
            
            gain.gain.setValueAtTime(0, startTime);
            // Quick attack
            gain.gain.linearRampToValueAtTime(0.12, startTime + 0.02);
            // Slow decay
            gain.gain.exponentialRampToValueAtTime(0.001, startTime + duration);
            
            osc.start(startTime);
            osc.stop(startTime + duration);
        };
        
        const now = ctx.currentTime;
        // G-major high chime arpeggio
        playTone(392.00, now, 0.4);        // G4
        playTone(493.88, now + 0.08, 0.4); // B4
        playTone(587.33, now + 0.16, 0.4); // D5
        playTone(783.99, now + 0.24, 0.6); // G5
    } catch (err) {
        console.warn('Web Audio Playback blocked or unsupported:', err);
    }
}

function initModal() {
    const modal = document.getElementById('notification-modal');
    const closeX = document.getElementById('modal-close-x');
    const ackBtn = document.getElementById('modal-ack-btn');

    const closeHandler = () => {
        modal.classList.remove('active');
    };

    closeX.addEventListener('click', closeHandler);
    ackBtn.addEventListener('click', closeHandler);
    
    // Close on clicking overlay outside card
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            closeHandler();
        }
    });
}

function openModal(task) {
    const modal = document.getElementById('notification-modal');
    document.getElementById('modal-task-title').textContent = task.title;
    document.getElementById('modal-task-desc').textContent = task.description || 'No description provided for this scheduled task.';
    
    const scheduledDate = new Date(task.scheduled_time);
    document.getElementById('modal-task-time').textContent = scheduledDate.toLocaleTimeString();

    modal.classList.add('active');
}

// ==========================================================================
// TOAST ALERT CONTROLLER
// ==========================================================================

function showToast(message, type = 'info') {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    toast.innerHTML = `
        <span>${escapeHTML(message)}</span>
        <button class="toast-close">&times;</button>
    `;
    
    const closeToast = () => {
        toast.style.animation = 'toastIn 0.25s reverse ease-out forwards';
        toast.addEventListener('animationend', () => toast.remove());
    };

    toast.querySelector('.toast-close').addEventListener('click', (e) => {
        e.stopPropagation();
        closeToast();
    });

    // Auto-remove after 4 seconds
    setTimeout(closeToast, 4000);

    container.appendChild(toast);
}

// ==========================================================================
// HELPERS
// ==========================================================================

function escapeHTML(str) {
    if (!str) return '';
    return str
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#039;');
}
