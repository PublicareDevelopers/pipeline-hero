package notifier

type Client interface {
	BuildBlocks() error
	Notify() error
}
