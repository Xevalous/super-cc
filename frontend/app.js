/**
 * Super CC - Main Controller
 */

// ============================================
// Utilities
// ============================================
function escapeHTML(str) {
    if (typeof str !== 'string') return String(str);
    const map = { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#039;' };
    return str.replace(/[&<>"']/g, c => map[c]);
}

// ============================================
// State
// ============================================
const state = {
    selectedVersion: null,
    allVersions: [],
    filteredVersions: [],
    currentFilter: 'all',
    searchQuery: '',
};

// ============================================
// Navigation
// ============================================
function navigateTo(pageId) {
    const navLinks = document.querySelectorAll('.nav-link');
    const pages = document.querySelectorAll('.page');

    navLinks.forEach(l => l.classList.remove('active'));
    const targetLink = document.querySelector(`.nav-link[data-page="${pageId}"]`);
    if (targetLink) targetLink.classList.add('active');

    pages.forEach(page => {
        page.classList.remove('active');
        if (page.id === `page-${pageId}`) {
            page.classList.add('active');
            switch(pageId) {
                case 'dashboard':
                    loadInstallationInfo();
                    break;
                case 'download':
                    loadVersions();
                    break;
                case 'lock':
                    loadLockStatus();
                    break;
                case 'crack':
                    loadCrackStatus();
                    break;
                case 'cleanup':
                    loadCleanupStatus();
                    break;
                case 'debug':
                    loadDebugInfo();
                    break;
            }
        }
    });
}

document.addEventListener('DOMContentLoaded', () => {
    const navLinks = document.querySelectorAll('.nav-link');

    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const pageId = link.getAttribute('data-page');
            navigateTo(pageId);
        });
    });

    // Modal handlers
    document.getElementById('modal-cancel')?.addEventListener('click', () => modal.hide(false));
    document.getElementById('modal-confirm')?.addEventListener('click', () => modal.hide(true));
    document.getElementById('modal-overlay')?.addEventListener('click', (e) => {
        if (e.target.id === 'modal-overlay') modal.hide(false);
    });

    // Close modal on Escape
    document.addEventListener('keydown', (e) => {
        if (e.key === 'Escape') {
            const overlay = document.getElementById('modal-overlay');
            if (overlay && overlay.style.display !== 'none') {
                modal.hide(false);
            }
        }
    });

    // Search & filter handlers
    setupSearchAndFilter();

    // Download action buttons
    document.getElementById('btn-download')?.addEventListener('click', downloadSelected);

    // Initial load
    if (window.go) {
        loadInstallationInfo();
    } else {
        console.log('Running in browser mode (not in Wails)');
        document.getElementById('app-status').textContent = 'Browser mode';
        document.getElementById('app-status-detail').textContent = 'Wails not available';
    }
});

// ============================================
// Modal
// ============================================
const modal = {
    overlay: null,
    resolvePromise: null,

    show({ title, message, confirmText = 'Confirm', cancelText = 'Cancel', danger = false, iconName = 'warning-circle' }) {
        return new Promise((resolve) => {
            if (this.resolvePromise) {
                resolve(false);
                return;
            }
            this.resolvePromise = resolve;
            this.overlay = document.getElementById('modal-overlay');

            const modalIcon = document.getElementById('modal-icon');
            const modalTitle = document.getElementById('modal-title');
            const modalMessage = document.getElementById('modal-message');
            const confirmBtn = document.getElementById('modal-confirm');
            const cancelBtn = document.getElementById('modal-cancel');

            modalIcon.innerHTML = `<i class="ph ph-${iconName}"></i>`;
            modalIcon.className = danger ? 'modal-icon danger' : 'modal-icon';
            modalTitle.textContent = title;
            modalMessage.textContent = message;
            confirmBtn.textContent = confirmText;
            cancelBtn.textContent = cancelText;

            this.overlay.style.display = 'flex';
        });
    },

    hide(result) {
        if (this.overlay) {
            this.overlay.style.display = 'none';
        }
        if (this.resolvePromise) {
            this.resolvePromise(result);
            this.resolvePromise = null;
        }
    }
};

