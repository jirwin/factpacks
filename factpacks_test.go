package factpacks

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFactStore_LoadFactPack(t *testing.T) {
	fs := MakeFactStore()

	err := fs.LoadFactPack("examples/facts.txt")
	require.NoError(t, err)

	require.Equal(t, "roses are red", fs.GetFact("roses").Output())
	require.Equal(t, "violets are blue", fs.GetFact("violets").Output())
	require.Equal(t, "one is fish", fs.GetFact("one").Output())
	require.Equal(t, "red fish is blue fish", fs.GetFact("red fish").Output())
	require.Nil(t, fs.GetFact("elephant"))
	require.Equal(t, "monkey is a banana eater => yum", fs.GetFact("monkey").Output())
}

func TestFactStore_SetFact(t *testing.T) {
	fs := MakeFactStore()

	fs.SetFact(&Fact{
		name:     "test fact",
		value:    "these are fact details",
		isPlural: false,
	})
	require.NotNil(t, fs.GetFact("test fact"))
	require.Equal(t, "test fact is these are fact details", fs.GetFact("test fact").Output())
}

func TestFactStore_DeleteFact(t *testing.T) {
	fs := MakeFactStore()

	fs.SetFact(&Fact{
		name:     "42",
		value:    "life, universe, and everything",
		isPlural: false,
	})
	require.NotNil(t, fs.GetFact("42"))
	require.Equal(t, "42 is life, universe, and everything", fs.GetFact("42").Output())

	fs.DeleteFact("42")
	require.Nil(t, fs.GetFact("test fact"))
}

func TestFactStore_HumanFactSet(t *testing.T) {
	fs := MakeFactStore()

	fs.HumanFactSet("roses are red")
	require.Equal(t, "roses are red", fs.GetFact("roses").Output())

	fs.HumanFactSet("the quick brown fox is jumping over the lazy dog")
	require.NotNil(t, fs.GetFact("the quick brown fox"))
	require.Equal(t, "the quick brown fox is jumping over the lazy dog", fs.GetFact("the quick brown fox").Output())

	fs.HumanFactSet("42 is the answer to life, the universe, and everything")
	require.NotNil(t, fs.GetFact("42"))
	require.Equal(t, "42 is the answer to life, the universe, and everything", fs.GetFact("42").Output())

	fs.HumanFactSet("monkeys are animals that live in trees and are animals with tails")
	require.NotNil(t, fs.GetFact("monkeys"))
	require.Equal(t, "monkeys are animals that live in trees and are animals with tails", fs.GetFact("monkeys").Output())
}
