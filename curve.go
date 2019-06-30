package bls12

import (
	"math/big"
)

const (
	compressedFormMask  uint8 = 1 << 7
	pointAtInfinityMask uint8 = 1 << 6
)

var (
	fqCurveB, _                   = new(fq).SetString("4")
	fqCurveBPlusOne, _            = new(fq).SetString("5")
	fqSqrtNegThree, _             = new(fq).SetString("1586958781458431025242759403266842894121773480562120986020912974854563298150952611241517463240701")
	fqHalfSqrtNegThreeMinusOne, _ = new(fq).SetString("793479390729215512621379701633421447060886740281060493010456487427281649075476305620758731620350")

	iso11XNum = []*fq{
		// 2712959285290305970661081772124144179193819192423276218370281158706191519995889425075952244140278856085036081760695
		&fq{},
		// 3564859427549639835253027846704205725951033235539816243131874237388832081954622352624080767121604606753339903542203
		&fq{},
		// 2051387046688339481714726479723076305756384619135044672831882917686431912682625619320120082313093891743187631791280
		&fq{},
		// 3612713941521031012780325893181011392520079402153354595775735142359240110423346445050803899623018402874731133626465
		&fq{},
		// 2247053637822768981792833880270996398470828564809439728372634811976089874056583714987807553397615562273407692740057
		&fq{},
		// 3415427104483187489859740871640064348492611444552862448295571438270821994900526625562705192993481400731539293415811
		&fq{},
		// 2067521456483432583860405634125513059912765526223015704616050604591207046392807563217109432457129564962571408764292
		&fq{},
		// 3650721292069012982822225637849018828271936405382082649291891245623305084633066170122780668657208923883092359301262
		&fq{},
		// 1239271775787030039269460763652455868148971086016832054354147730155061349388626624328773377658494412538595239256855
		&fq{},
		// 3479374185711034293956731583912244564891370843071137483962415222733470401948838363051960066766720884717833231600798
		&fq{},
		// 2492756312273161536685660027440158956721981129429869601638362407515627529461742974364729223659746272460004902959995
		&fq{},
		// 1058488477413994682556770863004536636444795456512795473806825292198091015005841418695586811009326456605062948114985
		&fq{},
	}

	iso11XDen = []*fq{
		// 1353092447850172218905095041059784486169131709710991428415161466575141675351394082965234118340787683181925558786844
		&fq{},
		//2822220997908397120956501031591772354860004534930174057793539372552395729721474912921980407622851861692773516917759
		&fq{},
		// 1717937747208385987946072944131378949849282930538642983149296304709633281382731764122371874602115081850953846504985
		&fq{},
		// 501624051089734157816582944025690868317536915684467868346388760435016044027032505306995281054569109955275640941784
		&fq{},
		// 3025903087998593826923738290305187197829899948335370692927241015584233559365859980023579293766193297662657497834014
		&fq{},
		// 2224140216975189437834161136818943039444741035168992629437640302964164227138031844090123490881551522278632040105125
		&fq{},
		// 1146414465848284837484508420047674663876992808692209238763293935905506532411661921697047880549716175045414621825594
		&fq{},
		// 3179090966864399634396993677377903383656908036827452986467581478509513058347781039562481806409014718357094150199902
		&fq{},
		// 1549317016540628014674302140786462938410429359529923207442151939696344988707002602944342203885692366490121021806145
		&fq{},
		// 1442797143427491432630626390066422021593505165588630398337491100088557278058060064930663878153124164818522816175370
		&fq{},
		// 1
		&fq{},
	}

	iso11YNum = []*fq{
		// 1393399195776646641963150658816615410692049723305861307490980409834842911816308830479576739332720113414154429643571
		&fq{},
		// 2968610969752762946134106091152102846225411740689724909058016729455736597929366401532929068084731548131227395540630
		&fq{},
		// 122933100683284845219599644396874530871261396084070222155796123161881094323788483360414289333111221370374027338230
		&fq{},
		// 303251954782077855462083823228569901064301365507057490567314302006681283228886645653148231378803311079384246777035
		&fq{},
		// 1353972356724735644398279028378555627591260676383150667237975415318226973994509601413730187583692624416197017403099
		&fq{},
		// 3443977503653895028417260979421240655844034880950251104724609885224259484262346958661845148165419691583810082940400
		&fq{},
		// 718493410301850496156792713845282235942975872282052335612908458061560958159410402177452633054233549648465863759602
		&fq{},
		// 1466864076415884313141727877156167508644960317046160398342634861648153052436926062434809922037623519108138661903145
		&fq{},
		// 1536886493137106337339531461344158973554574987550750910027365237255347020572858445054025958480906372033954157667719
		&fq{},
		// 2171468288973248519912068884667133903101171670397991979582205855298465414047741472281361964966463442016062407908400
		&fq{},
		// 3915937073730221072189646057898966011292434045388986394373682715266664498392389619761133407846638689998746172899634
		&fq{},
		// 3802409194827407598156407709510350851173404795262202653149767739163117554648574333789388883640862266596657730112910
		&fq{},
		// 1707589313757812493102695021134258021969283151093981498394095062397393499601961942449581422761005023512037430861560
		&fq{},
		// 349697005987545415860583335313370109325490073856352967581197273584891698473628451945217286148025358795756956811571
		&fq{},
		// 885704436476567581377743161796735879083481447641210566405057346859953524538988296201011389016649354976986251207243
		&fq{},
		// 3370924952219000111210625390420697640496067348723987858345031683392215988129398381698161406651860675722373763741188
		&fq{},
	}

	iso11YDen = []*fq{
		// 3396434800020507717552209507749485772788165484415495716688989613875369612529138640646200921379825018840894888371137
		&fq{},
		// 3907278185868397906991868466757978732688957419873771881240086730384895060595583602347317992689443299391009456758845
		&fq{},
		// 854914566454823955479427412036002165304466268547334760894270240966182605542146252771872707010378658178126128834546
		&fq{},
		// 3496628876382137961119423566187258795236027183112131017519536056628828830323846696121917502443333849318934945158166
		&fq{},
		// 1828256966233331991927609917644344011503610008134915752990581590799656305331275863706710232159635159092657073225757
		&fq{},
		// 1362317127649143894542621413133849052553333099883364300946623208643344298804722863920546222860227051989127113848748
		&fq{},
		// 3443845896188810583748698342858554856823966611538932245284665132724280883115455093457486044009395063504744802318172
		&fq{},
		// 3484671274283470572728732863557945897902920439975203610275006103818288159899345245633896492713412187296754791689945
		&fq{},
		// 3755735109429418587065437067067640634211015783636675372165599470771975919172394156249639331555277748466603540045130
		&fq{},
		// 3459661102222301807083870307127272890283709299202626530836335779816726101522661683404130556379097384249447658110805
		&fq{},
		// 742483168411032072323733249644347333168432665415341249073150659015707795549260947228694495111018381111866512337576
		&fq{},
		// 1662231279858095762833829698537304807741442669992646287950513237989158777254081548205552083108208170765474149568658
		&fq{},
		// 1668238650112823419388205992952852912407572045257706138925379268508860023191233729074751042562151098884528280913356
		&fq{},
		// 369162719928976119195087327055926326601627748362769544198813069133429557026740823593067700396825489145575282378487
		&fq{},
		// 2164195715141237148945939585099633032390257748382945597506236650132835917087090097395995817229686247227784224263055
		&fq{},
		// 1
		&fq{},
	}

	// Values taken from the execution of https://eprint.iacr.org/2019/403.pdf - A The isogeny maps.
	iso11K = [][]*fq{iso11XNum, iso11XDen, iso11YNum, iso11YDen}
)

