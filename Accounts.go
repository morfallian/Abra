package Slices_and_maps

import (
	"sort"
)

type account struct{
	uname string
	money uint64
}

type accountList map[string]account

func (acc *accountList) AddNewAccount(login string, accountNew account) error {

	if _, ok := (*acc)[login]; ok{
		return errAccountAlreadyExists
	}

	if login == "" {
		return errEmptyAccount
	}

	(*acc)[login] = account{
		uname: accountNew.uname,
		money: accountNew.money,
	}

	return nil
}
type tempSortStruct struct {
	login string
	value account
}
//typeSort = 1 по алфавиту поле: Имя
//typeSort = 2 в обратном порядке поле: Имя
//typeSort = 3 вобратном порядке поле: Деньги
func (acc *accountList) SortAcc(typeSort int) []tempSortStruct{

	tempSlice := make([]tempSortStruct, 0)

	for ind, val := range *acc {
		tempSlice = append(tempSlice, tempSortStruct{
			login: ind,
			value: val,
		})
	}

	sort.Slice(tempSlice, func(i, j int) bool {
		return tempSlice[i].value.uname < tempSlice[j].value.uname
	})

	switch typeSort {
	case 1:
		sort.Slice(tempSlice, func(i, j int) bool {
			return tempSlice[i].value.uname < tempSlice[j].value.uname
		})
	case 2:
		sort.Slice(tempSlice, func(i, j int) bool {
			return tempSlice[i].value.uname > tempSlice[j].value.uname
		})
	case 3:
		sort.Slice(tempSlice, func(i, j int) bool {
			return tempSlice[i].value.money > tempSlice[j].value.money
		})
	}

	//for _, value := range tempSlice {
	//	fmt.Print(value.value.uname, ":", value.value.money, "\n")
	//}

	return tempSlice
}

