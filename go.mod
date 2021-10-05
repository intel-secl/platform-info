module intel/isecl/lib/platform-info/v3

require (
	github.com/pkg/errors v0.9.1
	intel/isecl/lib/common/v3 v3.6.1
	github.com/intel-secl/intel-secl/v3 v3.6.1
)

replace (
    intel/isecl/lib/common/v3 => github.com/intel-secl/common/v3 v3.6.1
)
