package magick

import (
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

func TestActionInlineImages(t *testing.T) {
	sandbox.Action(t, func() carapace.Action {
		return ActionInlineImages()
	})(func(s *sandbox.Sandbox) {
		s.Run("").Expect(carapace.ActionValuesDescribed(
			"canvas:", "solid color canvas",
			"gradient:", "linear gradient",
			"magick:", "read from stdin",
			"null:", "transparent null image",
			"pattern:", "built-in pattern",
			"plasma:", "plasma fractal image",
			"radial-gradient:", "radial gradient",
			"tile:", "tiled image",
			"xc:", "solid color canvas",
		).Tag("inline images"))
	})
}

func TestActionInlineImagesFiltering(t *testing.T) {
	sandbox.Action(t, func() carapace.Action {
		return ActionInlineImages()
	})(func(s *sandbox.Sandbox) {
		s.Run("gr").Expect(carapace.ActionValuesDescribed(
			"gradient:", "linear gradient",
		).Tag("inline images"))
	})
}

func TestActionInlineImagesPartialXc(t *testing.T) {
	sandbox.Action(t, func() carapace.Action {
		return ActionInlineImages()
	})(func(s *sandbox.Sandbox) {
		s.Run("xc").Expect(carapace.ActionValuesDescribed(
			"xc:", "solid color canvas",
		).Tag("inline images"))
	})
}
