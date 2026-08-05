package output

import (
	"os"
	"sysinfo3-go/internal/collector"

	"gopkg.in/yaml.v3"
)

func RenderYAML(s *collector.SystemSnapshot) {
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	_ = enc.Encode(s)
}
