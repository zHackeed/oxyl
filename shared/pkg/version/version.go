package version

import _ "embed"

//go:generate bash -c "printf %s $(git rev-parse --short HEAD | head -c8) > .commit_id"
//go:embed .commit_id
var CommitID string
