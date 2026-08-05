package types

import "time"

// ServiceZone is the wall clock an MPP building runs on.
//
// Storage and the API stay UTC (see docs/03-architecture); this only decides
// where a calendar DAY begins and ends — which quota, queue reset and QR
// validity all depend on. Computing those in UTC makes a "day" end at 07:00
// local the next morning.
//
// A fixed offset rather than time.LoadLocation: Indonesia has no DST, and the
// scratch container images ship without tzdata.
//
// ponytail: one zone for the whole deployment until mpp.instansi carries its
// own timezone column (a WITA/WIT building needs its own offset).
var ServiceZone = time.FixedZone("WIB", 7*60*60)
