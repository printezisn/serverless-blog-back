package model

import (
	"testing"
)

// TestValidateWithRequiredErrors tests that Validate returns errors for required fields.
func TestValidateWithRequiredErrors(t *testing.T) {
	testCases := []struct {
		post      BlogPost
		hasErrors bool
	}{
		{BlogPost{}, true},
		{BlogPost{ID: "test_id"}, true},
		{BlogPost{ID: "test_id", Title: "test_title"}, true},
		{BlogPost{ID: "test_id", Title: "test_title", Description: "test_description"}, true},
		{BlogPost{ID: "test_id", Title: "test_title", Description: "test_descr", Revision: 1}, true},
		{BlogPost{ID: "test_id", Title: "test_title", Description: "test_descr", Revision: 1, Body: "test_body"}, false},
	}

	for _, testCase := range testCases {
		errs := testCase.post.Validate()
		if testCase.hasErrors && len(errs) == 0 {
			t.Error("The following test case was supposed to have errors, but it didn't: ", testCase)
		} else if !testCase.hasErrors && len(errs) > 0 {
			t.Error("The following test case wasn't supposed to have errors, but it did: ", testCase)
		}
	}
}

// TestValidateWithLengthErrors tests that Validate returns errors for fields that exceed a certain length.
func TestValidateWithLengthErrors(t *testing.T) {
	longStr := "a"
	for i := 0; i < 300; i++ {
		longStr += "a"
	}

	testCases := []struct {
		post      BlogPost
		hasErrors bool
	}{

		{BlogPost{ID: "test_id", Title: "test_title", Description: "test_descr", Body: "test_body", Revision: 1}, false},
		{BlogPost{ID: longStr, Title: "test_title", Description: "test_descr", Body: "test_body", Revision: 1}, true},
		{BlogPost{ID: "test_id", Title: longStr, Description: "test_descr", Body: "test_body", Revision: 1}, true},
		{BlogPost{ID: "test_id", Title: "test_title", Description: longStr, Body: "test_body", Revision: 1}, true},
		{BlogPost{ID: "test_id", Title: "test_title", Description: "test_descr", Body: longStr, Revision: 1}, false},
	}

	for _, testCase := range testCases {
		errs := testCase.post.Validate()
		if testCase.hasErrors && len(errs) == 0 {
			t.Error("The following test case was supposed to have errors, but it didn't: ", testCase)
		} else if !testCase.hasErrors && len(errs) > 0 {
			t.Error("The following test case wasn't supposed to have errors, but it did: ", testCase)
		}
	}
}
