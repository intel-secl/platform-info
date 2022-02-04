module intel/isecl/lib/platform-info/v4

require (
	github.com/pkg/errors v0.9.1
	intel/isecl/lib/common/v4 v4.0.2
	github.com/intel-secl/intel-secl/v4 v4.0.2
)

replace (
    intel/isecl/lib/common/v4 => github.com/intel-innersource/libraries.security.isecl.common/v4 v4.0.2/develop
    github.com/intel-secl/intel-secl/v4 => github.com/intel-innersource/applications.security.isecl.intel-secl/v4 v4.0.2/develop
)
