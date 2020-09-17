package lights_test

import (
	"github.com/keep94/marvin2/lights"
	"reflect"
	"testing"
)

func TestSlice(t *testing.T) {
	islice, ok := lights.All.Slice()
	if len(islice) > 0 || !ok {
		t.Error("Expected empty int slice and true.")
	}
	// Other tests are covered via the String method.
}

func TestFormatLights(t *testing.T) {
	assertStrEqual(t, "All", lights.All.String())
	assertStrEqual(t, "None", lights.None.String())
	lightSet := lights.Set{
		1: true, 2: true, 3: true, 5: true, 6: false, 8: true, 9: false}
	assertStrEqual(t, "1,2,3,5,8", lightSet.String())
}

func TestOverlapWith(t *testing.T) {
	oneThreeFive := lights.New(1, 3, 5)
	if oneThreeFive.OverlapsWith(lights.None) {
		t.Error("Can't overlap with empty set.")
	}
	if lights.None.OverlapsWith(oneThreeFive) {
		t.Error("Can't overlap with empty set.")
	}
	if !oneThreeFive.OverlapsWith(lights.All) {
		t.Error("Everything should overlap with all lights.")
	}
	if !lights.All.OverlapsWith(oneThreeFive) {
		t.Error("Everything should overlap with all lights.")
	}
	if !lights.All.OverlapsWith(lights.All) {
		t.Error("All lights should overlap with itself.")
	}
	if lights.None.OverlapsWith(lights.None) {
		t.Error("No lights should never overlap with itself.")
	}
	if lights.All.OverlapsWith(lights.None) {
		t.Error("All lights should not overlap with no lights.")
	}
	if lights.None.OverlapsWith(lights.All) {
		t.Error("All lights should not overlap with no lights.")
	}
	if oneThreeFive.OverlapsWith(lights.New(2, 4)) {
		t.Error("They don't overlap")
	}
	if !oneThreeFive.OverlapsWith(lights.New(5, 7, 9)) {
		t.Error("These should overlap")
	}
	if oneThreeFive.OverlapsWith(lights.Set{5: false}) {
		t.Error("These don't overlap")
	}
	if !oneThreeFive.OverlapsWith(lights.New(1, 7, 9, 12)) {
		t.Error("These should overlap")
	}
	if !oneThreeFive.OverlapsWith(lights.New(3)) {
		t.Error("These should overlap")
	}
}

func TestIsNoneIsAll(t *testing.T) {
	if !lights.None.IsNone() || lights.None.IsAll() {
		t.Error("No lights should have no lights")
	}
	if lights.All.IsNone() || !lights.All.IsAll() {
		t.Error("All lights should have all lights")
	}
	if !lights.New().IsNone() || lights.New().IsAll() {
		t.Error("No listed lights shouldhave no lights")
	}
	l := lights.Set{3: false, 6: false}
	if !l.IsNone() || l.IsAll() {
		t.Error("Expect light set to be none")
	}
	l = lights.Set{3: false, 6: true}
	if l.IsNone() || l.IsAll() {
		t.Error("Expect light to be neither none or all")
	}
}

func TestInvString(t *testing.T) {
	verifyInvString(t, lights.None)
	verifyInvString(t, lights.All)
	verifyInvString(t, nil)
	verifyInvString(t, lights.New(2, 1, 4, 4))
}

func TestParseLights(t *testing.T) {
	actual, err := lights.Parse("")
	if err != nil {
		t.Errorf("Got error parsing %v", err)
		return
	}
	assertLightSetEqual(t, lights.All, actual)
	actual, err = lights.Parse("9")
	if err != nil {
		t.Errorf("Got error parsing %v", err)
		return
	}
	assertLightSetEqual(t, lights.New(9), actual)
	actual, err = lights.Parse("9, 3, 9, 3, 5, 8, 2, 4, 10")
	if err != nil {
		t.Errorf("Got error parsing %v", err)
		return
	}
	assertLightSetEqual(
		t,
		lights.New(2, 3, 4, 5, 8, 9, 10),
		actual)
	_, err = lights.Parse("asdfj ksdfj")
	if err == nil {
		t.Errorf("Expected error parsing.")
	}
	_, err = lights.Parse("2, 0, 3")
	if err == nil {
		t.Errorf("Expected error parsing need positive light Ids.")
	}
}

