package wotwhb

type UniqueStringList struct {
	contents []string
	tracking map[string]bool
}

func NewUniqueStringList(elements ...string) *UniqueStringList {
	uniqueList := &UniqueStringList{[]string{}, map[string]bool{}}
	return uniqueList.Add(elements...)
}

func (u *UniqueStringList) Add(elements ...string) *UniqueStringList {
	for _, element := range elements {
		if "" == element {
			continue
		}
		if _, ok := u.tracking[element]; !ok {
			u.tracking[element] = true
			u.contents = append(u.contents, element)
		}
	}
	return u
}

func (u *UniqueStringList) Contents() []string {
	return u.contents
}

func (u *UniqueStringList) Size() int {
	return len(u.contents)
}
