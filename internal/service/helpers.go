package service

func AddCategory(arr []int) []int {
	for i, nb := range arr {
		arr[i] = nb + 1
	}
	return arr
}