// ============================================
// Dashboard
// ============================================
async function loadInstallationInfo() {
    try {
        const result = await window.go.main.App.GetInstallationInfo();
        const statusEl = document.getElementById('app-status');
        const statusDetail = document.getElementById('app-status-detail');
        const iconWrapper = document.getElementById('status-icon-wrapper');
        const iconMain = document.getElementById('status-icon-main');

        if (result.status === 'Installed') {
            iconMain.className = 'ph ph-check-circle';
            iconWrapper.className = 'status-icon-wrapper installed';
            statusEl.textContent = 'Installed';
            statusDetail.textContent = 'Ready to use';
            document.getElementById('app-version').textContent = result.version || 'Unknown';
            document.getElementById('app-path').textContent = result.path;
            document.getElementById('app-size').textContent = result.size;
        } else {
            iconMain.className = 'ph ph-x-circle';
            iconWrapper.className = 'status-icon-wrapper not-installed';
            statusEl.textContent = 'Not Installed';
            statusDetail.textContent = 'Not found on this system';
            document.getElementById('app-version').textContent = '-';
            document.getElementById('app-path').textContent = '-';
            document.getElementById('app-size').textContent = '-';
        }
    } catch (error) {
        console.error('Failed to load installation info:', error);
        document.getElementById('app-status').textContent = 'Error';
        document.getElementById('app-status-detail').textContent = error.message;
    }
}

