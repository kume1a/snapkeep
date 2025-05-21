package backup

import (
	"os/exec"
)

func DumpDatabase(dbURL, outFile string) error {
	cmd := exec.Command("pg_dump", dbURL, "-f", outFile)
	return cmd.Run()
}
