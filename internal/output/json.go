package output

import (
	"encoding/json"
	"os"
	"sysinfo3-go/internal/collector"
)

func RenderJSON(s *collector.SystemSnapshot) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(s)
}

func RenderJSONCompact(s *collector.SystemSnapshot) {
	enc := json.NewEncoder(os.Stdout)
	_ = enc.Encode(s)
}
