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
		decTxt[TextMilli]{Value: TextMilliRef(MustFromUInt64(1, Micro))},
		decTxt[TextMicro]{Value: TextMicroRef(MustFromUInt64(1, Milli))},
		decTxt[TextNano]{Value: TextNanoRef(MustFromUInt64(1, Milli))},
		decTxt[TextPico]{Value: TextPicoRef(MustFromUInt64(1, Milli))},
		decTxt[TextFemto]{Value: TextFemtoRef(MustFromUInt64(1, Milli))},
		decTxt[TextAtto]{Value: TextAttoRef(MustFromUInt64(1, Milli))},
		decTxt[TextZepto]{Value: TextZeptoRef(MustFromUInt64(1, Milli))},
		decTxt[TextYocto]{Value: TextYoctoRef(MustFromUInt64(1, Milli))},
		decTxt[TextRonto]{Value: TextRontoRef(MustFromUInt64(1, Milli))},
		decTxt[TextQuecto]{Value: TextQuectoRef(MustFromUInt64(1, Milli))},
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
