package googlefree

import (
	"errors"
	"math/rand"
	"time"

	translator "github.com/zijiren233/google-translator"

	"github.com/metatube-community/metatube-sdk-go/translate"
)

func Translate(q, source, target string) (string, error) {
	data, err := translator.Translate(q, target, translator.TranslationParams{
		From:       source,
		Retry:      2,
		RetryDelay: 3 * time.Second,
	})
	if err != nil /* fallback */ {
		if data, err = translator.TranslateWithClienID(q, target, translator.TranslationWithClienIDParams{
			From:       source,
			Retry:      2,
			ClientID:   rand.Intn(5) + 1,
			RetryDelay: 3 * time.Second,
		}); err != nil {
			return "", err
		}
	}
	if data == nil {
		return "", errors.New("data is nil")
	}
	return data.Text, nil
}

func init() {
	translate.Register("googlefree",
		func(text, from, to string, _ struct{}) (string, error) {
			return Translate(text, from, to)
		},
		func() struct{} {
			return struct{}{}
		},
	)
}