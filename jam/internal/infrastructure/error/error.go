package error

type RespnoseError struct {
	err error
}

func (re *RespnoseError) Error() string {
	return re.err.Error()
}
