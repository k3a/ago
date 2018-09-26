package currency

import "strings"

// Based on https://github.com/rmg/eagle/blob/master/Godeps/_workspace/src/github.com/rmg/iso4217/constants.go
// License of this file only is MIT
// Copyright 2018 K3A.me

// Currency identifier
type Currency uint16

var names = map[uint16]string{
	0:   "---",
	8:   "ALL",
	12:  "DZD",
	32:  "ARS",
	36:  "AUD",
	44:  "BSD",
	48:  "BHD",
	50:  "BDT",
	51:  "AMD",
	52:  "BBD",
	60:  "BMD",
	64:  "BTN",
	68:  "BOB",
	72:  "BWP",
	84:  "BZD",
	90:  "SBD",
	96:  "BND",
	104: "MMK",
	108: "BIF",
	116: "KHR",
	124: "CAD",
	132: "CVE",
	136: "KYD",
	144: "LKR",
	152: "CLP",
	156: "CNY",
	170: "COP",
	174: "KMF",
	188: "CRC",
	191: "HRK",
	192: "CUP",
	203: "CZK",
	208: "DKK",
	214: "DOP",
	222: "SVC",
	230: "ETB",
	232: "ERN",
	238: "FKP",
	242: "FJD",
	262: "DJF",
	270: "GMD",
	292: "GIP",
	320: "GTQ",
	324: "GNF",
	328: "GYD",
	332: "HTG",
	340: "HNL",
	344: "HKD",
	348: "HUF",
	352: "ISK",
	356: "INR",
	360: "IDR",
	364: "IRR",
	368: "IQD",
	376: "ILS",
	388: "JMD",
	392: "JPY",
	398: "KZT",
	400: "JOD",
	404: "KES",
	408: "KPW",
	410: "KRW",
	414: "KWD",
	417: "KGS",
	418: "LAK",
	422: "LBP",
	426: "LSL",
	430: "LRD",
	434: "LYD",
	440: "LTL",
	446: "MOP",
	454: "MWK",
	458: "MYR",
	462: "MVR",
	478: "MRO",
	480: "MUR",
	484: "MXN",
	496: "MNT",
	498: "MDL",
	504: "MAD",
	512: "OMR",
	516: "NAD",
	524: "NPR",
	532: "ANG",
	533: "AWG",
	548: "VUV",
	554: "NZD",
	558: "NIO",
	566: "NGN",
	578: "NOK",
	586: "PKR",
	590: "PAB",
	598: "PGK",
	600: "PYG",
	604: "PEN",
	608: "PHP",
	634: "QAR",
	643: "RUB",
	646: "RWF",
	654: "SHP",
	678: "STD",
	682: "SAR",
	690: "SCR",
	694: "SLL",
	702: "SGD",
	704: "VND",
	706: "SOS",
	710: "ZAR",
	728: "SSP",
	748: "SZL",
	752: "SEK",
	756: "CHF",
	760: "SYP",
	764: "THB",
	776: "TOP",
	780: "TTD",
	784: "AED",
	788: "TND",
	800: "UGX",
	807: "MKD",
	818: "EGP",
	826: "GBP",
	834: "TZS",
	840: "USD",
	858: "UYU",
	860: "UZS",
	882: "WST",
	886: "YER",
	901: "TWD",
	931: "CUC",
	932: "ZWL",
	934: "TMT",
	936: "GHS",
	937: "VEF",
	938: "SDG",
	940: "UYI",
	941: "RSD",
	943: "MZN",
	944: "AZN",
	946: "RON",
	947: "CHE",
	948: "CHW",
	949: "TRY",
	950: "XAF",
	951: "XCD",
	952: "XOF",
	953: "XPF",
	955: "XBA",
	956: "XBB",
	957: "XBC",
	958: "XBD",
	959: "XAU",
	960: "XDR",
	961: "XAG",
	962: "XPT",
	963: "XTS",
	964: "XPD",
	965: "XUA",
	967: "ZMW",
	968: "SRD",
	969: "MGA",
	970: "COU",
	971: "AFN",
	972: "TJS",
	973: "AOA",
	974: "BYR",
	975: "BGN",
	976: "CDF",
	977: "BAM",
	978: "EUR",
	979: "MXV",
	980: "UAH",
	981: "GEL",
	984: "BOV",
	985: "PLN",
	986: "BRL",
	990: "CLF",
	994: "XSU",
	997: "USN",
	999: "XXX",
}

func (c Currency) String() string {
	nm, _ := names[uint16(c)]
	return nm
}