// curvePoint is an elliptic curve point in projective coordinates. The elliptic
// curve is defined by the following equation y²=x³+3.
type curvePoint struct {
	x, y, z fq
}

// Set sets c to the value of a and returns c.
func (c *curvePoint) Set(a *curvePoint) *curvePoint {
	c.x, c.y, c.z = a.x, a.y, a.z
	return c
}

// Equal reports whether a is equal to b.
func (a *curvePoint) Equal(b *curvePoint) bool {
	return *a == *b
}

// IsInfinity reports whether the point is at infinity.
func (a *curvePoint) IsInfinity() bool {
	return a.z == fq{}
}

// Add sets c to the sum a+b and returns c.
func (c *curvePoint) Add(a, b *curvePoint) *curvePoint {
	if a.IsInfinity() {
		return c.Set(b)
	}
	if b.IsInfinity() {
		return c.Set(a)
	}

	// faster than Add
	if a.Equal(b) {
		return c.Double(a)
	}

	// See https://hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#addition-add-2007-bl
	z1z1, z2z2 := new(fq), new(fq)
	fqMul(z1z1, &a.z, &a.z)
	fqMul(z2z2, &b.z, &b.z)

	u1, u2 := new(fq), new(fq)
	fqMul(u1, &a.x, z2z2)
	fqMul(u2, &b.x, z1z1)

	s1, s2 := new(fq), new(fq)
	fqMul(s1, &a.y, &b.z)
	fqMul(s1, s1, z2z2)
	fqMul(s2, &b.y, &a.z)
	fqMul(s2, s2, z1z1)

	h, i, j, r, v := new(fq), new(fq), new(fq), new(fq), new(fq)
	fqSub(h, u2, u1)
	fqAdd(i, h, h)
	fqMul(i, i, i)
	fqMul(j, h, i)
	fqSub(r, s2, s1)
	fqAdd(r, r, r)
	fqMul(v, u1, i)

	p, t0, t1 := new(curvePoint), new(fq), new(fq)
	fqAdd(t0, v, v)
	fqAdd(t0, t0, j)
	fqMul(&p.x, r, r)
	fqSub(&p.x, &p.x, t0)

	fqAdd(t0, s1, s1)
	fqMul(t0, t0, j)
	fqSub(t1, v, &p.x)
	fqMul(t1, t1, r)
	fqSub(&p.y, t1, t0)

	fqAdd(&p.z, &a.z, &b.z)
	fqMul(&p.z, &p.z, &p.z)
	fqAdd(t0, z1z1, z2z2)
	fqSub(&p.z, &p.z, t0)
	fqMul(&p.z, &p.z, h)

	return c.Set(p)
}

