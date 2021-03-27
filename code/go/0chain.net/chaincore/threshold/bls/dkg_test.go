package bls

import (
	"0chain.net/core/logging"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"0chain.net/chaincore/wallet"
	"github.com/0chain/gosdk/bls"
)

type DKGID = bls.ID

type DKGKeyShareImpl struct {
	wallet *wallet.Wallet
	id     DKGID

	//This is private to the party
	Msk []bls.SecretKey // aik, coefficients of the polynomial, Si(x). ai0 = secret value of i

	//This is public information
	mpk []bls.PublicKey // Aik (Fik in some papers) = (g^aik), coefficients of the polynomial, Pi(x). Ai0 = public value of i

	//This is aggregate private info to the party from others
	sij map[bls.ID]bls.SecretKey // Sij = Si(j) for j in all parties

	//Below info is from Qual Set

	//This is arregate private info computed from sij
	si bls.SecretKey // Secret key share - Sigma(Sij) , for j in qual parties

	//This is publicly computable information for others and for self, it's the public key of the private share key
	pi *bls.PublicKey // Public key share - g^Si

	//This is the group public key
	gmpk []bls.PublicKey // Sigma(Aik) , for i in qual parties
}

type DKGSignatureShare struct {
	id        DKGID
	signature bls.Sign
}

type DKGSignature struct {
	shares []DKGSignatureShare
}

var wallets []*wallet.Wallet

var dkgShares []*DKGKeyShareImpl

// Debug-only function to set the 'msk' to custom values, for more reliable
// unit test repro.
func SetMsk(msk []bls.SecretKey, i, t int) {
	m := make(map[int][]string)

	// m[0] = []string{
	//  "08b73302a182654ba1d6c40f9a621d61d324c74d10037c262dbb0407d1e2c623",
	//   "121b2d51a81d0d81a34c3580167939b24331d2c422acd1cdab072939487f6960",
	// }
	// m[1] = []string{
	//   "0cc8bd96c78404532ab9f57cf22746b74ee95a2536795144b8ff6cd471c4837f",
	//   "157242a9337c581e4270ba26a41a94f3a5c94e7c8b7639f3cf49d4955c8485fb",
	// }
	// m[2] = []string{
	//   "18de5c3b6c738dd43579f0f921119c054cc502448615bbb2ba313f6a8eb668ac",
	//   "0bf5059875921e668a5bdf2c7fc4844592d2572bcd0668d2d6c52f5054e2d083",
	// }

	m[0] = []string{
		"c730018f37b32f82aed59d8b73c8f48718784f69b6d896fecd92fbed5facb817",
		"1c297c05e0e1f70777e2d394060ead76a3fd0a4b4dc4f1f0212b748090adf711",
	}
	m[1] = []string{
		"5612847e14bfabc04b958eee142767bd82a013e610ad0937bf53a6f187638f0a",
		"577e52f6e300814b9597341111991e7837d58c167551a4f0e35b8d3eb2991f03",
	}
	m[2] = []string{
		"b51d6d76b427634bf07edb48e09336be884c0aa3a6a0078fe4e9d31dd8e14720",
		"8db1f24e753d623fcbf2d44c4c4fd84baf391bd3e6148caab7dd69c73c33e51e",
	}

	for z := 0; z < t; z++ {
		msk[z].DeserializeHexStr(m[i][z])
	}
}

// Prints out a []string that can be used for SetMsk, for more reliable unit
// test repro.
func PrintMsk(msk []bls.SecretKey, i int) {
	fmt.Println("m[", i,"] = []string{")
	for _, msk := range msk {
		fmt.Println("\t\""+msk.SerializeToHexStr()+"\"")
	}
	fmt.Println("}")
}

//GenerateWallets - generate the wallets used to participate in DKG
func GenerateWallets(n int) {
	for i := 0; i < n; i++ {
		wallet := &wallet.Wallet{}
		if err := wallet.Initialize("bls0chain"); err != nil {
			panic(err)
		}
		wallets = append(wallets, wallet)
	}
}

