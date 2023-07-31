package utils

func ConvertUintSlice2Map(uintSlice []uint) map[uint]struct{} {

	set := make(map[uint]struct{}, len(uintSlice))

	for _, v := range uintSlice {
		set[v] = struct{}{}
	}
	return set

}
