package output

import (
	"fmt"
	"os"
	"sysinfo3-go/internal/collector"

	"gopkg.in/yaml.v3"
)

func RenderYAML(s *collector.SystemSnapshot) {
	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	enc.SetIndent(2)
	if err := enc.Encode(s); err != nil {
		fmt.Fprintf(os.Stderr, "Error: YAML encoding failed: %v\n", err)
		os.Exit(1)
	}
}
