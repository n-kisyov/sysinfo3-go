package output

import (
	"encoding/json"
	"fmt"
	"os"
	"sysinfo3-go/internal/collector"
)

func RenderJSON(s *collector.SystemSnapshot) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		fmt.Fprintf(os.Stderr, "Error: JSON encoding failed: %v\n", err)
		os.Exit(1)
	}
}