// Currency constant
const (
	Unknown = Currency(0)
	ALL     = Currency(8)
	DZD     = Currency(12)
	ARS     = Currency(32)
	AUD     = Currency(36)
	BSD     = Currency(44)
	BHD     = Currency(48)
	BDT     = Currency(50)
	AMD     = Currency(51)
	BBD     = Currency(52)
	BMD     = Currency(60)
	BTN     = Currency(64)
	BOB     = Currency(68)
	BWP     = Currency(72)
	BZD     = Currency(84)
	SBD     = Currency(90)
	BND     = Currency(96)
	MMK     = Currency(104)
	BIF     = Currency(108)
	KHR     = Currency(116)
	CAD     = Currency(124)
	CVE     = Currency(132)
	KYD     = Currency(136)
	LKR     = Currency(144)
	CLP     = Currency(152)
	CNY     = Currency(156)
	COP     = Currency(170)
	KMF     = Currency(174)
	CRC     = Currency(188)
	HRK     = Currency(191)
	CUP     = Currency(192)
	CZK     = Currency(203)
	DKK     = Currency(208)
	DOP     = Currency(214)
	SVC     = Currency(222)
	ETB     = Currency(230)
	ERN     = Currency(232)
	FKP     = Currency(238)
	FJD     = Currency(242)
	DJF     = Currency(262)
	GMD     = Currency(270)
	GIP     = Currency(292)
	GTQ     = Currency(320)
	GNF     = Currency(324)
	GYD     = Currency(328)
	HTG     = Currency(332)
	HNL     = Currency(340)
	HKD     = Currency(344)
	HUF     = Currency(348)
	ISK     = Currency(352)
	INR     = Currency(356)
	IDR     = Currency(360)
	IRR     = Currency(364)
	IQD     = Currency(368)
	ILS     = Currency(376)
	JMD     = Currency(388)
	JPY     = Currency(392)
	KZT     = Currency(398)
	JOD     = Currency(400)
	KES     = Currency(404)
	KPW     = Currency(408)
	KRW     = Currency(410)
	KWD     = Currency(414)
	KGS     = Currency(417)
	LAK     = Currency(418)
	LBP     = Currency(422)
	LSL     = Currency(426)
	LRD     = Currency(430)
	LYD     = Currency(434)
	LTL     = Currency(440)
	MOP     = Currency(446)
	MWK     = Currency(454)
	MYR     = Currency(458)
	MVR     = Currency(462)
	MRO     = Currency(478)
	MUR     = Currency(480)
	MXN     = Currency(484)
	MNT     = Currency(496)
	MDL     = Currency(498)
	MAD     = Currency(504)
	OMR     = Currency(512)
	NAD     = Currency(516)
	NPR     = Currency(524)
	ANG     = Currency(532)
	AWG     = Currency(533)
	VUV     = Currency(548)
	NZD     = Currency(554)
	NIO     = Currency(558)
	NGN     = Currency(566)
	NOK     = Currency(578)
	PKR     = Currency(586)
	PAB     = Currency(590)
	PGK     = Currency(598)
	PYG     = Currency(600)
	PEN     = Currency(604)
	PHP     = Currency(608)
	QAR     = Currency(634)
	RUB     = Currency(643)
	RWF     = Currency(646)
	SHP     = Currency(654)
	STD     = Currency(678)
	SAR     = Currency(682)
	SCR     = Currency(690)
	SLL     = Currency(694)
	SGD     = Currency(702)
	VND     = Currency(704)
	SOS     = Currency(706)
	ZAR     = Currency(710)
	SSP     = Currency(728)
	SZL     = Currency(748)
	SEK     = Currency(752)
	CHF     = Currency(756)
	SYP     = Currency(760)
	THB     = Currency(764)
	TOP     = Currency(776)
	TTD     = Currency(780)
	AED     = Currency(784)
	TND     = Currency(788)
	UGX     = Currency(800)
	MKD     = Currency(807)
	EGP     = Currency(818)
	GBP     = Currency(826)
	TZS     = Currency(834)
	USD     = Currency(840)
	UYU     = Currency(858)
	UZS     = Currency(860)
	WST     = Currency(882)
	YER     = Currency(886)
	TWD     = Currency(901)
	CUC     = Currency(931)
	ZWL     = Currency(932)
	TMT     = Currency(934)
	GHS     = Currency(936)
	VEF     = Currency(937)
	SDG     = Currency(938)
	UYI     = Currency(940)
	RSD     = Currency(941)
	MZN     = Currency(943)
	AZN     = Currency(944)
	RON     = Currency(946)
	CHE     = Currency(947)
	CHW     = Currency(948)
	TRY     = Currency(949)
	XAF     = Currency(950)
	XCD     = Currency(951)
	XOF     = Currency(952)
	XPF     = Currency(953)
	XBA     = Currency(955)
	XBB     = Currency(956)
	XBC     = Currency(957)
	XBD     = Currency(958)
	XAU     = Currency(959)
	XDR     = Currency(960)
	XAG     = Currency(961)
	XPT     = Currency(962)
	XTS     = Currency(963)
	XPD     = Currency(964)
	XUA     = Currency(965)
	ZMW     = Currency(967)
	SRD     = Currency(968)
	MGA     = Currency(969)
	COU     = Currency(970)
	AFN     = Currency(971)
	TJS     = Currency(972)
	AOA     = Currency(973)
	BYR     = Currency(974)
	BGN     = Currency(975)
	CDF     = Currency(976)
	BAM     = Currency(977)
	EUR     = Currency(978)
	MXV     = Currency(979)
	UAH     = Currency(980)
	GEL     = Currency(981)
	BOV     = Currency(984)
	PLN     = Currency(985)
	BRL     = Currency(986)
	CLF     = Currency(990)
	XSU     = Currency(994)
	USN     = Currency(997)
	XXX     = Currency(999)
)

// FromString returns currency representing
// a string identifier or Unknown for invalid code
func FromString(code string) Currency {
	code = strings.ToUpper(code)

	for id, nm := range names {
		if nm == code {
			return Currency(id)
		}
	}

	return Unknown
}
