module intel/isecl/lib/platform-info/v4

require (
	github.com/pkg/errors v0.9.1
	intel/isecl/lib/common/v4 v4.2.0-Beta
	github.com/intel-secl/intel-secl/v4 v4.2.0-Beta
)

replace (
    intel/isecl/lib/common/v4 => github.com/intel-secl/common/v4 v4.2.0-Beta
)
