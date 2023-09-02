package cli

type Cli struct {
	Clean            clean            `cmd:"clean" aliases:"c" help:"Remove previous database file."`
	PrintFranchises  printLocations   `cmd:"print" aliases:"p" help:"Print saved locations (requires previous run)."`
	SearchFranchises searchFranchises `cmd:"franchises" aliases:"f,franchises" help:"Search franchises."`
}
