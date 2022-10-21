package solvers

import (
	"testing"

	"github.com/tgrindinger/go-graphql-3sat-solver/graph/factories"
)

func TestGeneratePopulationLessThanMax(t *testing.T) {
	cases := []struct {
		desc string
		input []string
		want population
	}{
		{ "single variable yields two members", []string{ "var" }, population{
				{ "var": true },
				{ "var": false },
			},
		},
		{ "two variables yields four members", []string{ "var1", "var2" }, population{
				{ "var1": true, "var2": true },
				{ "var1": false, "var2": false },
				{ "var1": true, "var2": false },
				{ "var1": false, "var2": true },
			},
		},
		{ "three variables yields eight members", []string{ "var1", "var2", "var3" }, population{
				{ "var1": true, "var2": true, "var3": true },
				{ "var1": true, "var2": true, "var3": false },
				{ "var1": true, "var2": false, "var3": false },
				{ "var1": true, "var2": false, "var3": true },
				{ "var1": false, "var2": true, "var3": true },
				{ "var1": false, "var2": true, "var3": false },
				{ "var1": false, "var2": false, "var3": true },
				{ "var1": false, "var2": false, "var3": false },
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			maxPopulation := 10
			randomFactory := &factories.ZeroRandomFactory{}
			sut := &PopulationGenerator{randomFactory: randomFactory}

			// act
			got := sut.generatePopulation(maxPopulation, tc.input)

			// assert
			assertPopulationsAreEqual(t, got, tc.want)
		})
	}
}

func TestGeneratePopulationOverMax(t *testing.T) {
	cases := []struct {
		desc string
		input []string
	}{
		{ "four variables yields max pop size", []string{ "var1", "var2", "var3", "var4" } },
	}
	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			// arrange
			maxPopulation := 10
			randomFactory := &factories.ZeroRandomFactory{}
			sut := &PopulationGenerator{randomFactory: randomFactory}

			// act
			got := sut.generatePopulation(maxPopulation, tc.input)

			// assert
			assertPopulationHasCase(t, got, map[string]bool{ "var1": true, "var2": true, "var3": true, "var4": true })
			assertPopulationHasCase(t, got, map[string]bool{ "var1": false, "var2": false, "var3": false, "var4": false })
			assertMembersAreUnique(t, got)
		})
	}
}

func assertPopulationHasCase(t testing.TB, got population, member map[string]bool) {
	for _, m := range got {
		found := true
		for k, v := range m {
			if v != member[k] {
				found = false
			}
		}
		if found {
			return
		}
	}
	t.Fatalf("Unable to find member in population: %v", member)
}

func assertMembersAreUnique(t testing.TB, got population) {
	for _, m := range got {
		assertMemberOccursOnce(t, m, got)
	}
}

func assertMemberOccursOnce(t testing.TB, member member, population population) {
	found := false
	for _, m := range population {
		if member.matches(m) {
			if found {
				t.Fatalf("Found multiple occurrences of %v", member)
			}
			found = true
		}
	}
}

func assertPopulationsAreEqual(t testing.TB, got population, want population) {
	if len(got) != len(want) {
		t.Fatalf("incorrect population size: got %d want %d", len(got), len(want))
	}
	for _, member := range got {
		assertMemberIsInPopulation(t, member, want)
	}
}

func assertMemberIsInPopulation(t testing.TB, member member, pop population) {
	for _, popMember := range pop {
		if member.matches(popMember) {
			return
		}
	}
	t.Fatalf("unable to find member in population: %v", member)
}
