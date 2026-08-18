package messagebus

import "strings"

type RouteKey string

func NewRouteKeyFromParts(parts []string) RouteKey {
	return RouteKey(strings.Join(parts, TopicSeparator))
}

func NewRouteKeyFromString(rk string) RouteKey {
	return RouteKey(rk)
}

func (t RouteKey) Parts() []string {
	return strings.Split(string(t), TopicSeparator)
}

type Topic string

const (
	TopicSeparator      = "."
	MultiLevelWildcard  = "#"
	SingleLevelWildcard = "*"
)

func NewTopicFromParts(parts []string) Topic {
	return Topic(strings.Join(parts, TopicSeparator))
}

func NewTopicFromString(topic string) Topic {
	return Topic(topic)
}

func (t Topic) Parts() []string {
	return strings.Split(string(t), TopicSeparator)
}

//nolint:cyclop // todo
func (t Topic) Matches(rk RouteKey) bool {
	topicParts := t.Parts()
	routeParts := rk.Parts()

	var match func(int, int) bool

	match = func(topicIndex, routeKeyIndex int) bool {
		if topicIndex == len(topicParts) && routeKeyIndex == len(routeParts) {
			return true
		}

		if topicIndex == len(topicParts) {
			return false
		}

		switch topicParts[topicIndex] {
		case MultiLevelWildcard:
			// '#' matches zero or more parts.
			if match(topicIndex+1, routeKeyIndex) {
				return true
			}

			if routeKeyIndex < len(routeParts) {
				return match(topicIndex, routeKeyIndex+1)
			}

			return false

		case SingleLevelWildcard:
			// '*' matches exactly one non-empty part.
			if routeKeyIndex >= len(routeParts) || routeParts[routeKeyIndex] == "" {
				return false
			}

			return match(topicIndex+1, routeKeyIndex+1)

		default:
			if routeKeyIndex >= len(routeParts) || topicParts[topicIndex] != routeParts[routeKeyIndex] {
				return false
			}

			return match(topicIndex+1, routeKeyIndex+1)
		}
	}

	return match(0, 0)
}
