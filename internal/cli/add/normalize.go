package add

import "strings"

func normalizeTargetSection(section string) string {
	targetSection := section
	if targetSection == "" {
		return "## Next Up"
	}
	if !strings.HasPrefix(targetSection, "##") {
		return "## " + targetSection
	}
	return targetSection
}
