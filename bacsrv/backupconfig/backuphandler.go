package backupconfig

// Backup specifies a backup
type Backup struct {
	Paths     []string `json:"paths"`
	Excluded  []string `json:"excludedPaths"`
	Retention string   `json:"retentionTime"`
}

type BackupTriggerMessage struct {
	ClientName   string `json:"clientName"`
	BackupConfig Backup
}