// ============================================
// Download - Version Data
// ============================================
const VERSION_DATABASE = [
    // Latest
    { label: "5.4.0 (Beta6)", version: "5.4.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1991_beta6_capcutpc_beta_creatortool.exe", type: "latest", tag: "Latest", risk: "High" },
    { label: "5.4.0 (Beta5)", version: "5.4.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1988_beta5_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.4.0 (Beta4)", version: "5.4.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1982_beta4_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.4.0 (Beta3)", version: "5.4.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1979_beta3_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.4.0 (Beta2)", version: "5.4.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1978_beta2_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.4.0 (Beta1)", version: "5.4.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_4_0_1976_beta1_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    // 5.3.0
    { label: "5.3.0 (Latest)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1964_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "High" },
    { label: "5.3.0 (Test2)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1961_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "High" },
    { label: "5.3.0 (Test1)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1957_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "High" },
    { label: "5.3.0 (Beta5)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1962_beta5_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.3.0 (Beta4)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1956_beta4_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.3.0 (Beta3)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1952_beta3_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.3.0 (Beta2 Latest)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1949_beta2_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.3.0 (Beta2 Test1)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1947_beta2_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.3.0 (Beta1 Latest)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1942_beta1_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    { label: "5.3.0 (Beta1 Test1)", version: "5.3.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_3_0_1941_beta1_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    // 5.2.0
    { label: "5.2.0 (Latest)", version: "5.2.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1950_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "High" },
    { label: "5.2.0 (Test3)", version: "5.2.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1946_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "High" },
    { label: "5.2.0 (Test2)", version: "5.2.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1940_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "High" },
    { label: "5.2.0 (Test1)", version: "5.2.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1939_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "High" },
    { label: "5.2.0 (Beta8)", version: "5.2.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_5_2_0_1945_beta8_capcutpc_beta_creatortool.exe", type: "beta", tag: "Beta", risk: "High" },
    // Older stable
    { label: "4.8.0", version: "4.8.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_4_8_0_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "Medium" },
    { label: "4.5.0", version: "4.5.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_4_5_0_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "Medium" },
    { label: "3.5.0", version: "3.5.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_3_5_0_capcutpc_0_creatortool.exe", type: "stable", tag: "Stable", risk: "Low" },
    // First version
    { label: "1.0.0 (Latest)", version: "1.0.0", url: "https://lf16-capcut.faceulv.com/obj/capcutpc-packages-us/packages/CapCut_1_0_0_44_capcutpc_0.exe", type: "stable", tag: "Stable", risk: "Low" },
];

function loadVersions() {
    state.allVersions = VERSION_DATABASE;
    state.filteredVersions = [...VERSION_DATABASE];
    state.selectedVersion = null;

    document.getElementById('btn-download').disabled = true;
    document.getElementById('version-details-panel').style.display = 'none';

    applyFilters();
}

function setupSearchAndFilter() {
    const searchInput = document.getElementById('version-search');
    if (searchInput) {
        let searchTimeout;
        searchInput.addEventListener('input', (e) => {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                state.searchQuery = e.target.value;
                applyFilters();
            }, 300);
        });
    }

    document.querySelectorAll('.filter-chip').forEach(chip => {
        chip.addEventListener('click', () => {
            document.querySelectorAll('.filter-chip').forEach(c => c.classList.remove('active'));
            chip.classList.add('active');
            state.currentFilter = chip.dataset.filter;
            applyFilters();
        });
    });
}

function applyFilters() {
    let filtered = [...state.allVersions];

    if (state.currentFilter !== 'all') {
        filtered = filtered.filter(v => v.type === state.currentFilter);
    }

    if (state.searchQuery) {
        const query = state.searchQuery.toLowerCase();
        filtered = filtered.filter(v =>
            v.label.toLowerCase().includes(query) ||
            v.version.toLowerCase().includes(query) ||
            v.tag.toLowerCase().includes(query)
        );
    }

    state.filteredVersions = filtered;
    renderVersionList();
}

function renderVersionList() {
    const container = document.getElementById('version-list');
    const countEl = document.getElementById('version-count');

    if (!container) return;

    countEl.textContent = `${state.filteredVersions.length} version${state.filteredVersions.length !== 1 ? 's' : ''} found`;

    if (state.filteredVersions.length === 0) {
        container.innerHTML = `
            <div class="version-empty">
                <i class="ph ph-magnifying-glass"></i>
                <span>No versions found</span>
            </div>
        `;
        return;
    }

    const fragment = document.createDocumentFragment();

    state.filteredVersions.forEach((v, i) => {
        const row = document.createElement('div');
        row.className = 'version-row';
        row.tabIndex = 0;
        row.dataset.index = i;

        const riskClass = v.risk === 'High' ? 'badge-risk-high' : v.risk === 'Medium' ? 'badge-risk-medium' : 'badge-risk-low';
        const iconBg = v.type === 'latest' ? 'var(--accent-green)' : v.type === 'beta' ? 'var(--accent-orange)' : 'var(--accent-indigo)';

        row.innerHTML = `
            <div class="version-icon" style="background: ${iconBg}">
                <i class="ph ph-package"></i>
            </div>
            <div class="version-content">
                <span class="version-title">${escapeHTML(v.label)}</span>
                <span class="version-subtitle">${escapeHTML(v.tag)} <span class="version-badge ${riskClass}">${escapeHTML(v.risk)} Risk</span></span>
            </div>
            <i class="ph ph-check version-check"></i>
        `;

        row.addEventListener('click', () => selectVersion(i));
        row.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                selectVersion(i);
            }
        });

        fragment.append(row);
    });

    container.innerHTML = '';
    container.append(fragment);
}

function selectVersion(idx) {
    state.selectedVersion = state.filteredVersions[idx];

    document.querySelectorAll('.version-row').forEach((row, i) => {
        row.classList.toggle('selected', i === idx);
    });

    const panel = document.getElementById('version-details-panel');
    panel.style.display = 'block';

    document.getElementById('detail-label').textContent = state.selectedVersion.label;
    document.getElementById('detail-version').textContent = state.selectedVersion.version;
    document.getElementById('detail-risk').textContent = state.selectedVersion.risk;
    document.getElementById('detail-type').textContent = state.selectedVersion.tag;

    document.getElementById('btn-download').disabled = false;
}

async function downloadSelected() {
    if (!state.selectedVersion) return;

    // Check if version is locked — download URLs are blocked when locked
    try {
        const locked = await window.go.main.App.IsLocked();
        if (locked) {
            await modal.show({
                title: 'Version Locked',
                message: 'Download URLs are blocked while the version is locked. Please unlock the version first from the Lock page before downloading.',
                confirmText: 'OK',
                cancelText: 'Close',
                danger: true,
                iconName: 'shield-warning'
            });
            return;
        }
    } catch (error) {
        // If we can't check lock status, proceed with download
    }

    const confirmed = await modal.show({
        title: 'Download Version?',
        message: `This will download ${state.selectedVersion.label}. The installer will open in your browser.`,
        confirmText: 'Download',
        cancelText: 'Cancel',
        danger: false,
        iconName: 'download-simple'
    });

    if (!confirmed) return;

    try {
        await window.go.main.App.OpenURL(state.selectedVersion.url);
    } catch (error) {
        await modal.show({
            title: 'Error',
            message: `Failed to open browser: ${error.message || error}`,
            confirmText: 'OK',
            cancelText: 'Close',
            danger: true,
            iconName: 'x-circle'
        });
    }
}

// ============================================
// Lock
// ============================================
async function loadLockStatus() {
    try {
        const locked = await window.go.main.App.IsLocked();
        const version = await window.go.main.App.GetLockedVersion();

        const statusEl = document.getElementById('lock-status');
        const detailEl = document.getElementById('lock-detail');
        const iconWrapper = document.getElementById('lock-icon-wrapper');
        const iconMain = document.getElementById('lock-icon-main');

        if (locked) {
            iconMain.className = 'ph ph-shield-check';
            iconWrapper.className = 'status-icon-wrapper locked';
            statusEl.textContent = 'Version Locked';
            detailEl.textContent = 'Updates are blocked via hosts file';
        } else {
            iconMain.className = 'ph ph-shield-warning';
            iconWrapper.className = 'status-icon-wrapper unlocked';
            statusEl.textContent = 'Version Unlocked';
            detailEl.textContent = 'Auto-updates are enabled';
        }

        document.getElementById('lock-version').textContent = version || '-';

        const toggleBtn = document.getElementById('lock-toggle-btn');
        if (toggleBtn) {
            if (locked) {
                toggleBtn.innerHTML = '<i class="ph ph-lock-open"></i> Unlock';
            } else {
                toggleBtn.innerHTML = '<i class="ph ph-lock"></i> Lock';
            }
        }
    } catch (error) {
        console.error('Failed to load lock status:', error);
        document.getElementById('lock-status').textContent = 'Error';
        document.getElementById('lock-detail').textContent = error.message || 'Could not check lock status';
        document.getElementById('lock-version').textContent = '-';
    }
}

async function toggleLock() {
    try {
        const locked = await window.go.main.App.IsLocked();
        const version = await window.go.main.App.GetLockedVersion();
        const enable = !locked;

        const confirmed = await modal.show({
            title: enable ? 'Lock Version?' : 'Unlock Version?',
            message: enable
                ? 'This requires Administrator privileges. Make sure you are running this app as Administrator.'
                : 'This requires Administrator privileges. Make sure you are running this app as Administrator.',
            confirmText: enable ? 'Lock' : 'Unlock',
            cancelText: 'Cancel',
            danger: !enable,
            iconName: enable ? 'lock' : 'lock-open'
        });

        if (!confirmed) return;

        await window.go.main.App.SetLock(enable, version);
        loadLockStatus();
    } catch (error) {
        await modal.show({
            title: 'Error',
            message: `Failed to toggle lock: ${error.message || error}`,
            confirmText: 'OK',
            cancelText: 'Close',
            danger: true,
            iconName: 'x-circle'
        });
    }
}

// ============================================
// Crack
// ============================================
async function loadCrackStatus() {
    try {
        const status = await window.go.main.App.GetCrackStatus();
        const mode = status.mode;
        const vipEnabled = status.vipEnabled;

        const statusEl = document.getElementById('crack-vip-status');
        const detailEl = document.getElementById('crack-mode');
        const iconWrapper = document.getElementById('crack-icon-wrapper');
        const iconMain = document.getElementById('crack-icon-main');

        if (vipEnabled) {
            iconMain.className = 'ph ph-check-circle';
            iconWrapper.className = 'status-icon-wrapper installed';
            statusEl.textContent = 'Unlocked';
            detailEl.textContent = mode ? `Mode: ${mode}` : '';
        } else {
            iconMain.className = 'ph ph-lock-simple';
            iconWrapper.className = 'status-icon-wrapper unlocked';
            statusEl.textContent = 'Locked';
            detailEl.textContent = 'Premium features are restricted';
        }

        const toggleBtn = document.getElementById('crack-toggle-btn');
        if (toggleBtn) {
            if (vipEnabled) {
                toggleBtn.innerHTML = '<i class="ph ph-arrow-counter-clockwise"></i> Revert Patch';
            } else {
                toggleBtn.innerHTML = '<i class="ph ph-magic-wand"></i> Apply Patch';
            }
        }
    } catch (error) {
        console.error('Failed to load crack status:', error);
        document.getElementById('crack-vip-status').textContent = 'Error';
        document.getElementById('crack-mode').textContent = error.message || 'Could not detect';
    }
}

async function toggleCrack() {
    try {
        const status = await window.go.main.App.GetCrackStatus();
        const vipEnabled = status.vipEnabled;
        const enable = !vipEnabled;

        const confirmed = await modal.show({
            title: enable ? 'Apply Patch?' : 'Revert Patch?',
            message: enable
                ? 'Make sure the application is running before applying the patch. This requires Administrator privileges.'
                : 'This will restore the original VECreator.dll. This requires Administrator privileges.',
            confirmText: enable ? 'Apply' : 'Revert',
            cancelText: 'Cancel',
            danger: enable,
            iconName: enable ? 'magic-wand' : 'arrow-counter-clockwise'
        });

        if (!confirmed) return;

        await window.go.main.App.ApplyCrack(enable);
        loadCrackStatus();
    } catch (error) {
        await modal.show({
            title: 'Error',
            message: `Failed to toggle patch: ${error.message || error}`,
            confirmText: 'OK',
            cancelText: 'Close',
            danger: true,
            iconName: 'x-circle'
        });
    }
}

// ============================================
// Debug
// ============================================
async function loadDebugInfo() {
    const logEl = document.getElementById('debug-log');
    if (!logEl) return;

    logEl.innerHTML = `
        <div class="debug-loading">
            <i class="ph ph-circle-notch spin"></i>
            <span>Loading debug info...</span>
        </div>
    `;

    try {
        const info = await window.go.main.App.GetDebugInfo();
        const keys = Object.keys(info);
        logEl.innerHTML = '';
        keys.forEach(key => {
            const row = document.createElement('div');
            row.className = 'debug-row';
            row.innerHTML = `
                <span class="debug-key">${escapeHTML(key)}</span>
                <span class="debug-value">${escapeHTML(info[key])}</span>
            `;
            logEl.appendChild(row);
        });
    } catch (error) {
        logEl.innerHTML = `
            <div class="debug-loading">
                <i class="ph ph-x-circle"></i>
                <span>Error: ${error.message || error}</span>
            </div>
        `;
    }
}

// ============================================
// Cleanup
// ============================================
async function loadCleanupStatus() {
    try {
        const result = await window.go.main.App.GetCleanupStatus();

        document.getElementById('cleanup-files-found').textContent = result.filesFound;
        document.getElementById('cleanup-files-cleaned').textContent = result.filesCleaned;

        const pathsEl = document.getElementById('cleanup-paths');
        if (pathsEl) {
            if (result.residualPaths && result.residualPaths.length > 0) {
                pathsEl.innerHTML = '';
                result.residualPaths.forEach(p => {
                    const li = document.createElement('li');
                    li.className = 'cleanup-path-item';
                    li.innerHTML = `
                        <span class="cleanup-path-text">${escapeHTML(p)}</span>
                        <button class="btn-icon" title="Open in Explorer">
                            <i class="ph ph-folder-open"></i>
                        </button>
                    `;
                    li.querySelector('.btn-icon').addEventListener('click', async () => {
                        try {
                            await window.go.main.App.OpenFolder(p);
                        } catch (err) {
                            console.error('Failed to open folder:', err);
                        }
                    });
                    pathsEl.appendChild(li);
                });
            } else {
                pathsEl.innerHTML = '<li class="cleanup-path-item empty">No residual paths found</li>';
            }
        }
    } catch (error) {
        console.error('Failed to load cleanup status:', error);
        document.getElementById('cleanup-files-found').textContent = '-';
        document.getElementById('cleanup-files-cleaned').textContent = '-';
        const pathsEl = document.getElementById('cleanup-paths');
        if (pathsEl) {
            pathsEl.innerHTML = `<li class="cleanup-path-item empty">Error: ${error.message || 'Could not scan'}</li>`;
        }
    }
}

async function runCleanup() {
    const confirmed = await modal.show({
        title: 'Run Cleanup?',
        message: 'This will remove all residual files and reset the configuration. This action cannot be undone.',
        confirmText: 'Cleanup',
        cancelText: 'Cancel',
        danger: true,
        iconName: 'broom'
    });

    if (!confirmed) return;

    try {
        const result = await window.go.main.App.RunCleanup();
        await modal.show({
            title: 'Cleanup Complete',
            message: `${result.filesCleaned} item${result.filesCleaned !== 1 ? 's' : ''} cleaned successfully.`,
            confirmText: 'OK',
            cancelText: 'Close',
            danger: false,
            iconName: 'check-circle'
        });
        loadCleanupStatus();
    } catch (error) {
        await modal.show({
            title: 'Error',
            message: `Failed to run cleanup: ${error.message || error}`,
            confirmText: 'OK',
            cancelText: 'Close',
            danger: true,
            iconName: 'x-circle'
        });
    }
}
