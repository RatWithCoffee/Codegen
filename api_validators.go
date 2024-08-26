package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Tag int

//* `required` - поле является обязательным
//* `paramname` - если указано - то берется из параметра с этим именем, иначе `lowercase` от имени
//* `enum` - "одно из"
//* `default` - если указано и приходит пустое значение (значение по-умолчанию) - устанавливать то что написано указано
//* `minTag` - >= X для типа `int`, для строк `len(str)` >=
//* `max` - <= X для типа `int`

const (
	required Tag = iota
	paramname
	enum
	defaultTag
	minTag
	maxTag
)

func (t Tag) String() string {
	return [...]string{"required", "paramname", "enum", "defaultTag", "min", "max"}[t]
}

// проверка тегов для численного поля структуры
func IsTagIntValid(tag map[string]string, param string, paramName string) (int, error) {
	num, err := strconv.Atoi(param)
	if err != nil {
		return 0, errors.New(paramName + " must be int")
	}
	for name, val := range tag {
		border, _ := strconv.Atoi(val)
		if name == minTag.String() {
			if num < border {
				return 0, fmt.Errorf("%s must be >= %s", paramName, val)
			}
		} else {
			if num > border {
				return 0, fmt.Errorf("%s must be <= %s", paramName, val)
			}
		}

	}
	return num, nil
}

// проверка тегов для строкового поля структуры
func IsTagStrValid(tag map[string]string, param string, paramName string) (string, error) {
	resVal := param

	val, ok := tag[defaultTag.String()]
	if ok && param == "" {
		resVal = val
	}

	_, ok = tag[required.String()]
	if ok && param == "" {
		return param, errors.New(paramName + " must be not empty")
	}

	for name, val := range tag {
		switch name {
		case minTag.String():
			border, _ := strconv.Atoi(val)
			if len(resVal) < border {
				return resVal, fmt.Errorf("%s len must be >= %s", paramName, val)
			}
		case maxTag.String():
			border, _ := strconv.Atoi(val)
			if len(resVal) > border {
				return resVal, fmt.Errorf("%s len must be <= %s", paramName, val)
			}
		case enum.String():
			values := strings.Split(val, "|")
			isValid := false
			for _, v := range values {
				if resVal == v {
					isValid = true
				}
			}
			if !isValid {
				arr := ""
				for _, v := range strings.Split(val, "|") {
					arr += v + ", "
				}
				arr = arr[:len(arr)-2]
				return resVal, fmt.Errorf("status must be one of [%s]", arr)
			}
		}
	}
	return resVal, nil
}
