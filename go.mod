module intel/isecl/lib/platform-info/v4

require (
	github.com/pkg/errors v0.9.1
	intel/isecl/lib/common/v4 v4.0.0
	github.com/intel-secl/intel-secl/v4 v4.0.0
)

replace (
    intel/isecl/lib/common/v4 => gitlab.devtools.intel.com/sst/isecl/lib/common.git/v4 v4.0/develop
    github.com/intel-secl/intel-secl/v4 => gitlab.devtools.intel.com/sst/isecl/intel-secl.git/v4 v4.0/develop
)
