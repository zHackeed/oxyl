package version

import _ "embed"

//go:generate sh -c "printf %s $(git rev-parse --short HEAD | head -c8) > .commit_id"
//go:embed .commit_id
var CommitID string

//go:generate sh -c "printf %s $(git rev-parse --abbrev-ref HEAD) > .branch_id"
//go:embed .branch_id
var Branch string
