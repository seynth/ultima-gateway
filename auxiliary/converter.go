package auxiliary

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"ultima/model"

	tea "github.com/charmbracelet/bubbletea"
	"howett.net/plist"
)

func Convert(seb string) tea.Cmd {
	return func() tea.Msg {
		file, err := os.Open(seb)
		if err != nil {
			return model.ReadAndConvert{Err: err}
		}
		defer file.Close()

		var data map[string]interface{}
		decoder := plist.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil {
			fmt.Println("Error decode: ", err)
		}
		delete(data, "originatorVersion")

		jsonStr := GenerateSortedJSON(data)
		return model.ReadAndConvert{EncodedContent: jsonStr}
	}
}

func GenerateSortedJSON(data interface{}) string {
	switch v := data.(type) {
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
		})

		var sb strings.Builder
		sb.WriteString("{")
		for i, key := range keys {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(`"` + key + `":`)
			sb.WriteString(GenerateSortedJSON(v[key]))
		}
		sb.WriteString("}")
		return sb.String()

	case []interface{}:
		var sb strings.Builder
		sb.WriteString("[")
		for i, item := range v {
			if i > 0 {
				sb.WriteString(",")
			}
			sb.WriteString(GenerateSortedJSON(item))
		}
		sb.WriteString("]")
		return sb.String()

	case string:
		return `"` + v + `"`
	case float64:
		if v == float64(int64(v)) {
			return strconv.FormatInt(int64(v), 10)
		}
		return strconv.FormatFloat(v, 'f', -1, 64)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	default:
		switch v := data.(type) {
		case uint:
			return strconv.FormatUint(uint64(v), 10)
		case uint64:
			return strconv.FormatUint(v, 10)
		default:
			return `""`
		}
	}
}
