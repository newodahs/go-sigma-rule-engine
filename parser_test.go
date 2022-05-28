package sigma

import (
	"encoding/json"
	"testing"

	"github.com/gobwas/glob"
	"github.com/markuskont/datamodels"
	"gopkg.in/yaml.v2"
)

var detection1 = `
detection:
  condition: "selection1 and not selection3"
  selection1:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
  selection3:
    CommandLine: "+R +H +S +A *.cui"
`

var detection1_positive = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +A asd.cui",
	"ParentImage": "C:\\test\\wmiprvse.exe",
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "aaa",
	"ParentImage": "C:\\test\\wmiprvse.exe"
}
`

var detection1_negative1 = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +S +A lll.cui",
	"ParentImage": "C:\\test\\mshta.exe"
}
`
var detection1_negative2 = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +S +A lll.cui"
}
`
var detection2 = `
detection:
  condition: "(selection1 and selection2) and not selection3"
  selection1:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
  selection2:
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
  selection3:
    CommandLine: "+R +H +S +A *.cui"
`

var detection3 = `
detection:
  condition: "(selection1 or selection2) and not selection3"
  selection1:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
  selection2:
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
  selection3:
    CommandLine: "+R +H +S +A *.cui"
`

var detection3_positive1 = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +A asd.cui",
	"ParentImage": "C:\\test\\custom.exe",
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "aaa",
	"ParentImage": "C:\\test\\wmiprvse.exe"
}
`
var detection3_positive2 = `
{
	"Image":       "C:\\test\\custom.exe",
	"CommandLine": "+R +H +A asd.cui",
	"ParentImage": "C:\\test\\wmiprvse.exe",
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "aaa",
	"ParentImage": "C:\\test\\wmiprvse.exe"
}
`

var detection3_negative = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +S +A lll.cui",
	"ParentImage": "C:\\test\\mshta.exe"
}
`

var detection4 = `
detection:
  condition: "all of selection* and not filter"
  selection1:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
  selection2:
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
  filter:
    CommandLine: "+R +H +S +A *.cui"
`

var detection5 = `
detection:
  condition: "1 of selection* and not filter"
  selection1:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
  selection2:
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
  filter:
    CommandLine: "+R +H +S +A *.cui"
`

var detection6 = `
detection:
  condition: "all of them"
  selection_images:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
  selection_parent_images:
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
`

var detection6_positive = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +A asd.cui",
	"ParentImage": "C:\\test\\wmiprvse.exe",
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "aaa",
	"ParentImage": "C:\\test\\wmiprvse.exe"
}
`

var detection6_negative = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +S +A lll.cui",
	"ParentImage": "C:\\test\\mshta\\lll.exe"
}
`

var detection7 = `
detection:
  condition: "1 of them"
  selection_images:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
  selection_parent_images:
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
`

var detection7_negative1 = `
{
	"Image":       "C:\\test\\bytesadmin.exe",
	"CommandLine": "+R +H +S +A lll.cui",
	"ParentImage": "E:\\go\\bin\\gofmt"
}
`
var detection7_negative2 = `
{
	"Image":       "C:\\test\\bytesadmin.exe",
	"CommandLine": "+R +H +S +A lll.cui"
}
`

var detection8 = `
detection:
  condition: "selection1 and not selection3"
  selection1:
    Image:
    - '*\schtasks.exe'
    - '*\nslookup.exe'
    - '*\certutil.exe'
    - '*\bitsadmin.exe'
    - '*\mshta.exe'
    ParentImage:
    - '*\mshta.exe'
    - '*\powershell.exe'
    - '*\cmd.exe'
    - '*\rundll32.exe'
    - '*\cscript.exe'
    - '*\wscript.exe'
    - '*\wmiprvse.exe'
  selection3:
    CommandLine: "+R +H +S +A *.cui"
`

var detection8_positive = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +A asd.cui",
	"ParentImage": "C:\\test\\wmiprvse.exe",
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "aaa",
	"ParentImage": "C:\\test\\wmiprvse.exe"
}
`

