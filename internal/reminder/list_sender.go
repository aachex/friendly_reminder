package reminder

type ListSender interface {
	Send(email string) error
}

type listSender struct{}

func (ss *listSender) Send(email string) error {
	// По почте получить список из базы данных и отправить его на эту же почту
	return nil
}
