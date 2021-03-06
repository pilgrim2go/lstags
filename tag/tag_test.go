package tag

import (
	"testing"

	"strconv"
)

var params = map[string]string{
	"name":   "latest",
	"digest": "sha256:c92260fe6357ac1cdd79e86e23fa287701c5edd2921d243a253fd21c9f0012ae",
}

func TestNewTag(t *testing.T) {
	i, err := New(
		params["name"],
		params["digest"],
	)

	if err != nil {
		t.Fatalf("Unable to create new valid tag instance: %s", err.Error())
	}

	if i.GetName() != params["name"] {
		t.Fatalf("Unexpected name: '%s' (expected '%s')", i.GetName(), params["name"])
	}
	if i.GetDigest() != params["digest"] {
		t.Fatalf("Unexpected digest: '%s' (expected '%s')", i.GetDigest(), params["digest"])
	}
	if i.GetShortDigest() != params["digest"][0:40] {
		t.Fatalf("Unexpected short digest: '%s' (expected '%s')", i.GetShortDigest(), params["digest"][0:40])
	}
}

func TestNewTagWithEmptyName(t *testing.T) {
	_, err := New(
		"",
		params["digest"],
	)

	if err == nil {
		t.Fatalf("Was able to create tag with empty name")
	}
}

func TestNewTagWithEmptyDigest(t *testing.T) {
	_, err := New(
		params["name"],
		"",
	)

	if err == nil {
		t.Fatalf("Was able to create tag with empty image digest")
	}
}

func remoteTags() map[string]*Tag {
	tags := make(map[string]*Tag, 0)

	tg1, _ := New(
		"latest",
		"sha256:c92260fe6357ac1cdd79e86e23fa287701c5edd2921d243a253fd21c9f0012ae",
	)
	tags["latest"] = tg1

	tg2, _ := New(
		"v1.1",
		"sha256:7abd16433f3bec5ee4c566ddbfc0e5255678498d5e7e2da8f41393bfe84bfcac",
	)
	tags["v1.1"] = tg2

	tg3, _ := New(
		"v1.2",
		"sha256:7f7f94f26d23f7aca80a33732161af068f9f62fbe0e824a58cf3a39d209cfa77",
	)
	tags["v1.2"] = tg3

	tg4, _ := New(
		"v1.3.1",
		"sha256:9fb0e8a4f629b72a0a69aef357e637e4145b6588f04f1540a31a0d2e030ea7da",
	)
	tags["v1.3.1"] = tg4

	tg5, _ := New(
		"v1.3.2",
		"sha256:fc41473fc36c09222a29ffce9eaf5732fae91c3fabfa40aa878f600e13c7fed3",
	)
	tg5.SetContainerID("16dcde7895c73c98161aa6981c4ea0df027697cd")
	tags["v1.3.2"] = tg5

	return tags
}

func localTags() map[string]*Tag {
	tags := make(map[string]*Tag, 0)

	tg1, _ := New(
		"latest",
		"sha256:8ffc20b5be0e391f07f270bf79441fbea3c8b67200e5316bdefad9e0ca80277b",
	)
	tg1.SetImageID("sha256:883e3a5b24d7b46f81436bfc85564a676aa021a2c8adedc3ac6ab12ac06fdd95")
	tags["latest"] = tg1

	tg2, _ := New(
		"v1.0",
		"sha256:fe4286e7b852dc6aad6225239ecb32691f15f20b0d4354defb4ca4957958b2f0",
	)
	tg2.SetImageID("sha256:c9a69a36ff3ce76d7970df83bd438f0f1bc0363a3b4707b42542ea20ba4282f4")
	tags["v1.0"] = tg2

	tg3, _ := New(
		"v1.2",
		"sha256:7f7f94f26d23f7aca80a33732161af068f9f62fbe0e824a58cf3a39d209cfa77",
	)
	tg3.SetImageID("sha256:4c4ebb9614ef823bd04e5eba65e59286a4314d3a063e2eaa221d38fc21723cea")
	tags["v1.2"] = tg3

	tg4, _ := New(
		"v1.3.1",
		"sha256:7264ba7450b6be1bfba9ab29f506293bb324f4764c41ff32dcc04379c1a69c91",
	)
	tags["v1.3.1"] = tg4

	tg5, _ := New(
		"v1.3.2",
		"sha256:70fbfacca0ab3ec01258b1787b02de77474c6f120b86bb8743b81b7dc37d4aab",
	)
	tg5.SetContainerID("16dcde7895c73c98161aa6981c4ea0df027697cd")
	tags["v1.3.2"] = tg5

	return tags
}

