package user

import "golang.org/x/text/language"

// Locale is a user's preferred languages. It is used for localization.
// The format of the string is specified by RFC 2616 but is validated by
// [language.ParseAcceptLanguage], which is more lax.
type Locale string

// ParseLocale parses a locale string into a [Locale] type.
func ParseLocale(locale string) (Locale, error) {
	l := Locale(locale)
	return l, l.Validate()
}

// Tags returns the Locale as a list of language tags.
// If l is empty or invalid, then this function returns one [language.Und]. The
// returned list is never empty.
func (l Locale) Tags() []language.Tag {
	tags, _, _ := language.ParseAcceptLanguage(string(l))
	if len(tags) == 0 {
		return []language.Tag{language.Und}
	}
	return tags
}

// Validate checks if the Locale is valid.
func (l Locale) Validate() error {
	_, _, err := language.ParseAcceptLanguage(string(l))
	return err
}

// String implements the [fmt.Stringer] interface.
func (l Locale) String() string {
	return string(l)
}
