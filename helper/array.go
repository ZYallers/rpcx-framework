package helper

//  RemoveDuplicateWithInt 去除重复值
//  @author Cloud|2021-12-12 18:20:52
//  @param arr []int ...
//  @return []int ...
func RemoveDuplicateWithInt(arr []int) []int {
	var result []int      // 存放返回的不重复切片
	tmp := map[int]byte{} // 存放不重复主键
	for _, val := range arr {
		l := len(tmp)
		tmp[val] = 0 // 当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		// 如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tmp) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, val) // 当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

//  RemoveDuplicateWithString 去除重复值
//  @author Cloud|2021-12-13 09:33:05
//  @param arr []string ...
//  @return []string ...
func RemoveDuplicateWithString(arr []string) []string {
	var result []string      // 存放返回的不重复切片
	tmp := map[string]byte{} // 存放不重复主键
	for _, val := range arr {
		l := len(tmp)
		tmp[val] = 0 // 当e存在于tempMap中时，再次添加是添加不进去的，，因为key不允许重复
		// 如果上一行添加成功，那么长度发生变化且此时元素一定不重复
		if len(tmp) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, val) // 当元素不重复时，将元素添加到切片result中
		}
	}
	return result
}

//  RemoveWithString 删除数组中的指定值
//  @author Cloud|2021-12-12 18:20:21
//  @param arr []string ...
//  @param in string ...
//  @return []string ...
func RemoveWithString(arr []string, in string) []string {
	for k, v := range arr {
		if v == in {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	return arr
}

//  RemoveWithInt 删除数组中的指定值
//  @author Cloud|2021-12-14 12:15:15
//  @param arr []int ...
//  @param in int ...
//  @return []int ...
func RemoveWithInt(arr []int, in int) []int {
	for k, v := range arr {
		if v == in {
			arr = append(arr[:k], arr[k+1:]...)
			break
		}
	}
	return arr
}