var detection8_negative1 = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"CommandLine": "+R +H +S +A lll.cui",
	"ParentImage": "C:\\test\\mshta.exe"
}
`
var detection8_negative2 = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"ParentImage": "C:\\test\\mshta.exe"
}
`

var detection9 = `
detection:
  condition: "selection"
  selection:
    - PipeName|re: '\\\\SomePipeName[0-9a-f]{2}'
    - PipeName2|re: '\\\\AnotherPipe[0-9a-f]*Name'
`

var detection9_positive = `
{
	"PipeName":       "\\\\SomePipeNamea4",
	"PipeName2":       "\\\\AnotherPipe0af3Name"
}
`

var detection9_negative = `
{
	"PipeName":       "\\\\SomePipeNameZZ",
	"PipeName2":       "\\\\AnotherPipe01zzName"
}
`

var detection10 = `
detection:
  condition: "selection1 and selection2"
  selection1:
    - SomeName|startswith: 'TestStart'
  selection2:
    - SomeName|endswith: 'TestEnd'
`

var detection10_positive = `
{
	"SomeName":       "TestStart-Value-TestEnd"
}
`

var detection10_negative = `
{
	"SomeName":       "TestStart-Value"
}
`

var detection11 = `
detection:
  condition: "selection1 and selection2"
  selection1:
    SomeName|contains|all: 
      - 'mark1'
      - 'mark2'
  selection2:
    SomeName|contains:
      - 'version1'
      - 'version2'
`

var detection11_positive = `
{
	"SomeName":       "Some mark1 mark2 String version2"
}
`

var detection11_negative = `
{
	"SomeName":       "mark1 mark2 mark3 non-matching string"
}
`

var detection12 = `
detection:
  condition: "selection1 and selection2"
  selection1:
    SomeKey|contains|all:
      - 'val1'
      - 'val2'
  selection2:
    SomeKey2:
      - 'mustMatch1'
      - 'mustMatch2'
`

var detection12_positive = `
{
	"SomeKey":       "val1 val2",
	"SomeKey2":      "mustMatch1"
}
`

var detection12_negative = `
{
	"SomeKey":       "val1 val2",
	"SomeKey2":      "mustMatch3"
}
`

//this test is a bit tricky:
//the '*\bits\*admin.exe' is looking to match '[wildCard]\bits*admin.exe' (one wildcard at head, one escaped wildcard)
//the '\\\\DoubleBackslash\\some*.exe' is looking to match '\\DoubleBackslash\some[wildCard].exe' (multiple backslashes, one wildcard)
//the '\leadingBackslash\\*.exe' is looking to match '\leadingBackslash\[wildCard].exe' (one wildcard and leading backslash)
//the 'full\\\*plaintext.exe' is looking to match 'full\*plaintext.exe' (no wildcards exact match)
var detection13 = `
detection:
  condition: "all of them"
  selection_images:
    Image:
    - '*\bits\*admin.exe'
    - '\\\\DoubleBackslash\\some*.exe'
    - '[Windows-*]\image.???'
  selection_parent_images:
    ParentImage:
    - '\leadingBackslash\\*.exe'
    - 'full\\\*plaintext.exe'
    - '{000-aaa-*}\\*.exe'
`

var detection13_positive = `
{
	"Image":       "C:\\test\\bits*admin.exe",
	"ParentImage": "\\leadingBackslash\\something.exe"
}
`
var detection13_positive2 = `
{
	"Image":       "\\\\DoubleBackslash\\someOther.exe",
	"ParentImage": "full\\*plaintext.exe"
}
`

var detection13_positive3 = `
{
	"Image":       "C:\\test\\bits*admin.exe",
	"ParentImage": "full\\*plaintext.exe"
}
`

var detection13_positive4 = `
{
	"Image":       "[Windows-Security]\\image.cmd",
	"ParentImage": "{000-aaa-123}\\evil.exe"
}
`

//won't match as Image is looking for '*\bits*admin.exe' witha leading wildcard and an escaped '*' between bits and admin
//this provides 'C:\test\bitsadmin.exe', which matches the leading wildcard but fails to present the escaped '*'
var detection13_negative = `
{
	"Image":       "C:\\test\\bitsadmin.exe",
	"ParentImage": "\\leadingBackslash\\something.exe"
}
`

