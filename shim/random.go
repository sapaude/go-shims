package shim

import (
    "math/rand"
    "time"
)

// RandElem 从给定的切片中随机选择一个元素
func RandElem[T any](list []T) T {
    if len(list) == 0 {
        var zero T
        return zero // 返回类型 T 的零值
    }
    index := rand.Intn(len(list)) // 生成一个随机索引
    return list[index]            // 返回随机选择的元素
}

// Shuffle 打乱一个切片并返回一个新的打乱后的切片。
// 它接受一个泛型切片作为输入，并返回一个相同类型的泛型切片，原始切片不会被修改。
func Shuffle[T any](input []T) []T {
    // 1. 创建一个输入切片的副本，以确保不修改原始切片。
    shuffled := make([]T, len(input))
    copy(shuffled, input)

    // 2. 初始化一个随机数生成器。
    // 使用当前时间的纳秒作为种子，以确保每次运行都能得到不同的结果。
    // 注意：在 Go 1.20+ 中，rand.NewSource 推荐使用 rand.New(rand.NewSource(seed))
    // 在 Go 1.22+ 中，rand.New(rand.NewSource(seed)) 已经被 rand.New(rand.NewPCG(seed1, seed2)) 或 rand.New(rand.NewEntropy()) 替代，
    // 但对于大多数简单用例，time.Now().UnixNano() 仍然足够。
    // 为了兼容性，我们使用 rand.New(rand.NewSource(...))
    r := rand.New(rand.NewSource(time.Now().UnixNano()))

    // 3. 实现 Fisher-Yates 洗牌算法。
    // 从最后一个元素开始，向前遍历到第二个元素。
    for i := len(shuffled) - 1; i > 0; i-- {
        // 随机选择一个索引 j，范围从 0 到 i（包括 i）。
        j := r.Intn(i + 1)

        // 交换当前元素 shuffled[i] 和随机选择的元素 shuffled[j]。
        shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
    }

    return shuffled
}