//InitializeDKGShares - initialize DKG Share structures
func InitializeDKGShares(t int) {
	dkgShares = dkgShares[:0]
	for _, wallet := range wallets {
		dkgShare := &DKGKeyShareImpl{wallet: wallet}
		if err := dkgShare.id.SetHexString("1" + wallet.ClientID[:31]); err != nil {
			fmt.Printf("client id: %v\n", wallet.ClientID)
			panic(err)
		}
		dkgShare.sij = make(map[bls.ID]bls.SecretKey)
		dkgShares = append(dkgShares, dkgShare)
	}
}

//GenerateDKGKeyShare - create Si(x) and corresponding Pi(x) polynomial coefficients
func (dkgs *DKGKeyShareImpl) GenerateDKGKeyShare(t int) {
	var dsk bls.SecretKey
	dsk.SetByCSPRNG()
	dkgs.Msk = dsk.GetMasterSecretKey(t)
	dkgs.mpk = bls.GetMasterPublicKey(dkgs.Msk)
}

//GenerateSij - generate secret key shares from i for each party j
func (dkgs *DKGKeyShareImpl) GenerateSij(ids []DKGID) {
	dkgs.sij = make(map[bls.ID]bls.SecretKey)
	for _, id := range ids {
		var sij bls.SecretKey
		if err := sij.Set(dkgs.Msk, &id); err != nil {
			panic(err)
		}
		dkgs.sij[id] = sij
	}
}

//ValidateShare - validate Sij using Pj coefficients
func (dkgs *DKGKeyShareImpl) ValidateShare(jpk []bls.PublicKey, sij bls.SecretKey) bool {
	var expectedSijPK bls.PublicKey
	if err := expectedSijPK.Set(jpk, &dkgs.id); err != nil {
		panic(err)
	}
	sijPK := sij.GetPublicKey()
	return expectedSijPK.IsEqual(sijPK)
}

//AggregateSecretShares - compute Si = Sigma(Sij), j in qual and Pi = g^Si
//Useful to compute self secret key share and associated public key share
//For other parties, the public key share can be derived using the Pj(x) coefficients
func (dkgs *DKGKeyShareImpl) AggregateSecretKeyShares(qual []DKGID, dkgShares map[bls.ID]*DKGKeyShareImpl) {
	sk := bls.NewSecretKey()
	for _, id := range qual {
		dkgsj, ok := dkgShares[id]
		if !ok {
			panic("no share")
		}
		sij := dkgsj.sij[dkgs.id]
		sk.Modadd(&sij)
	}
	dkgs.si = *sk
	dkgs.pi = dkgs.si.GetPublicKey()
}

//ComputePublicKeyShare - compute the public key share of any party j, based on the coefficients of Pj(x)
func (dkgs *DKGKeyShareImpl) ComputePublicKeyShare(qual []DKGID, dkgShares map[bls.ID]*DKGKeyShareImpl) bls.PublicKey {
	pk := bls.NewPublicKey()
	for _, id := range qual {
		dkgsj, ok := dkgShares[id]
		if !ok {
			panic("no share")
		}
		var pkj bls.PublicKey
		pkj.Set(dkgsj.mpk, &dkgs.id)
		pk.Add(&pkj)
	}
	return *pk
}

//AggregatePublicKeyShares - compute Sigma(Aik, i in qual)
func (dkgs *DKGKeyShareImpl) AggregatePublicKeyShares(qual []DKGID, dkgShares map[bls.ID]*DKGKeyShareImpl) {
	dkgs.gmpk = dkgs.gmpk[:0]
	for k := 0; k < len(dkgs.mpk); k++ {
		pk := bls.NewPublicKey()
		for _, id := range qual {
			dkgsj, ok := dkgShares[id]
			if !ok {
				panic("no share")
			}
			pk.Add(&dkgsj.mpk[k])
		}
		dkgs.gmpk = append(dkgs.gmpk, *pk)
	}
}