// Double sets c to the sum a+a and returns c.
// See http://www.hyperelliptic.org/EFD/g1p/auto-shortw-jacobian-0.html#doubling-dbl-2009-l
func (c *curvePoint) Double(a *curvePoint) *curvePoint {
	d, e, f, g, h, i := new(fq), new(fq), new(fq), new(fq), new(fq), new(fq)
	fqMul(d, &a.x, &a.x)
	fqMul(e, &a.y, &a.y)
	fqMul(f, e, e)
	fqMul(g, &a.x, e)
	fqAdd(g, g, g)
	fqAdd(g, g, g)
	fqAdd(h, d, d)
	fqAdd(h, h, d)
	fqMul(i, h, h)

	p, t0 := new(curvePoint), new(fq)
	fqAdd(&p.x, g, g)
	fqSub(&p.x, i, &p.x)
	fqAdd(t0, f, f)
	fqAdd(t0, t0, t0)
	fqAdd(t0, t0, t0)
	fqSub(&p.y, g, &p.x)
	fqMul(&p.y, &p.y, h)
	fqSub(&p.y, &p.y, t0)
	fqMul(&p.z, &a.y, &a.z)
	fqAdd(&p.z, &p.z, &p.z)

	return c.Set(p)
}

// ScalarMult returns b*(Ax,Ay) where b is a number in big-endian form.
// See https://en.wikipedia.org/wiki/Elliptic_curve_point_multiplication#Double-and-add.
func (c *curvePoint) ScalarMult(a *curvePoint, b *big.Int) *curvePoint {
	p := new(curvePoint)
	for i := b.BitLen() - 1; i >= 0; i-- {
		p.Double(p)
		if b.Bit(i) == 1 {
			p.Add(p, a)
		}
	}

	return c.Set(p)
}

