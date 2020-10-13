package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/wolfeidau/cf-security-transform/internal/cf"
	lmw "github.com/wolfeidau/lambda-go-extras/middleware"
	"github.com/wolfeidau/lambda-go-extras/middleware/raw"
	zlog "github.com/wolfeidau/lambda-go-extras/middleware/zerolog"
)

// assigned during build time with -ldflags
var (
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	flds := lmw.FieldMap{"commit": commit, "buildDate": buildDate}

	tr := new(cf.Transform)

	ch := lmw.New(
		raw.New(raw.Fields(flds)),   // raw event logger primarily used during development
		zlog.New(zlog.Fields(flds)), // inject zerolog into the context
	).Then(tr)

	lambda.StartHandler(ch)
}
