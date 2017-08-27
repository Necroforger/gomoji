package gomoji

import (
	"bytes"
)

// Generate emoji constants.
//go:generate go run gen/gen.go

// var emojiRegex = regexp.MustCompile(":[a-zA-Z0-9_]+:")

// FormatRegex replaces indexes of :emojiname: with the corresponding emoji
// func FormatRegex(text string) string {
// 	for _, emojiText := range emojiRegex.FindAllString(text, -1) {
// 		if em := Emoji(emojiText[1 : len(emojiText)-1]); em != "" {
// 			text = strings.Replace(text, emojiText, em, 1)
// 		}
// 	}
// 	return text
// }

// Format replaces indexes of :emojiname: with the corresponding emoji
func Format(text string) string {
	return format(text)
}

func format(text string) string {
	input := bytes.NewBufferString(text)
	var output bytes.Buffer

	for r, _, err := input.ReadRune(); err == nil; r, _, err = input.ReadRune() {
		switch r {
		case ':':
			output.Write([]byte(replace(input)))
		default:
			output.WriteRune(r)
		}
	}

	return string(output.Bytes())
}

func replace(input *bytes.Buffer) string {
	var emojiName bytes.Buffer

	var count int
	for r, _, err := input.ReadRune(); err == nil; r, _, err = input.ReadRune() {
		switch r {
		case ':':
			if count == 0 {
				return ":" + emojiName.String() + ":"
			}
			if e := Emoji(emojiName.String()); e != "" {
				return e
			}
			emojiName.WriteRune(r)
			return ":" + emojiName.String()

		case ' ':
			emojiName.WriteRune(r)
			return ":" + emojiName.String()

		default:
			emojiName.WriteRune(r)
		}
		count++
	}
	return ":" + emojiName.String()
}
