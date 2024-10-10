package arr

import "github.com/kbrownehs18/goutil/random"

// RandArray random string slice
func RandArray(arr []string) string {
	return arr[random.NewRand().Intn(len(arr))]
}

// RangeArray generate array
func RangeArray(m, n int) (b []int) {
	if m >= n || m < 0 {
		return b
	}

	c := make([]int, 0, n-m)
	for i := m; i < n; i++ {
		c = append(c, i)
	}

	return c
}
