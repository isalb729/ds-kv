package utils

const (
	HoleNum = 137
)

func Label(labelList []int) ([]int, int) {
	length := len(labelList)
	if length == 0 {
		return []int{0}, 0
	} else if length == 1 {
		label := (HoleNum / 2 + labelList[0]) % HoleNum
		if label > labelList[0] {
			return []int{labelList[0], label}, label
		} else {
			return []int{label, labelList[0]}, label
		}
	}
	max := Dist(labelList[length - 1], labelList[0], HoleNum)
	i := 0
	for j := 1; j < len(labelList); j++ {
		if labelList[j] - labelList[j - 1] > max {
			i = j
			max = labelList[j] - labelList[j - 1]
		}
	}
	// max gap is from i-1 to i
	var label int
	if i == 0 {
		label = (max / 2 + labelList[length - 1]) % HoleNum
		labelList = append(labelList, label)
	} else {
		label = max / 2 + labelList[i - 1]
		labelList = Insert(labelList, i - 1, label)
	}
	return labelList, label
}


func Insert(list []int, index, val int) []int {
	last := len(list) - 1
	if last == -1 {
		return []int{val}
	}
	if index == last {
		return append(list, val)
	}
	list = append(list, list[last])
	copy(list[index + 2:last + 1], list[index + 1:last])
	list[index + 1] = val
	return list
}


func Dist(a, b, max int) int {
	var dist int
	if a <= b {
		dist = b - a
		if max - b + a < dist {
			dist = max - b + a
		}
	} else {
		dist = a - b
		if max - a + b < dist {
			dist = max - a + b
		}
	}
	return dist
}