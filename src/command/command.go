package command

type GitCommand interface {
    Init(args []string) error

    Add(args []string) error
    Status(args []string) error
    Commit(args []string) error

    Branch(args []string) error
    Checkout(args []string) error

    Reset(args []string) error
    CherryPick(args []string) error

    Log(args []string) error
    Diff(args []string) error
}