//won't match as the ParentImage is looking for '\leadingBackslash\*.exe' with a wildcard
//this provides 'leadingBackslash\something.exe', missing the leading backslash
var detection13_negative2 = `
{
	"Image":       "C:\\test\\bits*admin.exe",
	"ParentImage": "leadingBackslash\\something.exe"
}
`

//won't match as the ParentImage is looking for an exact match (no wildcards) to 'full\*plaintext.exe'
//this provides 'full\\*plaintext', the extra backslash kills it
var detection13_negative3 = `
{
	"Image":       "C:\\test\\bits*admin.exe",
	"ParentImage": "full\\\\*plaintext"
}
`

//shouldn't match on either of these (Image is missing 'Windows' in the bracket, ParentImage is missing the
//a vaule of 000-aaa in the brackets)
var detection13_negative4 = `
{
	"Image":       "[-Security]\\image.cmd",
	"ParentImage": "{000-aaa}\\evil.exe"
}
`

//this has a hacky test; we set 'noCollapseWSNeg' in the parseTestCast struct for this case specifically
//doing so will turn off collapsing the whitespace for the negative test and cause this to fail detection
var detection14 = `
detection:
  condition: "selection"
  selection:
    SomeName|contains:
      - 'whitespace   collapse	testing'
`

var detection14_case = `
{
	"SomeName":       "whitespace\t\tcollapse         testing"
}
`

var detection15 = `
detection:
  condition: "selection1 and selection2 and selection3"
  selection1:
    SomeKey: '%ValuePlaceholder%'
  selection2:
    SomeKey2:
    - '%ValuePlaceholder2%'
    - '%%NotAPlaceholder'
    - '%NotAPlaceholder%/APath/In/Windows'
  selection3:
    SomeKey3|contains:
    - '%PlaceholderCannotBeInContains%'
    - 'StillnotA%Placeholder%'
    - 'Wild%Placeholder'
`

//case 1: first two keys should have values in their respective maps; third key can never be a placeholder so must match as a contains
var detection15_positive = `
{
	"SomeKey":		"ThisValueShouldBeInMap",
	"SomeKey2":		"ThisValueShouldBeInMap2",
	"SomeKey3":		"%PlaceholderCannotBeInContains%"
}
`

//case 2: Mix of placeholder (SomeKey) and non-placeholder matching where there is a placeholder mixed in (SomeKey2)
var detection15_positive2 = `
{
	"SomeKey":		"ThisIsAlsoInTheMap",
	"SomeKey2":		"%NotAPlaceholder%/APath/In/Windows",
	"SomeKey3":		"In Contains StillnotA%Placeholder%"
}
`

//case 3: Mix of placeholder (SomeKey) and non-placeholder matching where there is a placeholder mixed in (SomeKey2); different matching of non-placeholder
var detection15_positive3 = `
{
	"SomeKey":		"ThisValueShouldBeInMap",
	"SomeKey2":		"%%NotAPlaceholder",
	"SomeKey3":		"Wild%Placeholder In Contains"
}
`

//case 1: SomeKey not found (placeholder)
var detection15_negative = `
{
	"SomeKey":		"ThisValueIsNOTInMap",
	"SomeKey2":		"ThisValueShouldBeInMap2",
	"SomeKey3":		"%PlaceholderCannotBeInContains%"
}
`

//case 2: SomeKey2 not found
var detection15_negative2 = `
{
	"SomeKey":		"ThisValueShouldBeInMap",
	"SomeKey2":		"ThisValueIsNOTInMap",
	"SomeKey3":		"%PlaceholderCannotBeInContains%"
}
`

//case 3: Placeholders match, non-placeholder does not
var detection15_negative3 = `
{
	"SomeKey":		"ThisValueShouldBeInMap",
	"SomeKey2":		"ThisValueShouldBeInMap2",
	"SomeKey3":		"NotFoundAnywhere"
}
`

