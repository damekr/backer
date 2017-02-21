package backupconfig

type Backup struct {
	ID        string
	Paths     []string
	Excluded  []string
	Retention string
}
