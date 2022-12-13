package build

import (
	_ "embed"
)

//go:generate bash get_version.sh
//go:embed version.txt
var Version string
var Commit string
