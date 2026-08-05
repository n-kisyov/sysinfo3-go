package collector

import "testing"

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{512, "512 B"},
		{1023, "1023 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
		{1099511627776, "1.00 TB"},
		{2147483648, "2.00 GB"},       // 2 GB
		{8589934592, "8.00 GB"},       // 8 GB - typical VRAM
		{17179869184, "16.00 GB"},     // 16 GB
	}

	for _, tc := range tests {
		result := FormatBytes(tc.input)
		if result != tc.expected {
			t.Errorf("FormatBytes(%d) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}

func TestFormatUptime(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{0, "0m"},
		{59, "0m"},
		{60, "1m"},
		{3599, "59m"},
		{3600, "1h 0m"},
		{3660, "1h 1m"},
		{86400, "1d 0h 0m"},
		{90061, "1d 1h 1m"},
		{172800, "2d 0h 0m"},
		{259200, "3d 0h 0m"},
	}

	for _, tc := range tests {
		result := FormatUptime(tc.input)
		if result != tc.expected {
			t.Errorf("FormatUptime(%d) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}
