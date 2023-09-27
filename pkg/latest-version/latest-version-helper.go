package latestVersion

import (
	_ "embed"
)

//go:embed latest-version.txt
var LatestVersion string