func TestJoinLength(t *testing.T) {
	const expected = 6

	_, _, tags := Join(remoteTags(), localTags())

	c := len(tags)

	if c != expected {
		t.Fatalf(
			"Unexpected number of joined tags: %s (expected: %s)",
			strconv.Itoa(c),
			strconv.Itoa(expected),
		)
	}
}

func TestJoinDigest(t *testing.T) {
	expected := map[string]string{
		"latest": "sha256:c92260fe6357ac1cdd79e86e23fa287701c5edd2921d243a253fd21c9f0012ae",
		"v1.0":   "sha256:fe4286e7b852dc6aad6225239ecb32691f15f20b0d4354defb4ca4957958b2f0",
		"v1.1":   "sha256:7abd16433f3bec5ee4c566ddbfc0e5255678498d5e7e2da8f41393bfe84bfcac",
		"v1.2":   "sha256:7f7f94f26d23f7aca80a33732161af068f9f62fbe0e824a58cf3a39d209cfa77",
	}

	_, _, tags := Join(remoteTags(), localTags())

	for name, digest := range expected {
		if tags[name].GetDigest() != digest {
			t.Fatalf(
				"Unexpected digest [%s]: %s (expected: %s)",
				name,
				tags[name].GetDigest(),
				digest,
			)
		}
	}
}

func TestJoinImageID(t *testing.T) {
	expected := map[string]string{
		"latest": "883e3a5b24d7",
		"v1.0":   "c9a69a36ff3c",
		"v1.1":   "n/a",
		"v1.2":   "4c4ebb9614ef",
	}

	_, _, tags := Join(remoteTags(), localTags())

	for name, imageID := range expected {
		if tags[name].GetImageID() != imageID {
			t.Fatalf(
				"Unexpected image ID [%s]: %s (expected: %s)",
				name,
				tags[name].GetImageID(),
				imageID,
			)
		}
	}
}

func TestJoinState(t *testing.T) {
	expected := map[string]string{
		"latest": "CHANGED",
		"v1.0":   "LOCAL-ONLY",
		"v1.1":   "ABSENT",
		"v1.2":   "PRESENT",
		"v1.3.1": "CHANGED",
		"v1.3.2": "PRESENT",
	}

	_, _, tags := Join(remoteTags(), localTags())

	for name, state := range expected {
		if tags[name].GetState() != state {
			t.Fatalf(
				"Unexpected state [%s]: %s (expected: %s)",
				name,
				tags[name].GetState(),
				state,
			)
		}
	}
}

func TestJoinTagNeedsPushWithoutPushUpdate(t *testing.T) {
	examples := map[string]bool{
		"v1.3.1": false,
		"v1.3.2": false,
	}
	_, _, tags := Join(remoteTags(), localTags())

	for name, expected := range examples {
		needsPush := tags[name].NeedsPush(false)

		if needsPush != expected {
			t.Fatalf(
				"Unexpected push need [%s]: %v (expected: %v)",
				name,
				needsPush,
				expected,
			)
		}
	}
}

func TestJoinTagNeedsPushWithPushUpdate(t *testing.T) {
	examples := map[string]bool{
		"v1.3.1": true,
		"v1.3.2": false,
	}
	_, _, tags := Join(remoteTags(), localTags())

	for name, expected := range examples {
		needsPush := tags[name].NeedsPush(true)

		if needsPush != expected {
			t.Fatalf(
				"Unexpected push need [%s]: %v (expected: %v)",
				name,
				needsPush,
				expected,
			)
		}
	}
}
