package validators

type FormValidation struct {
	FieldErros map[string]string
}

func (fv *FormValidation) Valid() bool {
	return len(fv.FieldErros) == 0
}

func (fv *FormValidation) AddFieldErrors(field, message string) {
	if fv.FieldErros == nil {
		fv.FieldErros = make(map[string]string)
	}
	fv.FieldErros[field] = message

}
