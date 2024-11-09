package dec

import (
	"errors"
	"fmt"
)

type (
	TextDeci   Decimal
	TextCenti  Decimal
	TextMilli  Decimal
	TextMicro  Decimal
	TextNano   Decimal
	TextPico   Decimal
	TextFemto  Decimal
	TextAtto   Decimal
	TextZepto  Decimal
	TextYocto  Decimal
	TextRonto  Decimal
	TextQuecto Decimal
)

var ErrTypePrecisionMismatch = errors.New("decimal value and type precisions mismatch")

type typePrecisionWrapper interface {
	TypePrecision() Precision
	GetDecimal() Decimal
	SetDecimal(d Decimal)
}

func marshalText[T typePrecisionWrapper](d T) ([]byte, error) {
	value := d.GetDecimal()
	if value.Precision() != d.TypePrecision() {
		return nil, ErrTypePrecisionMismatch
	}
	return []byte(fmt.Sprintf("%q", value.String())), nil
}

func unmarshalText[T typePrecisionWrapper](d *T, data []byte) error {
	if len(data) < 1 {
		return errors.New("invalid decimal number")
	}

	coins, err := Parse(string(data), (*d).TypePrecision())
	if err != nil {
		return err
	}

	(*d).SetDecimal(coins)

	return nil
}

type DecimalTextTrait interface {
	TextDeci | TextCenti |
		TextMilli | TextMicro |
		TextNano | TextPico |
		TextFemto | TextAtto |
		TextZepto | TextYocto |
		TextRonto | TextQuecto
}

func ref[T DecimalTextTrait](d Decimal) *T {
	v := T(d)
	return &v
}

func TextDeciRef(d Decimal) *TextDeci                     { return ref[TextDeci](d) }
func (_ *TextDeci) TypePrecision() Precision              { return Deci }
func (d *TextDeci) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextDeci) SetDecimal(v Decimal)                  { *d = TextDeci(v) }
func (d *TextDeci) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextDeci) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextCentiRef(d Decimal) *TextCenti                    { return ref[TextCenti](d) }
func (_ *TextCenti) TypePrecision() Precision              { return Centi }
func (d *TextCenti) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextCenti) SetDecimal(v Decimal)                  { *d = TextCenti(v) }
func (d *TextCenti) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextCenti) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextMilliRef(d Decimal) *TextMilli                    { return ref[TextMilli](d) }
func (_ *TextMilli) TypePrecision() Precision              { return Milli }
func (d *TextMilli) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextMilli) SetDecimal(v Decimal)                  { *d = TextMilli(v) }
func (d *TextMilli) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextMilli) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextMicroRef(d Decimal) *TextMicro                    { return ref[TextMicro](d) }
func (_ *TextMicro) TypePrecision() Precision              { return Micro }
func (d *TextMicro) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextMicro) SetDecimal(v Decimal)                  { *d = TextMicro(v) }
func (d *TextMicro) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextMicro) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextNanoRef(d Decimal) *TextNano                     { return ref[TextNano](d) }
func (_ *TextNano) TypePrecision() Precision              { return Nano }
func (d *TextNano) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextNano) SetDecimal(v Decimal)                  { *d = TextNano(v) }
func (d *TextNano) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextNano) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextPicoRef(d Decimal) *TextPico                     { return ref[TextPico](d) }
func (_ *TextPico) TypePrecision() Precision              { return Pico }
func (d *TextPico) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextPico) SetDecimal(v Decimal)                  { *d = TextPico(v) }
func (d *TextPico) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextPico) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextFemtoRef(d Decimal) *TextFemto                    { return ref[TextFemto](d) }
func (_ *TextFemto) TypePrecision() Precision              { return Femto }
func (d *TextFemto) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextFemto) SetDecimal(v Decimal)                  { *d = TextFemto(v) }
func (d *TextFemto) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextFemto) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextAttoRef(d Decimal) *TextAtto                     { return ref[TextAtto](d) }
func (_ *TextAtto) TypePrecision() Precision              { return Atto }
func (d *TextAtto) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextAtto) SetDecimal(v Decimal)                  { *d = TextAtto(v) }
func (d *TextAtto) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextAtto) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextZeptoRef(d Decimal) *TextZepto                    { return ref[TextZepto](d) }
func (_ *TextZepto) TypePrecision() Precision              { return Zepto }
func (d *TextZepto) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextZepto) SetDecimal(v Decimal)                  { *d = TextZepto(v) }
func (d *TextZepto) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextZepto) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextYoctoRef(d Decimal) *TextYocto                    { return ref[TextYocto](d) }
func (_ *TextYocto) TypePrecision() Precision              { return Yocto }
func (d *TextYocto) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextYocto) SetDecimal(v Decimal)                  { *d = TextYocto(v) }
func (d *TextYocto) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextYocto) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextRontoRef(d Decimal) *TextRonto                    { return ref[TextRonto](d) }
func (_ *TextRonto) TypePrecision() Precision              { return Ronto }
func (d *TextRonto) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextRonto) SetDecimal(v Decimal)                  { *d = TextRonto(v) }
func (d *TextRonto) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextRonto) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }

func TextQuectoRef(d Decimal) *TextQuecto                   { return ref[TextQuecto](d) }
func (_ *TextQuecto) TypePrecision() Precision              { return Quecto }
func (d *TextQuecto) GetDecimal() Decimal                   { return Decimal(*d) }
func (d *TextQuecto) SetDecimal(v Decimal)                  { *d = TextQuecto(v) }
func (d *TextQuecto) MarshalText() ([]byte, error)          { return marshalText(d) }
func (d *TextQuecto) UnmarshalText(data []byte) (err error) { return unmarshalText(&d, data) }