// ToAffine sets a to its affine value and returns a.
func (a *curvePoint) ToAffine() *curvePoint {
	if a.z == *new(fq).SetUint64(1) {
		return a
	}

	// TODO create new curve point
	if a.IsInfinity() {
		//  If this bit is set, the remaining bits of the group element's encoding should be set to zero.
		//pointAtInfinityMask
		return nil
	}

	zInv, zInvSqr, zInvCube := new(fq), new(fq), new(fq)
	fqInv(zInv, &a.z)
	fqMul(zInvSqr, zInv, zInv)
	fqMul(zInvCube, zInvSqr, zInv)
	fqMul(&a.x, &a.x, zInvSqr)
	fqMul(&a.y, &a.y, zInvCube)
	a.z = *new(fq).SetUint64(1)

	return a
}

// CompressedEncode converts a curve point into the uncompressed form specified in
// See https://github.com/zkcrypto/pairing/tree/master/src/bls12_381#serialization.
func (a *curvePoint) Marshal() []byte {
	a.ToAffine()

	x := new(fq).MontgomeryDecode(&a.x)
	y := new(fq).MontgomeryDecode(&a.y)

	ret := make([]byte, fqByteLen*2)
	copy(ret, x.Bytes())
	copy(ret[fqByteLen:], y.Bytes())

	// TODO review
	//if cp.IsInfinity() {
	//	ret[0] |= pointAtInfinityMask
	//}

	return ret
}

// Unmarshal decodes a curve point, serialized by Marshal.
// It is an error if the point is not on the curve.
func (cp *curvePoint) Unmarshal(data []byte) error {
	if len(data) != 2*fqByteLen {
		// TODO error
		return nil
	}

	if data[0]&compressedFormMask != 0 { // uncompressed form
		// TODO error
		return nil
	}

	if data[0]&pointAtInfinityMask == 1 {

	} else {
		cp.z = *new(fq).SetUint64(1)
	}

	var err error
	fqX, fqY := new(fq), new(fq)
	_, err = fqX.SetInt(new(big.Int).SetBytes(data[:fqByteLen]))
	if err != nil {
		return err
	}
	cp.x = *fqX
	_, err = fqY.SetInt(new(big.Int).SetBytes(data[fqByteLen:]))
	if err != nil {
		return err
	}
	cp.y = *fqY

	return nil
}

// HashToPoint sets c to the curve point that results from the given slice of bytes
// and returns c. The point is not guaranteed to be in a particular subgroup.
// See https://eprint.iacr.org/2019/403.pdf - Section 5, Construction #2.
func (c *curvePoint) HashToPoint(msg []byte) *curvePoint {
	t0 := new(curvePoint).SWUMap(hp(msg))
	t1 := new(curvePoint).SWUMap(hp(msg))
	// add note about sum before converting to E(Fq)
	t0.iso11(t0.Add(t0, t1))

	return c
}