type parseTestCase struct {
	Rule            string
	Pos, Neg        []string
	noCollapseWSNeg bool
	lookupMap       PlaceholderLookup
}

//a simplistic lookup function for testing placeholders
func PlaceholderLookupT1(key string) ([]string, bool) {
	var testMap = map[string][]string{
		"ValuePlaceholder":  {"ThisValueShouldBeInMap", "ThisIsAlsoInTheMap"},
		"ValuePlaceholder2": {"ThisValueShouldBeInMap2", "AnotherValue"},
	}

	set, found := testMap[key]
	if !found {
		return nil, false
	}
	return set, found
}

var parseTestCases = []parseTestCase{
	{
		Rule: detection1,
		Pos:  []string{detection1_positive},
		Neg:  []string{detection1_negative1, detection1_negative2},
	},
	{
		Rule: detection2,
		Pos:  []string{detection1_positive},
		Neg:  []string{detection1_negative1, detection1_negative2},
	},
	{
		Rule: detection3,
		Pos:  []string{detection3_positive1, detection3_positive2},
		Neg:  []string{detection3_negative},
	},
	{
		Rule: detection4,
		Pos:  []string{detection1_positive},
		Neg:  []string{detection1_negative1, detection1_negative2},
	},
	{
		Rule: detection5,
		Pos:  []string{detection3_positive1, detection3_positive2},
		Neg:  []string{detection3_negative},
	},
	{
		Rule: detection6,
		Pos:  []string{detection6_positive},
		Neg:  []string{detection6_negative},
	},
	{
		Rule: detection7,
		Pos:  []string{detection3_positive1, detection3_positive2},
		Neg:  []string{detection7_negative1, detection7_negative2},
	},
	{
		Rule: detection8,
		Pos:  []string{detection8_positive},
		Neg:  []string{detection8_negative1, detection8_negative2},
	},
	{
		Rule: detection9,
		Pos:  []string{detection9_positive},
		Neg:  []string{detection9_negative},
	},
	{
		Rule: detection10,
		Pos:  []string{detection10_positive},
		Neg:  []string{detection10_negative},
	},
	{
		Rule: detection11,
		Pos:  []string{detection11_positive},
		Neg:  []string{detection11_negative},
	},
	{
		Rule: detection12,
		Pos:  []string{detection12_positive},
		Neg:  []string{detection12_negative},
	},
	{
		Rule: detection13,
		Pos:  []string{detection13_positive, detection13_positive2, detection13_positive3, detection13_positive4},
		Neg:  []string{detection13_negative, detection13_negative2, detection13_negative3, detection13_negative4},
	},
	{
		Rule:            detection14,
		Pos:             []string{detection14_case},
		noCollapseWSNeg: false, //ensures whitespace is collapsed and everything matches
	},
	{
		Rule:            detection14,
		Neg:             []string{detection14_case},
		noCollapseWSNeg: true, //turns off whitespace collapsing and causing a non-match
	},
	{
		Rule:      detection15, //placeholder testing
		Pos:       []string{detection15_positive, detection15_positive2, detection15_positive3},
		Neg:       []string{detection15_negative, detection15_negative2, detection15_negative3},
		lookupMap: PlaceholderLookupT1,
	},
}

func TestTokenCollect(t *testing.T) {
	for _, c := range LexPosCases {
		p := &parser{
			lex: lex(c.Expr),
		}
		if err := p.collect(); err != nil {
			switch err.(type) {
			case ErrUnsupportedToken:
			default:
				t.Fatal(err)
			}
		}
	}
}

