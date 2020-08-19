package plaintext

import (
	"fmt"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
	"github.com/SENERGY-Platform/marshaller/lib/marshaller/serialization/base"
)

type Marshaller struct {
}

const Format = "plain-text"

func init() {
	base.Register(Format, Marshaller{})
}

func (Marshaller) Marshal(in interface{}, variable model.ContentVariable) (out string, err error) {
	return fmt.Sprint(in), nil
}

func (Marshaller) Unmarshal(in string, variable model.ContentVariable) (out interface{}, err error) {
	return in, nil
}
