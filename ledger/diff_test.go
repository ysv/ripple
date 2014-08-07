package ledger

import (
	"testing"

	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/storage/memdb"
	. "launchpad.net/gocheck"
)

func Test(t *testing.T) { TestingT(t) }

type DiffSuite struct {
	db *memdb.MemoryDB
}

var _ = Suite(&DiffSuite{})

func (s *DiffSuite) SetUpSuite(c *C) {
	var err error
	s.db, err = memdb.NewMemoryDB([]string{"testdata/38129-32570.gz", "testdata/99943.gz"})
	c.Assert(err, IsNil)
}

var expectedDiff = []string{
	"D,Account Node,0,AF47E9E91A41621B0F8AC5A119A5AD8B9E892147381BEAF6F2186127B89A44FF",
	"A,Account Node,0,2C23D15B6B549123FB351E4B5CDE81C564318EB845449CD43C3EA7953C4DB452",
	"D,Account Node,1,271E1B9B1B1FB8C7FD860F8A3CCC4F76952B8A35625933433A7420802D62D456",
	"A,Account Node,1,29FD2F34869B2E46EA2FC996FE7CB94AF4C3B40CD9859232D682F8AE1C17DAD5",
	"D,Account Node,1,724CA5CAEB55D7948CAC67BF087E7F9B0B83587B6EEB7CE186055C860B076356",
	"A,Account Node,1,067A065323B98104D6A3CAA82FE77FEDB228F10BEF5E2AAD216608B424C3CC1D",
	"D,AccountRoot,2,2A75953DB729CC20A0AD2F585D95198059FD7851282035D7E1F4B53178297F93",
	"A,AccountRoot,2,C97303390D8FF71BC716710CC2D74B209C990D18E293736E2FB87383B8F914EF",
	"D,Account Node,2,C62002973CAB176FC006FDC0EA2DF873FB57CE6B77A80D1AE7B20074D42E061A",
	"A,Account Node,2,FFB92B95013668AEED6A77545E084F4117289103EF0D5384AD7BA03775FFD4B6",
	"D,LedgerHashes,3,28372696D26A27D4E7AB82096D129EC44685BF9977257439609018353B9174AD",
	"A,LedgerHashes,3,0E386F1549BD2B10BC53DDEEA1067AE55E6F804B84347E4724DE33F75C4A9726",
}

func (s *DiffSuite) TestDiff(c *C) {
	first, err := data.NewHash256("2C23D15B6B549123FB351E4B5CDE81C564318EB845449CD43C3EA7953C4DB452") // 38,129 Account Hash
	c.Assert(err, IsNil)
	second, err := data.NewHash256("AF47E9E91A41621B0F8AC5A119A5AD8B9E892147381BEAF6F2186127B89A44FF") // 38,128 Account Hash
	c.Assert(err, IsNil)
	diff, err := Diff(*first, *second, s.db)
	c.Assert(err, IsNil)
	c.Assert(diff.String(), DeepEquals, expectedDiff)
}

var expectedSummary = "1,1,0,0,0,0,0,0,0,145,137,65,0,2,4,53,0"

func (s *DiffSuite) TestSummary(c *C) {
	ledger, err := data.NewHash256("E6DB7365949BF9814D76BCC730B01818EB9136A89DB224F3F9F5AAE4569D758E") // 38,129 Ledger Hash
	c.Assert(err, IsNil)
	state, err := NewLedgerStateFromDB(*ledger, s.db)
	c.Assert(err, IsNil)
	c.Assert(state.Fill(), IsNil)
	summary, err := state.Summary()
	c.Assert(err, IsNil)
	c.Assert(summary, DeepEquals, expectedSummary)
}

var expectedUnfolded = []string{
	"D,AccountRoot,3,2EED41E7965F3D382A84A7238C12FEEF2DB3B52C539E8EA0117D2D322D98934A",
	"A,Account Node,3,65DDB4B5447ACD15A9DBF9418B9894B8CE4F7C52B59BBA75883CD2C160C23D3D",
	"A,RippleState,4,0D625AA28BA84F2EBC862493B5D0285321F2812D209CCE2B7801EDC776ED23F5",
	"A,AccountRoot,4,2EED41E7965F3D382A84A7238C12FEEF2DB3B52C539E8EA0117D2D322D98934A",
	"D,Account Node,1,B1ABCFAE4B9F907A9B97B6F88421F53FE3E70025825605EA55CC6B5C454DB745",
	"A,Account Node,1,4980EA45C8307D855918381F759012AF96B2C48596758A12CE882061FE0C6921",
}

var expectedFolded = []string{
	"D,AccountRoot,3,2EED41E7965F3D382A84A7238C12FEEF2DB3B52C539E8EA0117D2D322D98934A",
	"A,Account Node,3,65DDB4B5447ACD15A9DBF9418B9894B8CE4F7C52B59BBA75883CD2C160C23D3D",
	"A,RippleState,4,0D625AA28BA84F2EBC862493B5D0285321F2812D209CCE2B7801EDC776ED23F5",
	"A,AccountRoot,4,2EED41E7965F3D382A84A7238C12FEEF2DB3B52C539E8EA0117D2D322D98934A",
	"D,Account Node,1,B1ABCFAE4B9F907A9B97B6F88421F53FE3E70025825605EA55CC6B5C454DB745",
	"A,Account Node,1,4980EA45C8307D855918381F759012AF96B2C48596758A12CE882061FE0C6921",
}

func (s *DiffSuite) TestFold(c *C) {
	first, err := data.NewHash256("71E4F1363337A5A0305FEE7175EBACA458381AB77E820738336E084BE190B8C6") //99,943 account hash
	c.Assert(err, IsNil)
	second, err := data.NewHash256("D027189E2A84B636C6B812FE0FA808279212D8F59B2465C9BBF9E2665294D09A") //99,942 account hash
	diff, err := Diff(*first, *second, s.db)
	c.Assert(err, IsNil)
	c.Assert(diff.String(), DeepEquals, expectedUnfolded)

	// c.Assert(state.Fill(), IsNil)
	// c.Assert(summary, DeepEquals, expectedSummary)
}
