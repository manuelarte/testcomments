package analyzer

import (
	"errors"

	"github.com/manuelarte/testcomments/analyzer/checks"
)

const (
	EqualityComparisonCheckNamePrefix  = "equality-comparison"
	EqualityComparisonReflectCheckName = EqualityComparisonCheckNamePrefix + ".reflect"
	EqualityComparisonEqualCheckName   = EqualityComparisonCheckNamePrefix + ".equal"
	GotBeforeWantCheck                 = "got-before-want"
	IdentifyTheFunctionCHeck           = "identify-function"
	TableDrivenFormatCheckNamePrefix   = "table-driven-format"
	TableDrivenFormatCheckTypeName     = TableDrivenFormatCheckNamePrefix + ".type"
	TableDrivenFormatCheckInlinedName  = TableDrivenFormatCheckNamePrefix + ".inlined"
)

type (
	Settings struct {
		EqualityComparison EqualityComparison
		GotBeforeWant      bool
		IdentifyFunction   bool
		TableDrivenFormat  TableDrivenFormat
	}

	EqualityComparison struct {
		Reflect bool
		Equal   bool
	}

	TableDrivenFormat struct {
		FormatType string
		Inlined    bool
	}
)

func DefaultSettings() Settings {
	return Settings{
		EqualityComparison: EqualityComparison{
			Reflect: true,
			Equal:   true,
		},
		GotBeforeWant:     true,
		IdentifyFunction:  true,
		TableDrivenFormat: TableDrivenFormat{},
	}
}

func SettingsFrom(settings any) (Settings, error) {
	casted, ok := settings.(map[string]any)
	if !ok {
		return Settings{}, errors.New("invalid settings type, expected map[string]any")
	}

	s := DefaultSettings()

	if tableDrivenFormat, okTableDrivenFormat := casted[TableDrivenFormatCheckNamePrefix]; okTableDrivenFormat {
		mapTableDrivenFormat, okMapTableDrivenFormat := tableDrivenFormat.(map[string]any)
		if !okMapTableDrivenFormat {
			return Settings{}, errors.New("invalid TableDrivenFormat settings, expected map[string]any")
		}

		if formatType, okFormatType := mapTableDrivenFormat["type"]; okFormatType {
			stringFormatType, okStringFormatType := formatType.(string)
			if !okStringFormatType {
				return Settings{}, errors.New("invalid TableDrivenFormatType setting, expected string")
			}

			if stringFormatType != "map" && stringFormatType != "slice" && stringFormatType != "" {
				return Settings{}, errors.New("invalid TableDrivenFormatType setting, expected empty, map or slice")
			}

			s.TableDrivenFormat.FormatType = stringFormatType
		}

		if inlined, okInlined := mapTableDrivenFormat["inlined"]; okInlined {
			booleanInlined, okBooleanInlined := inlined.(bool)
			if !okBooleanInlined {
				return Settings{}, errors.New("invalid TableDrivenFormatInlined setting, expected bool")
			}

			s.TableDrivenFormat.Inlined = booleanInlined
		}
	}

	return s, nil
}

func (t TableDrivenFormat) getTableDrivenFormatPredicate() checks.TableDrivenFormatPredicate {
	f := checks.TableDrivenFormatType(t.FormatType)
	if f != checks.Map && f != checks.Slice {
		return checks.AlwaysValid()
	}

	pred, _ := checks.OfTypeAndInline(f, t.Inlined)

	return pred
}
