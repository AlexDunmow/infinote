package report

import "fmt"

type Console struct {
}

func (s *Console) LogExternal(err error) {
	fmt.Println(err)
}
