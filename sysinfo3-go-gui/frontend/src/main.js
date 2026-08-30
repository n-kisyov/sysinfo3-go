import './style.css';
import {GetSystemSnapshot} from '../wailsjs/go/main/App';

function renderRow(label, value) {
    return `
        <div class="row">
            <span class="label">${label}</span>
            <span class="value">${value}</span>
        </div>
    `;
}

function renderProgress(percent) {
    return `
        <div class="progress-container">
            <div class="progress-bar" style="width: ${percent}%"></div>
        </div>
    `;
}

async function updateData() {
    try {
        const snap = await GetSystemSnapshot();
        
        // Update Time
        document.getElementById('last-updated').innerText = `Last updated: ${new Date(snap.timestamp).toLocaleTimeString()}`;

        // Host
        document.getElementById('host-content').innerHTML = `
            ${renderRow('Hostname', snap.host.hostname)}
            ${renderRow('OS', snap.host.os)}
            ${renderRow('Uptime', snap.host.uptime)}
        `;

        // CPU
        document.getElementById('cpu-content').innerHTML = `
            ${renderRow('Model', snap.cpu.model)}
            ${renderRow('Cores', `${snap.cpu.physical_cores} / ${snap.cpu.logical_cores}`)}
            ${renderRow('Usage', `${snap.cpu.usage_percent.toFixed(1)}%`)}
            ${renderProgress(snap.cpu.usage_percent)}
        `;

        // Memory
        document.getElementById('memory-content').innerHTML = `
            ${renderRow('Total', snap.memory.total)}
            ${renderRow('Used', `${snap.memory.used} (${snap.memory.used_percent.toFixed(1)}%)`)}
            ${renderProgress(snap.memory.used_percent)}
            ${snap.memory.swap_total_bytes > 0 ? `
                ${renderRow('Swap Used', `${snap.memory.swap_used} (${snap.memory.swap_used_percent.toFixed(1)}%)`)}
                ${renderProgress(snap.memory.swap_used_percent)}
            ` : ''}
        `;

        // Battery
        if (snap.battery) {
            document.getElementById('battery-content').innerHTML = `
                ${renderRow('Status', snap.battery.status)}
                ${renderRow('Charge', `${snap.battery.percent}%`)}
                ${renderProgress(snap.battery.percent)}
                ${renderRow('Runtime', snap.battery.runtime)}
            `;
        } else {
            document.getElementById('battery-content').innerHTML = 'No battery detected.';
        }

        // GPU
        if (snap.gpu && snap.gpu.length > 0) {
            document.getElementById('gpu-content').innerHTML = snap.gpu.map(g => renderRow('GPU', g.name)).join('');
        } else {
            document.getElementById('gpu-content').innerHTML = 'No GPU detected.';
        }

        // Network
        document.getElementById('network-content').innerHTML = snap.network.filter(n => n.addresses.length > 0).map(n => `
            <div style="margin-bottom: 8px;">
                <strong>${n.name}</strong>
                ${renderRow('Sent', n.bytes_sent_h)}
                ${renderRow('Received', n.bytes_recv_h)}
            </div>
        `).join('');

        // Disks
        document.getElementById('disks-content').innerHTML = `
            <table>
                <tr><th>Mount</th><th>FS</th><th>Usage</th><th>Free</th><th>Total</th></tr>
                ${snap.disks.map(d => `
                    <tr>
                        <td>${d.mount_point}</td>
                        <td>${d.fstype}</td>
                        <td style="width:100px;">
                            <div style="font-size:0.8rem;text-align:right;">${d.used_percent.toFixed(1)}%</div>
                            ${renderProgress(d.used_percent)}
                        </td>
                        <td>${d.free}</td>
                        <td>${d.total}</td>
                    </tr>
                `).join('')}
            </table>
        `;

        // Processes
        document.getElementById('processes-content').innerHTML = `
            <table>
                <tr><th>PID</th><th>Name</th><th>CPU %</th><th>Memory</th></tr>
                ${snap.processes.map(p => `
                    <tr>
                        <td>${p.pid}</td>
                        <td>${p.name}</td>
                        <td>${p.cpu_percent.toFixed(1)}%</td>
                        <td>${p.memory}</td>
                    </tr>
                `).join('')}
            </table>
        `;

    } catch (err) {
        console.error(err);
    }
}

// Initial fetch and set interval
updateData();
setInterval(updateData, 2000);