func TestSubtract(t *testing.T) {
	ls := lights.New(1, 3, 5)
	assertStrEqual(
		t, "1,3,5", ls.Subtract(lights.None).String())
	assertStrEqual(
		t, "1,3,5", ls.Subtract(lights.New(2, 4)).String())
	assertStrEqual(
		t, "1,5", ls.Subtract(lights.New(3, 6)).String())
	assertStrEqual(
		t, "1,3,5", ls.Subtract(lights.Set{3: false}).String())
	assertStrEqual(
		t, "1,3,5", ls.Subtract(lights.None).String())
	assertStrEqual(
		t, "None", ls.Subtract(lights.All).String())
}

func TestIntersect(t *testing.T) {
	onethreefive := lights.New(1, 3, 5)
	twofour := lights.New(2, 4)
	fiveseven := lights.New(5, 7)
	assertStrEqual(
		t, "None", onethreefive.Intersect(twofour).String())
	assertStrEqual(
		t, "5", onethreefive.Intersect(fiveseven).String())
	assertStrEqual(
		t,
		"None",
		onethreefive.Intersect(lights.Set{3: false}).String())
	assertStrEqual(
		t,
		"None",
		onethreefive.Intersect(lights.None).String())
	assertStrEqual(
		t,
		"None",
		lights.None.Intersect(onethreefive).String())
	assertStrEqual(
		t,
		"None",
		lights.None.Intersect(lights.None).String())
	assertStrEqual(
		t,
		"1,3,5",
		onethreefive.Intersect(onethreefive).String())
	assertStrEqual(
		t,
		"1,3,5",
		onethreefive.Intersect(lights.All).String())
	assertStrEqual(
		t,
		"1,3,5",
		lights.All.Intersect(onethreefive).String())
	assertStrEqual(
		t,
		"All",
		lights.All.Intersect(lights.All).String())
}

func TestBuilder(t *testing.T) {
	var builder lights.Builder
	builder.AddOne(1).AddOne(3)
	lights_13 := builder.Build()
	builder.Add(lights.New(7, 9)).Add(lights.New(3, 5))
	lights_13579 := builder.Build()
	assertStrEqual(t, "1,3", lights_13.String())
	assertStrEqual(t, "1,3,5,7,9", lights_13579.String())
}

func TestEmpyBuilder(t *testing.T) {
	var builder lights.Builder
	empty := builder.Build()
	two := builder.AddOne(2).Build()
	five := builder.Clear().AddOne(5).Build()
	assertStrEqual(t, "None", empty.String())
	assertStrEqual(t, "2", two.String())
	assertStrEqual(t, "5", five.String())
}

func TestInitializedBuilder(t *testing.T) {
	origLights := lights.New(1, 2, 4)
	newLights := lights.NewBuilder(origLights).AddOne(5).Build()
	assertStrEqual(t, "1,2,4", origLights.String())
	assertStrEqual(t, "1,2,4,5", newLights.String())
}

func TestBuilderWithNil(t *testing.T) {
	var builder lights.Builder
	all := builder.AddOne(5).Add(nil).Add(lights.New(3)).AddOne(2).Build()
	allAgain := lights.NewBuilder(nil).AddOne(7).Build()
	assertStrEqual(t, "All", all.String())
	assertStrEqual(t, "All", allAgain.String())
}

func TestAdd(t *testing.T) {
	newls := lights.None.Add(
		lights.New(1, 2)).Add(lights.New(2, 3)).Add(lights.New(1, 3))
	assertStrEqual(t, "1,2,3", newls.String())
	assertStrEqual(t, "1,2,3", newls.Add(lights.Set{4: false}).String())
	assertStrEqual(t, "1,2,3", lights.None.Add(newls).String())
	assertStrEqual(t, "1,2,3", newls.Add(lights.None).String())
	assertStrEqual(t, "All", newls.Add(lights.All).String())
	assertStrEqual(t, "All", lights.All.Add(newls).String())
	assertStrEqual(t, "All", lights.All.Add(lights.All).String())
	assertStrEqual(t, "None", lights.None.Add(lights.None).String())
}

func TestMap(t *testing.T) {
	m := lights.Map{5: 6, 6: 9}
	assertIntEqual(t, 6, m.Convert(5))
	assertIntEqual(t, 9, m.Convert(6))
	assertIntEqual(t, 4, m.Convert(4))
}

func assertIntEqual(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func assertStrEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func assertLightSetEqual(t *testing.T, expected, actual lights.Set) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func verifyInvString(t *testing.T, lset lights.Set) {
	actual, err := lights.InvString(lset.String())
	if err != nil {
		t.Errorf("Got error %v", err)
		return
	}
	if !reflect.DeepEqual(lset, actual) {
		t.Errorf("Expected %v, got %v", lset, actual)
	}
}
