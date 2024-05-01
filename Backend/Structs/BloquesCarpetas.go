package Structs

type BloqueCarpetas struct {
	B_content [4]Content
}

func NewBloquesCarpetas() BloqueCarpetas {
	var bl BloqueCarpetas
	for i := 0; i < len(bl.B_content); i++ {
		bl.B_content[i] = NewContent()
	}
	return bl
}
