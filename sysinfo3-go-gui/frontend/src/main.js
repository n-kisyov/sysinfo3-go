import './style.css';
import {GetSystemSnapshot, KillProcess} from '../wailsjs/go/main/App';

// Theme Toggle
const themeToggleBtn = document.getElementById('theme-toggle');
themeToggleBtn.addEventListener('click', () => {
    const currentTheme = document.documentElement.getAttribute('data-theme') || 'dark';
    document.documentElement.setAttribute('data-theme', currentTheme === 'dark' ? 'light' : 'dark');
});

// Tab Switching
document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', (e) => {
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.tab-pane').forEach(p => p.classList.remove('active'));
        
        e.target.classList.add('active');
        const targetId = e.target.getAttribute('data-target');
        document.getElementById(targetId).classList.add('active');
    });
});

// Chart.js Setup
let cpuChart, memChart;

function initCharts() {
    Chart.defaults.color = '#9ca3af';
    Chart.defaults.font.family = 'Outfit';

    const commonOptions = {
        responsive: true,
        maintainAspectRatio: false,
        animation: false,
        scales: {
            x: { display: false },
            y: { min: 0, max: 100, grid: { color: 'rgba(156, 163, 175, 0.1)' } }
        },
        plugins: { legend: { display: false } }
    };

    const ctxCpu = document.getElementById('cpuChart').getContext('2d');
    cpuChart = new Chart(ctxCpu, {
        type: 'line',
        data: { labels: [], datasets: [{ data: [], borderColor: '#06b6d4', backgroundColor: 'rgba(6, 182, 212, 0.1)', fill: true, tension: 0.4 }] },
        options: commonOptions
    });

    const ctxMem = document.getElementById('memChart').getContext('2d');
    memChart = new Chart(ctxMem, {
        type: 'line',
        data: { labels: [], datasets: [{ data: [], borderColor: '#8b5cf6', backgroundColor: 'rgba(139, 92, 246, 0.1)', fill: true, tension: 0.4 }] },
        options: commonOptions
    });
}

function updateChart(chart, value) {
    const data = chart.data.datasets[0].data;
    const labels = chart.data.labels;
    if (data.length > 30) {
        data.shift();
        labels.shift();
    }
    data.push(value);
    labels.push('');
    chart.update();
}

// Helpers
function renderRow(label, value) {
    return `<div class="row"><span class="label">${label}</span><span class="value">${value}</span></div>`;
}

function renderProgress(percent) {
    return `<div class="progress-container"><div class="progress-bar" style="width: ${percent}%"></div></div>`;
}

