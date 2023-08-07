package misc

func MergeMaps(original map[string]any, overwriteWith map[string]any) map[string]any {
	merged := make(map[string]any)
	for k, v := range original {
		merged[k] = v
	}
	for key, value := range overwriteWith {
		merged[key] = value
	}
	return merged
}
