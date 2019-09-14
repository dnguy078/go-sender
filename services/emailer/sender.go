package emailer

type Emailer interface {
	Send() error
}
