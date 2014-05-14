package hashtree

import (
	"fmt"
	"io"
	"testing"
)

type htIOTest struct {
	out string
	in  string
}

var testdata = []htIOTest{
	{"0000000000000000000000000000000000000000000000000000000000000000",
		""},
	{"6100000000000000000000000000000000000000000000000000000000000000",
		"a"},
	{"6162000000000000000000000000000000000000000000000000000000000000",
		"ab"},
	{"6162630000000000000000000000000000000000000000000000000000000000",
		"abc"},
	{"6162636400000000000000000000000000000000000000000000000000000000",
		"abcd"},
	{"6162636465000000000000000000000000000000000000000000000000000000",
		"abcde"},
	{"6162636465660000000000000000000000000000000000000000000000000000",
		"abcdef"},
	{"6162636465666700000000000000000000000000000000000000000000000000",
		"abcdefg"},
	{"6162636465666768000000000000000000000000000000000000000000000000",
		"abcdefgh"},
	{"6162636465666768690000000000000000000000000000000000000000000000",
		"abcdefghi"},
	{"6162636465666768696a00000000000000000000000000000000000000000000",
		"abcdefghij"},
	{"1048275e1635b814b1cafdf96d8c8e2e8ff8210140d8a757e42b11a82e6993f1",
		"Discard medicine more than two years old."},
	{"a775fbc3ba5f2399f06c5487b80ec40cf1355830a8d2009eca137326fdaf057e",
		"He who has a shady past knows that nice guys finish last."},
	{"3d0accd175b46285a7b6e7cd1c50022efcce5ebcbd9dd61570cb8c5a2e605812",
		"I wouldn't marry him with a ten foot pole."},
	{"d2d4a59a1e51e7d330dc07ac04eb9522cc2e06977b22877451c465501c1e1d78",
		"Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{"4d548d425a0c4092dad521b4c733aa09d9e288b07c8cfd81251c12e2711dd2c8",
		"The days of the digital watch are numbered.  -Tom Stoppard"},
	{"4e6570616c207072656d69657220776f6e27742072657369676e2e0000000000",
		"Nepal premier won't resign."},
	{"82a44cd00d54bdfcd036fe52b20acc397b91cf32308e3d6c4489b4894d5ed9f1",
		"For every action there is an equal and opposite government program."},
	{"99a989c3ff4ce2f59426fd86a4d7d6035021ebb2e7ef330cb0127c4f25c62bb6",
		"His money is twice tainted: 'taint yours and 'taint mine."},
	{"e374e562ad7c95cdf5aa590da427196ae3dc14d323dde7a9e1b4655995b5f836",
		"There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{"8f6a8a24db1caafc13baa8b036f1767843f63b90487e50743df4f70b159dcb3f",
		"It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{"73697a653a2020612e6f75743a2020626164206d616769630000000000000000",
		"size:  a.out:  bad magic"},
	{"4288da1a6597d266ed35a2035a59011a5db0bac167db61ece553e212e8608064",
		"The major problem is with sendmail.  -Mark Horton"},
	{"938d25e487d153d557b606371eccda8d81fe2d3c03d43e9c27c3c3fdeae63aac",
		"Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{"f73067f72ff83d5b52832208eff65c1b177afec54268f81161f5ece048079e9c",
		"If the enemy is within range, then so are you."},
	{"c547f62e5f9b695e8f60814f1b9ad1771637ffa3cd334d9cc14fa452d9c343c4",
		"It's well we cannot hear the screams/That we create in others' dreams."},
	{"e625a1641127ea46a7e262c294f87185cd104d33a5f3587724cdc7d95953d30a",
		"You remind me of a TV show, but that's all right: I watch it anyway."},
	{"4320697320617320706f727461626c652061732053746f6e6568656467652121",
		"C is as portable as Stonehedge!!"},
	{"9bfa904681c4dc0b4a9f6649672b40366fac3bcd44974225b8373296fd0a3381",
		"Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{"c0009625c65c6e3dda048cc27675d598c9c98fa817f6409bd0a9ebb1670b26a8",
		"The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{"f40c942da95825e293842d27ec4d103177f435ae57e51de14b43b598617c12e8",
		"How can you write a big system without C++?  -Paul Glick"},
}

//test data generator
func printHashTestData() {
	c := NewTree()
	for i := 0; i < len(testdata); i++ {
		g := testdata[i]
		io.WriteString(c, g.in)
		s := fmt.Sprintf("%x", c.Sum(nil))
		fmt.Printf("{\"%v\",\n\"%v\"},\n", s, g.in)
		c.Reset()
	}
}

func TestHT(t *testing.T) {
	c := NewTree()
	for i := 0; i < len(testdata); i++ {
		g := testdata[i]
		io.WriteString(c, g.in)
		s := fmt.Sprintf("%x", c.Sum(nil))
		if s != g.out {
			t.Fatalf("testdata[%d](%s) = %s want %s", i, g.in, s, g.out)
		}
		c.Reset()
	}
}