func TestParse(t *testing.T) {
	for i, c := range parseTestCases {
		var rule Rule
		if err := yaml.Unmarshal([]byte(c.Rule), &rule); err != nil {
			t.Fatalf("rule parse case %d failed to unmarshal yaml, %s", i+1, err)
		}
		expr := rule.Detection["condition"].(string)
		p := &parser{
			lex:          lex(expr),
			sigma:        rule.Detection,
			noCollapseWS: c.noCollapseWSNeg,
		}
		if err := p.collect(); err != nil {
			t.Fatalf("rule parser case %d failed to collect lexical tokens, %s", i+1, err)
		}
		if err := p.parse(); err != nil {
			switch err.(type) {
			case ErrWip:
				t.Fatalf("WIP")
			default:
				t.Fatalf("rule parser case %d failed to parse lexical tokens, %s", i+1, err)
			}
		}
		var obj datamodels.Map
		// Positive cases
		for _, x := range c.Pos {
			if err := json.Unmarshal([]byte(x), &obj); err != nil {
				t.Fatalf("rule parser case %d positive case json unmarshal error %s", i+1, err)
			}
			m, _ := p.result.MatchEx(obj, c.lookupMap)
			if !m {
				t.Fatalf("rule parser case %d positive case did not match: %+v", i+1, obj)
			}
		}
		// Negative cases
		for _, x := range c.Neg {
			if err := json.Unmarshal([]byte(x), &obj); err != nil {
				t.Fatalf("rule parser case %d positive case json unmarshal error %s", i+1, err)
			}
			m, _ := p.result.MatchEx(obj, c.lookupMap)
			if m {
				t.Fatalf("rule parser case %d negative case matched", i+1)
			}
		}
	}
}

func TestSigmaEscape(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expected   string
		validMatch string
		skip       bool
	}{
		{
			name:       "No_Change",
			input:      `\\leadingBackslash\\*.exe`,
			expected:   `\\leadingBackslash\\*.exe`,
			validMatch: `\leadingBackslash\testing.exe`,
		},
		{
			name:       "Leading_Single_Backslash_Wildcard_After_Slash",
			input:      `\leadingBackslash\\*.exe`,
			expected:   `\\leadingBackslash\\*.exe`,
			validMatch: `\leadingBackslash\testing.exe`,
		},
		{
			name:       "Leading_Wildcard_Single_Backslash_Esc_Wildcard",
			input:      `*\bits\*admin.exe`,
			expected:   `*\\bits\*admin.exe`,
			validMatch: `leading\bits*admin.exe`,
		},
		{
			name:       "Double_Leading_Backslash_Single_Backslash_Wildcard",
			input:      `\\\\DoubleBackslash\some*.exe`,
			expected:   `\\\\DoubleBackslash\\some*.exe`,
			validMatch: `\\DoubleBackslash\sometMatch.exe`,
		},
		{
			name:       "Plaintext_Only_Esc_Wildcard",
			input:      `some\full\\\*plaintext.exe`,
			expected:   `some\\full\\\*plaintext.exe`,
			validMatch: `some\full\*plaintext.exe`,
		},
		{
			name:       "Double_Leading_Backslash_Complex_Mix_Esc",
			input:      `\\\\DoubleBackslash\?\some*Other\\*test.\\???`,
			expected:   `\\\\DoubleBackslash\?\\some*Other\\*test.\\???`,
			validMatch: `\\DoubleBackslash?\someMixOther\wildcardtest.\cmd`,
		},
		{
			name:       "Mixed_Wildcards_Single_Backslash_Brackets",
			input:      `[*]\*\aSetof\\\sigma{rule?}here*`,
			expected:   `\[*\]\*\\aSetof\\\\sigma\{rule?\}here*`,
			validMatch: `[testing]*\aSetof\\sigma{rules}hereWeGo`,
		},
	}
	for _, curTest := range tests {
		t.Run(curTest.name, func(t *testing.T) {
			if curTest.skip {
				t.Skip("test marked as skip")
			}

			escStr := escapeSigmaForGlob(curTest.input)
			if escStr != curTest.expected {
				t.Errorf("failed to validate escaped input; got: %s - expected: %s", escStr, curTest.expected)
			}

			//test as a glob to be sure
			globT, err := glob.Compile(escStr)
			if err != nil {
				t.Fatalf("failed to compile glob: %+v", err)
			}
			if !globT.Match(curTest.validMatch) {
				t.Errorf("compiled glob did NOT match valid input; glob: %s -- data: %s", escStr, curTest.validMatch)
			}
		})
	}
}
