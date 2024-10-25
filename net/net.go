package net

import (
	"fmt"
	"net"
)

// IsValidIP ip address is valid
func IsValidIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil
}

// IsValidNetmask
func IsValidNetmask(maskStr string) bool {
	// 解析掩码字符串
	mask := net.ParseIP(maskStr).To4()
	if mask == nil {
		return false // 无效的 IPv4 掩码
	}

	// 将掩码转换为 IPMask 并获取其前缀长度
	maskSize, bits := net.IPv4Mask(mask[0], mask[1], mask[2], mask[3]).Size()

	// 判断是否为合法的子网掩码
	return maskSize >= 0 && bits == 32
}

// ConvertMask 子网掩码转换
func BitsToMask(prefixLength int) (string, error) {
	if prefixLength < 0 || prefixLength > 32 {
		return "", fmt.Errorf("invalid prefix length: %d", prefixLength)
	}

	// 使用 net 包中的 IPv4Mask 函数根据前缀生成子网掩码
	mask := net.CIDRMask(prefixLength, 32)
	ipMask := net.IP(mask)

	return ipMask.String(), nil
}

// MaskToBits 将子网掩码字符串转换为掩码位数
func MaskToBits(maskStr string) (int, error) {
	// 解析掩码字符串为 IP 格式
	mask := net.ParseIP(maskStr).To4()
	if mask == nil {
		return 0, fmt.Errorf("invalid IPv4 mask")
	}

	// 将 IP 转为掩码，并获取其位数
	ones, bits := net.IPv4Mask(mask[0], mask[1], mask[2], mask[3]).Size()
	if bits != 32 {
		return 0, fmt.Errorf("invalid IPv4 mask length")
	}

	return ones, nil
}
