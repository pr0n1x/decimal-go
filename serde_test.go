package dec

import (
	"encoding/json"
	"errors"
	"testing"
)

var testJson = `{"a": "0.1","b": "0.2"}`

type twoDecValuesJson struct {
	A *TextNano `json:"a"`
	B *TextNano `json:"b"`
}

func Test_Encode_TypePrecisionMismatch(t *testing.T) {
	type decTxt[T DecimalTextTrait] struct {
		Value *T `json:"value"`
	}
	cases := []any{
		decTxt[TextDeci]{Value: TextDeciRef(FromUInt64(1, Centi))},
		decTxt[TextCenti]{Value: TextCentiRef(FromUInt64(1, Deci))},
		decTxt[TextMilli]{Value: TextMilliRef(FromUInt64(1, Deci))},
		decTxt[TextMicro]{Value: TextMicroRef(FromUInt64(1, Deci))},
		decTxt[TextNano]{Value: TextNanoRef(FromUInt64(1, Deci))},
		decTxt[TextPico]{Value: TextPicoRef(FromUInt64(1, Deci))},
		decTxt[TextFemto]{Value: TextFemtoRef(FromUInt64(1, Deci))},
		decTxt[TextAtto]{Value: TextAttoRef(FromUInt64(1, Deci))},
		decTxt[TextZepto]{Value: TextZeptoRef(FromUInt64(1, Deci))},
		decTxt[TextYocto]{Value: TextYoctoRef(FromUInt64(1, Deci))},
		decTxt[TextRonto]{Value: TextRontoRef(FromUInt64(1, Deci))},
		decTxt[TextQuecto]{Value: TextQuectoRef(FromUInt64(1, Deci))},
	}
	for _, testCase := range cases {
		if _, err := json.Marshal(testCase); err == nil || !errors.Is(err, ErrTypePrecisionMismatch) {
			t.Fatal("expected ErrTypePrecisionMismatch")
		}
	}
}

func Test_Decode(t *testing.T) {
	var values twoDecValuesJson
	if err := json.Unmarshal([]byte(testJson), &values); err != nil {
		t.Error(err)
	}

	if values.A.GetDecimal().String() != "0.1" {
		t.Error("unmarshal failed")
	}

	if values.B.GetDecimal().String() != "0.2" {
		t.Error("unmarshal failed")
	}
}
