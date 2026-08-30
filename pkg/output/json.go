package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sysinfo3-go/pkg/collector"
)

func RenderJSON(s *collector.SystemSnapshot, w io.Writer) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		fmt.Fprintf(os.Stderr, "Error: JSON encoding failed: %v\n", err)
		os.Exit(1)
	}
}
