package flag

import (
	"fmt"
	"github.com/ebauman/golicense/pkg/types"
	"regexp"
	"strconv"
)

var (
	GrantRegex *regexp.Regexp
	MetaRegex  *regexp.Regexp
)

func init() {
	GrantRegex = regexp.MustCompile(types.GrantRegex)
	MetaRegex = regexp.MustCompile(types.MetadataRegex)
}

func ParseGrantFlags(flags []string) (map[string]int, error) {
	var out = map[string]int{}
	for _, flag := range flags {
		parseResult, err := matchGrant(flag)
		if err != nil {
			return nil, err
		}

		if len(parseResult) < 3 {
			return nil, types.NewInvalidGrantError(flag) // regex didn't match
		}

		// second portion is an integer, so parse it
		i, err := strconv.ParseInt(parseResult["amount"], 10, 32)
		if err != nil {
			return nil, types.NewInvalidGrantError(flag)
		}

		out[fmt.Sprintf("%s/%s", parseResult["domain"], parseResult["unit"])] = int(i)
	}

	return out, nil
}

func ParseMetadataFlags(flags []string) (map[string]string, error) {
	var out = map[string]string{}

	for _, flag := range flags {
		parseResult, err := matchGrant(flag)
		if err != nil {
			return nil, err
		}

		if len(parseResult) < 3 {
			return nil, types.NewInvalidMetadataError(flag)
		}

		out[fmt.Sprintf("%s/%s", parseResult["domain"], parseResult["label"])] = parseResult["value"]
	}

	return out, nil
}

func matchGrant(grantString string) (map[string]string, error) {
	match := GrantRegex.FindStringSubmatch(grantString)
	if match == nil {
		return nil, types.NewInvalidGrantError(grantString)
	}

	return regexNames(GrantRegex, match), nil
}

func matchMeta(metaString string) (map[string]string, error) {
	match := MetaRegex.FindStringSubmatch(metaString)
	if match == nil {
		return nil, types.NewInvalidMetadataError(metaString)
	}

	return regexNames(MetaRegex, match), nil
}

func regexNames(r *regexp.Regexp, match []string) map[string]string {
	result := make(map[string]string)

	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result
}