// SWUMap maps a value of the finite field to a point in the elliptic curve.
// The point is not guaranteed to be in a particular subgroup.
// SWUMap implements an optimized version of the SWU map.
// See https://eprint.iacr.org/2019/403.pdf - Section 4.4.
func (a *curvePoint) SWUMap(t *fq) *curvePoint {
	n, t0, t1 := new(fq), new(fq), new(fq)
	fqMul(t0, t, t)
	fqMul(t1, t0, t0)
	fqNeg(t0, t0)
	fqAdd(t0, t0, t1)
	fqAdd(n, t0, new(fq).SetUint64(1))
	fqMul(n, n, fqCurveB)

	d := new(fq).Set(t0) // a = -1 for performance
	u := new(fq)
	fqMul(u, n, n)
	fqMul(u, u, n)
	fqMul(t0, d, d)
	fqMul(t1, n, t0)
	fqNeg(t1, t1)
	fqAdd(u, u, t1)
	fqMul(t0, d, t0)
	fqMul(t1, t0, fqCurveB)
	fqAdd(u, u, t1)

	v := new(fq).Set(t0)
	alpha := new(fq)
	fqMul(alpha, u, v)
	t0.Set(alpha)
	fqMul(t0, t0, v)
	fqMul(t0, t0, v)
	// t0 frob
	fqMul(alpha, alpha, t0)

	fqMul(t0, alpha, alpha)
	fqMul(t0, t0, v)
	fqSub(t0, t0, u)

	if (*t0 == fq{}) {
		fqMul(&a.x, n, d)
		fqMul(&a.y, alpha, v)
		a.z.Set(d)
	} else {
		// TODO cache Et^2, t^3
		fqMul(&a.x, n, d)
		fqMul(&a.y, alpha, v)
		a.z.Set(d)
	}

	// TODO to affine?

	return a
}

// look at isogeny of phore project simpler mul main loop or wrong?
// iso11 implements the 11-isogeny from E1´(Fp) to E1(Fp).
// See https://eprint.iacr.org/2019/403.pdf - 4.3 Isogeny maps.
func (a *curvePoint) iso11(b *curvePoint) *curvePoint {
	term := new(fq)
	mul := new(fq)
	var sum [4]fq
	for i, ki := range iso11K {
		sum[i].Set(ki[0])
		mul.SetUint64(1)
		for _, kij := range ki[1:] {
			fqMul(mul, mul, &b.x)
			fqMul(term, kij, mul)
			fqAdd(&sum[i], &sum[i], term)
		}
	}

	fqInv(&sum[1], &sum[1])
	fqMul(&a.x, &sum[0], &sum[1])
	fqInv(&sum[3], &sum[3])
	fqMul(term, &sum[2], &sum[3])
	fqMul(&a.y, &b.y, term)

	return a
}

// SWEncode implements the Shallue and van de Woestijne encoding.
// The point is not guaranteed to be in a particular subgroup.
// See https://www.di.ens.fr/~fouque/pub/latincrypt12.pdf - Algorithm 1.
func (a *curvePoint) SWEncode(b *fq) *curvePoint {
	w, inv := new(fq), new(fq)
	fqMul(w, fqSqrtNegThree, b)
	fqMul(w, w, b)
	fqMul(inv, b, b)
	fqAdd(inv, inv, fqCurveBPlusOne)
	fqInv(inv, inv)
	fqMul(w, w, inv)

	x, y := new(fq), new(fq)
	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			fqMul(x, b, w)
			fqSub(x, fqHalfSqrtNegThreeMinusOne, x)
		case 1:
			fqSub(x, &fq{0x43F5FFFFFFFCAAAE, 0x32B7FFF2ED47FFFD, 0x7E83A49A2E99D69, 0xECA8F3318332BB7A, 0xEF148D1EA0F4C069, 0x40AB3263EFF0206}, x)
		case 2:
			fqMul(x, w, w)
			fqInv(x, x)
			fqAdd(x, x, new(fq).SetUint64(1))
		}

		fqMul(y, x, x)
		fqMul(y, y, x)
		fqAdd(y, y, fqCurveB)
		if fqSqrt(y, y) {
			a.x, a.y, a.z = *x, *y, *new(fq).SetUint64(1)
			return a
		}
	}

	return a
}