// Update Loop
async function updateData() {
    try {
        const snap = await GetSystemSnapshot();
        const t = new Date(snap.timestamp);
        document.getElementById('last-updated').innerText = `Last updated: ${t.toLocaleTimeString()}`;

        // Overview / System
        document.getElementById('host-content').innerHTML = `
            ${renderRow('Hostname', snap.host.hostname)}
            ${renderRow('OS', snap.host.os)}
            ${renderRow('Platform', snap.host.platform)}
            ${renderRow('Uptime', snap.host.uptime)}
        `;

        // Battery
        if (snap.battery) {
            document.getElementById('battery-content').innerHTML = `
                ${renderRow('Status', snap.battery.status)}
                ${renderRow('Charge', `${snap.battery.percentage}%`)}
                ${renderProgress(snap.battery.percentage)}
                ${renderRow('Runtime', snap.battery.time_left)}
            `;
        } else {
            document.getElementById('battery-content').innerHTML = '<span class="text-muted">No battery detected.</span>';
        }

        // Charts
        updateChart(cpuChart, snap.cpu.usage_percent);
        updateChart(memChart, snap.memory.used_percent);

        // Hardware Tab
        document.getElementById('bios-content').innerHTML = `
            ${renderRow('Vendor', snap.bios.vendor)}
            ${renderRow('Version', snap.bios.version)}
            ${renderRow('Date', snap.bios.date)}
        `;

        document.getElementById('cpu-content').innerHTML = `
            ${renderRow('Model', snap.cpu.model)}
            ${renderRow('Cores', `${snap.cpu.physical_cores} / ${snap.cpu.logical_cores}`)}
            ${renderRow('Usage', `${snap.cpu.usage_percent.toFixed(1)}%`)}
            ${renderProgress(snap.cpu.usage_percent)}
        `;

        document.getElementById('memory-content').innerHTML = `
            ${renderRow('Total', snap.memory.total)}
            ${renderRow('Used', `${snap.memory.used} (${snap.memory.used_percent.toFixed(1)}%)`)}
            ${renderProgress(snap.memory.used_percent)}
        `;

        if (snap.gpu && snap.gpu.length > 0) {
            document.getElementById('gpu-content').innerHTML = snap.gpu.map(g => renderRow('GPU', g.name)).join('');
        } else {
            document.getElementById('gpu-content').innerHTML = '<span class="text-muted">No GPU detected.</span>';
        }

        // Storage & Network
        document.getElementById('network-content').innerHTML = snap.network.filter(n => n.addresses.length > 0).map(n => `
            <div style="margin-bottom: 12px; padding-bottom: 8px; border-bottom: 1px solid var(--glass-border);">
                <div style="font-weight: 600; margin-bottom: 8px; color: var(--accent-color);">${n.name}</div>
                ${renderRow('Total Sent', n.bytes_sent_h)}
                ${renderRow('Total Recv', n.bytes_recv_h)}
                ${renderRow('Upload', n.bytes_sent_per_sec_h || '0 B/s')}
                ${renderRow('Download', n.bytes_recv_per_sec_h || '0 B/s')}
            </div>
        `).join('');

        document.getElementById('phys-disks-content').innerHTML = `
            <table>
                <tr><th>Model</th><th>Size</th><th>Read</th><th>Write</th></tr>
                ${snap.physical_disks.map(d => `
                    <tr>
                        <td style="font-size: 0.85em;">${d.model}</td>
                        <td>${d.size}</td>
                        <td>${d.read_bytes_per_sec_h || '0 B/s'}</td>
                        <td>${d.write_bytes_per_sec_h || '0 B/s'}</td>
                    </tr>
                `).join('')}
            </table>
        `;

        document.getElementById('disks-content').innerHTML = `
            <table>
                <tr><th>Mount</th><th>FS</th><th>Usage</th><th>Free</th><th>Total</th></tr>
                ${snap.disks.map(d => `
                    <tr>
                        <td>${d.mount_point}</td>
                        <td>${d.fs_type}</td>
                        <td style="width:120px;">
                            <div style="font-size:0.8rem;text-align:right;">${d.used_percent.toFixed(1)}%</div>
                            ${renderProgress(d.used_percent)}
                        </td>
                        <td>${d.free}</td>
                        <td>${d.total}</td>
                    </tr>
                `).join('')}
            </table>
        `;

        // Processes Tab
        document.getElementById('processes-content').innerHTML = `
            <table>
                <tr><th>PID</th><th>Name</th><th>CPU %</th><th>Memory</th><th>Action</th></tr>
                ${snap.processes.map(p => `
                    <tr>
                        <td>${p.pid}</td>
                        <td>${p.name}</td>
                        <td>${p.cpu_percent.toFixed(1)}%</td>
                        <td>${p.memory}</td>
                        <td><button class="kill-btn" data-pid="${p.pid}">Kill</button></td>
                    </tr>
                `).join('')}
            </table>
        `;

    } catch (err) {
        console.error("Update Error:", err);
    }
}

// Event delegation for kill buttons
document.addEventListener('click', async (e) => {
    if (e.target && e.target.classList.contains('kill-btn')) {
        const pid = parseInt(e.target.getAttribute('data-pid'), 10);
        if (confirm(`Are you sure you want to kill process ${pid}?`)) {
            try {
                await KillProcess(pid);
                updateData(); // Refresh immediately
            } catch (err) {
                alert(`Failed to kill process: ${err}`);
            }
        }
    }
});

// Init
initCharts();
updateData();
setInterval(updateData, 2000);