//Sign - sign using the group secret key share
func (dkgs *DKGKeyShareImpl) Sign(msg string) string {
	return dkgs.si.Sign([]byte(msg)).SerializeToHexStr()
}

//VerifySignature - verify the signature using the group public key share
func (dkgs *DKGKeyShareImpl) VerifySignature(msg string, sig *bls.Sign) bool {
	return sig.Verify(dkgs.pi, []byte(msg))
}

//VerifyGroupSignature - verify group signature using group public key
func (dkgs *DKGKeyShareImpl) VerifyGroupSignature(msg string, sig *bls.Sign) bool {
	return sig.Verify(&dkgs.gmpk[0], []byte(msg))
}

//Recover - given t signature shares, recover the group signature (using lagrange interpolation)
func (dkgs *DKGKeyShareImpl) Recover(dkgSigShares []DKGSignatureShare) (*bls.Sign, error) {
	var aggSig bls.Sign
	var signatures []Sign
	var ids []bls.ID
	t := len(dkgSigShares)
	if t > len(dkgs.Msk) {
		t = len(dkgs.Msk)
	}
	for k := 0; k < t; k++ {
		ids = append(ids, dkgSigShares[k].id)
		signatures = append(signatures, dkgSigShares[k].signature)
	}
	if err := aggSig.Recover(signatures, ids); err != nil {
		return nil, err
	}
	return &aggSig, nil
}

func init() {
	logging.InitLogging("development")
}

