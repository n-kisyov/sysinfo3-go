package output

import (
	"fmt"
	"html/template"
	"io"
	"sysinfo3-go/pkg/collector"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>System Snapshot - {{ .Host.Hostname }}</title>
    <style>
        body { font-family: -apple-system, sans-serif; background-color: #0f111a; color: #f3f4f6; margin: 0; padding: 20px; }
        .container { max-width: 1000px; margin: 0 auto; }
        h1 { color: #8b5cf6; text-align: center; }
        h2 { color: #06b6d4; border-bottom: 1px solid #333; padding-bottom: 5px; margin-top: 30px; }
        .card { background: rgba(255,255,255,0.05); padding: 20px; border-radius: 8px; margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { text-align: left; padding: 8px; border-bottom: 1px solid #333; }
        th { color: #9ca3af; text-transform: uppercase; font-size: 0.85em; }
        .row { display: flex; justify-content: space-between; padding: 4px 0; }
        .row .label { color: #9ca3af; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Sysinfo3 Report</h1>
        <p style="text-align: center; color: #9ca3af;">Generated at: {{ .Timestamp.Format "2006-01-02 15:04:05" }}</p>

        <div class="card">
            <h2>System Overview</h2>
            <div class="row"><span class="label">Hostname</span><span>{{ .Host.Hostname }}</span></div>
            <div class="row"><span class="label">OS</span><span>{{ .Host.OS }} ({{ .Host.Platform }})</span></div>
            <div class="row"><span class="label">Uptime</span><span>{{ .Host.UptimeDisplay }}</span></div>
            <div class="row"><span class="label">CPU</span><span>{{ .CPU.Model }} ({{ .CPU.PhysicalCores }} Cores)</span></div>
            <div class="row"><span class="label">Memory</span><span>{{ .Memory.TotalH }}</span></div>
        </div>

        <div class="card">
            <h2>Network Interfaces</h2>
            <table>
                <tr><th>Name</th><th>Sent</th><th>Received</th></tr>
                {{ range .Network }}
                <tr>
                    <td>{{ .Name }}</td>
                    <td>{{ .BytesSentH }}</td>
                    <td>{{ .BytesRecvH }}</td>
                </tr>
                {{ end }}
            </table>
        </div>

        <div class="card">
            <h2>Disks (Volumes)</h2>
            <table>
                <tr><th>Mount</th><th>FS</th><th>Usage</th><th>Total</th></tr>
                {{ range .Disks }}
                <tr>
                    <td>{{ .MountPoint }}</td>
                    <td>{{ .FSType }}</td>
                    <td>{{ printf "%.1f" .UsedPercent }}%</td>
                    <td>{{ .TotalH }}</td>
                </tr>
                {{ end }}
            </table>
        </div>
        
        <div class="card">
            <h2>Top Processes</h2>
            <table>
                <tr><th>PID</th><th>Name</th><th>CPU %</th><th>Memory</th></tr>
                {{ range .Processes }}
                <tr>
                    <td>{{ .PID }}</td>
                    <td>{{ .Name }}</td>
                    <td>{{ printf "%.1f" .CPUPercent }}%</td>
                    <td>{{ .MemoryH }}</td>
                </tr>
                {{ end }}
            </table>
        </div>
    </div>
</body>
</html>`

func RenderHTML(s *collector.SystemSnapshot, w io.Writer) {
	t, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("Error creating HTML template: %v\n", err)
		return
	}
	err = t.Execute(w, s)
	if err != nil {
		fmt.Printf("Error executing HTML template: %v\n", err)
	}
}
