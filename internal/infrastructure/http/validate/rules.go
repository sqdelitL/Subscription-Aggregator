package validate

import (
	"regexp"

	"github.com/invopop/validation"
)

var uuidRegexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

func RuleUUID() validation.Rule {
	return validation.Match(uuidRegexp)
}

// MM-YYYY
var dateRegex = regexp.MustCompile(`^(0[1-9]|1[0-2])-\d{4}$`)

func RuleSubscribeDateFormat() validation.Rule {
	return validation.Match(dateRegex)
}
