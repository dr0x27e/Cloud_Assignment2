package utils

// Function to convert Interface slice to slice of type T:
func ConvSliceInter[T any](input []interface{}) []T {
	output := make([]T, len(input))
	for i, v := range input {
		output[i] = v.(T)
	}
	return output
}

// Function to cenvert map[string]interface{} to map[string]T
func ConvMapInter[T any](input map[string]interface{}) map[string]T {
	output := make(map[string]T)
	for key, v := range input {
		output[key] = v.(T)
	}
	return output
}
