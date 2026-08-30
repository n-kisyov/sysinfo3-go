package output

import (
	"fmt"
	"io"
	"os"
	"sysinfo3-go/internal/collector"

	"gopkg.in/yaml.v3"
)

func RenderYAML(s *collector.SystemSnapshot, w io.Writer) {
	enc := yaml.NewEncoder(w)
	defer enc.Close()
	enc.SetIndent(2)
	if err := enc.Encode(s); err != nil {
		fmt.Fprintf(os.Stderr, "Error: YAML encoding failed: %v\n", err)
		os.Exit(1)
	}
}
