package texttranslator

import (
	translator "github.com/Conight/go-googletrans"
	gt "github.com/bas24/googletranslatefree"
)

const (
	EN string = "en"
	ID string = "id"

	V1 string = "TranslateFree"
	V2 string = "TranslateConight"
)

type Translator struct {
	Version string
}

func New(version int) (*Translator, error) {
	var transVersion string
	if version == 1 {
		transVersion = V1
	} else {
		transVersion = V2
	}
	return &Translator{Version: transVersion}, nil
}

func (t *Translator) Translate(lang string, message string) string {

	if t.Version == V1 {
		result, err := TranslateFree(lang, message)
		if err != nil {
			return message
		}
		return result
	} else {
		result, err := TranslateConight(lang, message)
		if err != nil {
			return message
		}
		return result

	}
}

func TranslateFree(lang string, message string) (string, error) {
	// https://github.com/bas24/googletranslatefree
	if lang == ID {
		result, err := gt.Translate(message, "en", "id")
		if err != nil {
			return "", err
		}
		return result, nil
	} else {
		return message, nil
	}
}

func TranslateConight(lang string, message string) (string, error) {
	// https://github.com/Conight/go-googletrans
	if lang == ID {
		t := translator.New()
		result, err := t.Translate(message, "auto", "id")
		if err != nil {
			panic(err)
		}

		return result.Text, nil
	} else {
		return message, nil
	}
}
