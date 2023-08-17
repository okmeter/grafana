package thema

import (
	"bytes"
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/token"
	terrors "github.com/grafana/thema/errors"
)

type onesidederr struct {
	schpos, datapos []token.Pos
	code            terrors.ValidationCode
	coords          coords
	val             string
}

func (e *onesidederr) Error() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s: validation failed, data is not an instance:", e.coords)
	switch e.code {
	case terrors.MissingField:
		fmt.Fprintf(&buf, "\n\tschema specifies that field exists with type %v", e.val)
		for _, pos := range e.schpos {
			fmt.Fprintf(&buf, "\n\t\t%s", pos.String())
		}

		fmt.Fprintf(&buf, "\n\tbut field was absent from data")
		for _, pos := range e.datapos {
			fmt.Fprintf(&buf, "\n\t\t%s", pos.String())
		}
	case terrors.ExcessField:
		fmt.Fprintf(&buf, "\n\tschema is closed and does not specify field")
		for _, pos := range e.schpos {
			fmt.Fprintf(&buf, "\n\t\t%s", pos.String())
		}

		fmt.Fprintf(&buf, "\n\tbut field exists in data with value %v", e.val)
		for _, pos := range e.datapos {
			fmt.Fprintf(&buf, "\n\t\t%s", pos.String())
		}
	}

	return buf.String()
}

func (e *onesidederr) Unwrap() error {
	return terrors.ErrNotAnInstance
}

type twosidederr struct {
	schpos, datapos []token.Pos
	code            terrors.ValidationCode
	coords          coords
	sv, dv          string
}

func (e *twosidederr) Error() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s: validation failed, data is not an instance:\n\tschema expected `%s`", e.coords, e.sv)
	for _, pos := range e.schpos {
		fmt.Fprintf(&buf, "\n\t\t%s", pos.String())
	}

	fmt.Fprintf(&buf, "\n\tbut data contained `%s`", e.dv)
	for _, pos := range e.datapos {
		fmt.Fprintf(&buf, "\n\t\t%s", pos.String())
	}
	return buf.String()
}

func (e *twosidederr) Unwrap() error {
	return terrors.ErrNotAnInstance
}

// TODO differentiate this once we have generic composition to support trimming out irrelevant disj branches
type emptydisjunction struct {
	schpos, datapos []token.Pos
	coords          coords
	brancherrs      []error
}

func (e *emptydisjunction) Unwrap() error {
	return terrors.ErrNotAnInstance
}

type validationFailure []error

func (vf validationFailure) Unwrap() error {
	return terrors.ErrNotAnInstance
}

func (vf validationFailure) Error() string {
	var buf bytes.Buffer
	for _, e := range vf {
		fmt.Fprint(&buf, e.Error())
		fmt.Fprintf(&buf, "\n")
	}

	return buf.String()
}

func mungeValidateErr(err error, sch Schema) error {
	_, is := err.(errors.Error)
	if !is {
		return err
	}

	var errs validationFailure
	for _, ee := range errors.Errors(err) {
		schpos, datapos := splitTokens(ee.InputPositions())
		x := coords{
			sch:       sch,
			fieldpath: trimThemaPath(ee.Path()),
		}

		msg, vals := ee.Msg()
		switch len(vals) {
		case 1:
			val, ok := vals[0].(string)
			if !ok {
				break
			}
			err := &onesidederr{
				schpos:  schpos,
				datapos: datapos,
				coords:  x,
				val:     val,
			}

			if strings.Contains(msg, "incomplete") {
				err.code = terrors.MissingField
			} else if strings.Contains(msg, "not allowed") {
				err.code = terrors.ExcessField
			} else {
				break
			}

			errs = append(errs, err)
			continue
		case 4:
			schval, svok := vals[0].(string)
			dataval, dvok := vals[1].(string)
			schkind, skok := vals[2].(cue.Kind)
			datakind, dkok := vals[3].(cue.Kind)
			if !svok || !dvok || !skok || !dkok {
				break
			}

			err := &twosidederr{
				schpos:  schpos,
				datapos: datapos,
				coords:  x,
				sv:      schval,
				dv:      dataval,
			}
			if datakind.IsAnyOf(schkind) {
				err.code = terrors.OutOfBounds
			} else {
				err.code = terrors.KindConflict
			}

			errs = append(errs, err)
			continue
		}

		// We missed a case, wrap CUE err in a plea for help
		errs = append(errs, fmt.Errorf("no Thema handler for CUE error, please file an issue against github.com/grafana/thema\nto improve this error output!\n\n%w", ee))
	}
	return errs
}

func splitTokens(poslist []token.Pos) (schpos, datapos []token.Pos) {
	if len(poslist) == 0 {
		return
	}

	// We're assuming data is always last. ...Probably safe? Given that we
	// control the order of operands in the Schema.Validate() calls...
	dataname := poslist[len(poslist)-1].Filename()
	var split int
	for i, pos := range poslist {
		if pos.Filename() == dataname {
			split = i
			break
		}
	}

	return poslist[:split], poslist[split:]
}

func trimThemaPath(parts []string) []string {
	for i, s := range parts {
		if s == "seqs" {
			return parts[i+4:]
		}
	}

	// Otherwise, it's one of the defpath patterns - eliminate first element
	return parts[1:]
}

type coords struct {
	sch       Schema
	fieldpath []string
}

func (c coords) String() string {
	return fmt.Sprintf("<%s@v%s>.%s", c.sch.Lineage().Name(), c.sch.Version(), strings.Join(c.fieldpath, "."))
}
