package cmd

import (
	"sort"
	"testing"
	"time"
)

func TestLen(t *testing.T) {
	expected := []struct {
		images ByLastPushedDesc
		length int
	}{
		{
			images: ByLastPushedDesc([]Image{
				{
					Digest:     "three",
					LastPushed: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					Digest:     "one",
					LastPushed: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Digest:     "two",
					LastPushed: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
				},
			}),
			length: 3,
		},
		{
			images: ByLastPushedDesc([]Image{}),
			length: 0,
		},
	}

	for _, result := range expected {
		if result.images.Len() != result.length {
			log(t, result.images.Len(), result.length)
			t.Fail()
		}
	}
}

func TestLess(t *testing.T) {
	images := []Image{
		{
			Digest:     "three",
			LastPushed: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Digest:     "one",
			LastPushed: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Digest:     "two",
			LastPushed: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	if ByLastPushedDesc(images).Less(1, 0) {
		t.Logf("Expected %v to have been pushed before %v", images[1].LastPushed, images[0].LastPushed)
		t.Fail()
	}
}

func TestSort(t *testing.T) {
	images := []Image{
		{
			Digest:     "one",
			LastPushed: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			Digest:     "three",
			LastPushed: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Digest:     "two",
			LastPushed: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	expectedSort := []Image{
		{
			Digest:     "three",
			LastPushed: time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			Digest:     "two",
			LastPushed: time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Digest:     "one",
			LastPushed: time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	sort.Sort(ByLastPushedDesc(images))

	for i := 0; i < len(expectedSort); i++ {
		if expectedSort[i].Digest != images[i].Digest {
			t.Logf("For element %d, expected digest %s, got digest %s", i, expectedSort[i].Digest, images[i].Digest)
			t.Fail()
		}
	}
}