func TestGenerateDKG(tt *testing.T) {
	n := 3                                 //total participants at the beginning
	t := int(math.Round(0.67 * float64(n))) // threshold number of parties required to create aggregate signature
	q := int(math.Round(0.85 * float64(n))) // qualified to compute dkg based on DKG protocol execution
	if q == t && t < n {
		q++
	}
	fmt.Printf("n=%v t=%v q=%v\n", n, t, q)
	start := time.Now()
	GenerateWallets(n)
	fmt.Printf("time to generate wallets: %v\n", time.Since(start))
	InitializeDKGShares(t)
	fmt.Printf("time to initialize dkg shares: %v\n", time.Since(start))

	var ids []DKGID
	var qualIDs []DKGID
	var qualDKGSharesMap map[bls.ID]*DKGKeyShareImpl

	for _, dkgs := range dkgShares {
		ids = append(ids, dkgs.id)
	}

	//Generate aik for each party (the polynomial for sharing the secret)
	for _, dkgs := range dkgShares {
		dkgs.GenerateDKGKeyShare(t)
	}
	fmt.Printf("time to generate dkg key shares: %v\n", time.Since(start))

	//Generate Sij for each party (the p(id) value for a given id)
	for _, dkgs := range dkgShares {
		dkgs.GenerateSij(ids)
	}
	fmt.Printf("time to generate dkg key share sij: %v\n", time.Since(start))

	//Validate Sij shares received from others using P(x)
	for _, dkgsi := range dkgShares {
		for _, dkgsj := range dkgShares {
			sij := dkgsj.sij[dkgsi.id]
			valid := dkgsi.ValidateShare(dkgsj.mpk, sij)
			if !valid {
				fmt.Printf("%v -> %v share valid = %v\n", dkgsi.wallet.ClientID[:7], dkgsj.wallet.ClientID[:7], valid)
				tt.Fatal("Should have been a valid share.")
			}
		}
		//fmt.Printf("time to dkg key share validate: %v\n", time.Since(start))
	}
	fmt.Printf("time to dkg key share validate all: %v\n", time.Since(start))

	//Simulate Qual Set
	shuffled := make([]*DKGKeyShareImpl, n)
	perm := rand.New(rand.NewSource(int64(time.Now().Nanosecond()))).Perm(len(shuffled))
	for i, v := range perm {
		shuffled[v] = dkgShares[i]
	}
	qualDKGSharesMap = make(map[bls.ID]*DKGKeyShareImpl)
	for i := 0; i < q; i++ {
		qualIDs = append(qualIDs, shuffled[i].id)
		qualDKGSharesMap[shuffled[i].id] = shuffled[i]
	}

	//Compute si = Sigma sij, aggregate secret share of each qualified party
	for _, dkgsi := range qualDKGSharesMap {
		dkgsi.AggregateSecretKeyShares(qualIDs, qualDKGSharesMap)
		dkgsi.AggregatePublicKeyShares(qualIDs, qualDKGSharesMap)
	}
	fmt.Printf("time to aggregate secret/public key shares: %v\n", time.Since(start))
	for _, dkgsi := range qualDKGSharesMap {
		for _, dkgsj := range qualDKGSharesMap {
			if dkgsi == dkgsj {
				continue
			}
			pk := dkgsi.ComputePublicKeyShare(qualIDs, qualDKGSharesMap)
			if !pk.IsEqual(dkgsi.pi) {
				panic("public key share not valid")
			}
		}
	}
	fmt.Printf("time to compute public key share: %v\n", time.Since(start))

	msg := fmt.Sprintf("Hello 0Chain World %v", time.Now())
	falseMsg := fmt.Sprintf("Hello 0Chain World %v", time.Now())
	var falseCount int
	var signatures []DKGSignatureShare

	//Sign a message
	prng := rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	for idx, id := range qualIDs {
		dkgsi, ok := qualDKGSharesMap[id]
		if !ok {
			panic(fmt.Sprintf("no share: %v\n", idx))
		}
		var sign string
		//if rand.Float64() < float64(t)/float64(n) {
		if prng.Float64() < float64(t)/float64(q) {
			sign = dkgsi.Sign(msg)
		} else {
			sign = dkgsi.Sign(falseMsg)
			falseCount++
		}
		var blsSig bls.Sign
		if err := blsSig.DeserializeHexStr(sign); err != nil {
			panic(err)
		}
		signature := DKGSignatureShare{signature: blsSig, id: dkgsi.id}
		signatures = append(signatures, signature)
	}
	fmt.Printf("time to sign: %v\n", time.Since(start))
	fmt.Printf("Signatures: Correct: %v False: %v Total:%v\n", q-falseCount, falseCount, q)
	//Aggregate Signatures
	count := 0
	for _, id := range qualIDs {
		count++
		dkgsi, ok := qualDKGSharesMap[id]
		if !ok {
			panic("no share")
		}
		var dkgSignature DKGSignature
		for _, signature := range signatures {
			if signature.id != dkgsi.id {
				if prng.Float64() < 0.10 {
					//To simulate network/byzantine condition of not getting the shares
					continue
				}
				dkgsj, ok := qualDKGSharesMap[signature.id]
				if !ok {
					panic("no share")
				}
				if !dkgsj.VerifySignature(msg, &signature.signature) {
					fmt.Printf("\tInvalid signature from: %v\n", signature.id)
					continue
				}
			}
			dkgSignature.shares = append(dkgSignature.shares, signature)
		}
		if len(dkgSignature.shares) < t {
			fmt.Printf("signature %v %v(%v): insufficient signature shares\n", count, dkgsi.wallet.ClientID[:7], dkgsi.id)
			continue
		}
		shuffled := make([]DKGSignatureShare, len(dkgSignature.shares))
		perm := rand.New(rand.NewSource(int64(time.Now().Nanosecond()))).Perm(len(shuffled))
		//To simulate network condition and also self not having the share yet
		for i, v := range perm {
			shuffled[v] = dkgSignature.shares[i]
		}
		asign, err := dkgsi.Recover(shuffled)
		if err != nil {
			fmt.Printf("Error recovering signature %v\n", err)
			continue
		}
		gsigValid := dkgsi.VerifyGroupSignature(msg, asign)
		fmt.Printf("signature %v %v(%v): %v %v\n", count, dkgsi.wallet.ClientID[:7], dkgsi.id, asign.SerializeToHexStr()[:32], gsigValid)
	}
	fmt.Printf("time to finish: %v\n", time.Since(start))
}
