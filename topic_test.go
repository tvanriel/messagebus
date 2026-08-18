package messagebus_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tvanriel/messagebus"
)

//nolint:funlen // mostly tables.
func TestTopic_Matches(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		topic   string
		matches []string
		nomatch []string
	}{
		{
			name:  "no wildcard should only match exactly",
			topic: "this.is.a.topic",
			matches: []string{
				"this.is.a.topic",
			},
			nomatch: []string{
				"his.is.a.topic",
				"",
				"this.is.a.",
				"this.is.a",
				"this.is.a.topic.test",
			},
		},

		{
			name:  "one-level-wildcard at the end matches exactly one word",
			topic: "test.*",
			matches: []string{
				"test.test",
				"test.something",
				"test.aaaaaaaaaaaaaaaaa",
			},
			nomatch: []string{
				"test",
				"test.",
				"s.test",
				"s.",
				"test.test.test",
				"test.test.",
			},
		},

		{
			name:  "one-level-wildcard in the middle matches exactly one word",
			topic: "test.*.else",
			matches: []string{
				"test.something.else",
				"test.test.else",
			},
			nomatch: []string{
				"test..else",
				"test.happy.els",
				"test.happy",
				"is.something.else",
				"test.something.",
			},
		},

		{
			name:  "one-level-wildcard at the start matches exactly one word",
			topic: "*.something.else",
			matches: []string{
				"test.something.else",
				"is.something.else",
			},
			nomatch: []string{
				"test.something.else.",
				"test.something.else.with",
				"is.something.else.",
				".something.else",
				".something.else.test",
				"is.something",
				"is.something.",
			},
		},
		{
			name:  "multiple one-level-wildcard in the middle matches exactly one word",
			topic: "test.*.*.else",
			matches: []string{
				"test.test.something.else",
				"test.xx.s.else",
			},
			nomatch: []string{
				"test.test.something.else.",
				"test...else",
				"test.x..else",
				"test..x.else",
				"test.something.else",
				"test..else",
				"test.happy.els",
				"test.happy",
				"is.something.else",
				"test.something.",
			},
		},
		{
			name:  "multiple one-level-wildcard at the end matches exactly one word",
			topic: "test.*.*",
			matches: []string{
				"test.something.else",
				"test.xx.s",
			},
			nomatch: []string{
				"test..",
				"test.glamour",
				"test.glamour.too.many.words",
				"test.glamour.too.many.words.",
			},
		},
		{
			name:  "multiple one-level-wildcard in the middle matches exactly one word",
			topic: "test.*.*.else",
			matches: []string{
				"test.test.something.else",
				"test.xx.s.else",
			},
			nomatch: []string{
				"test.test.something.else.",
				"test...else",
				"test.x..else",
				"test..x.else",
				"test.something.else",
				"test..else",
				"",
				"test.happy.els",
				"test.happy",
				"is.something.else",
				"test.something.",
			},
		},
		{
			name:  "multiple one-level-wildcard at the end matches exactly one word",
			topic: "test.*.*",
			matches: []string{
				"test.something.else",
				"test.xx.s",
			},
			nomatch: []string{
				"test..",
				"test.glamour",
				"test.glamour.too.many.words",
				"",
				"test.glamour.too.many.words.",
			},
		},
		{
			name:  "multiple one-level-wildcard in the middle matches exactly one word",
			topic: "*.*.else",
			matches: []string{
				"test.something.else",
				"never.someone.else",
			},
			nomatch: []string{
				"contains.too.many.words.for.else",
				"words.for.else.",
				"for.else",
				"for.else.",
				"..else",
				"..else.",
				"",
			},
		},

		{
			name:  "multiple-level-card at the start matches zero-or-more words",
			topic: "#.test",
			matches: []string{
				"everything.that.ends.in.test",
				".everything.that.ends.in.test",
				"test",
				".test",
				"short.test",
				".short.test",
			},
			nomatch: []string{
				"doesnot.end.test.",
				"period",
				"",
				"typo.tst",
				".typo.tst",
			},
		},

		{
			name:  "multiple-level-card at the end matches zero-or-more words",
			topic: "test.#",
			matches: []string{
				"test.is.concise",
				"test.the.world.",
				"test",
				"test.",
				"test.me",
				"test.me.",
			},
			nomatch: []string{
				".test.it",
				"period",
				"",
				"tst.typoed",
				"tst.typoed.",
			},
		},

		{
			name:  "multiple-level-card in the middle matches zero-or-more words",
			topic: "test.#.it",
			matches: []string{
				"test.all.of.it",
				"test.it",
				"test..it",
				"test.some-of.it",
			},
			nomatch: []string{
				"test.some-of.it.",
				".test.some-of.it",
				".test.that",
				"be.it",
				"be.that",
				"",
			},
		},

		{
			name:  "only wildcard matches everything",
			topic: "#",
			matches: []string{
				"test.all.of.it",
				"test.it",
				"test..it",
				"test.some-of.it",
				"test.some-of.it.",
				".test.some-of.it",
				".test.that",
				"be.it",
				"be.that",
				"",
			},
			nomatch: []string{},
		},

		{
			name:  "only single-level-wildcard matches exactly one word",
			topic: "*",
			matches: []string{
				"it",
				"test",
			},
			nomatch: []string{
				"test..it",
				"test.some-of.it",
				"test.some-of.it.",
				".test.some-of.it",
				".test.that",
				"be.it",
				"be.that",
				"",
			},
		},

		{
			name:  "empty topic matches only emptystring",
			topic: "",
			matches: []string{
				"",
			},
			nomatch: []string{
				"it",
				"test",
				"test..it",
				"test.some-of.it",
				"test.some-of.it.",
				".test.some-of.it",
				".test.that",
				"be.it",
				"be.that",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			topic := messagebus.NewTopicFromString(tt.topic)

			for _, m := range tt.matches {
				r := messagebus.NewRouteKeyFromString(m)

				assert.True(t, topic.Matches(r))
			}

			for _, m := range tt.nomatch {
				r := messagebus.NewRouteKeyFromString(m)

				assert.False(t, topic.Matches(r))
			}
		})
	}
}
