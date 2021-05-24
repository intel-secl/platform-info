// +build linux

/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package platforminfo

import (
	"os"
	"os/exec"
	"strings"

	"github.com/intel-secl/intel-secl/v3/pkg/lib/common/utils"
	"github.com/pkg/errors"
)

const (
	BindingKeyCertificatePath = "/etc/workload-agent/bindingkey.pem"
	WLAGENT                   = "wlagent"
)

// Struct used to hold the current host's platform information that can be encoded/decoded to
// json (see example below).
//
// PLEASE NOTE THAT THE PLATFORMINFO NEEDS TO BE RUN AS ROOT ON LINUX.
//
// {
//     "errorCode": 0,
//     "os_name": "RedHatEnterpriseServer",
//     "os_version": "7.6",
//     "bios_version": "SE5C620.86B.00.01.0014.070920180847",
//     "vmm_name": "",
//     "vmm_version": "",
//     "processor_info": "54 06 05 00 FF FB EB BF",
//     "host_name": "Purley32",
//     "bios_name": "Intel Corporation",
//     "hardware_uuid": "809797df-6d2d-e711-906e-0017a4403562",
//     "processor_flags": "fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc aperfmperf eagerfpu pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid dca sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch epb cat_l3 cdp_l3 intel_ppin intel_pt ssbd mba ibrs ibpb stibp tpr_shadow vnmi flexpriority ept vpid fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm cqm mpx rdt_a avx512f avx512dq rdseed adx smap clflushopt clwb avx512cd avx512bw avx512vl xsaveopt xsavec xgetbv1 cqm_llc cqm_occup_llc cqm_mbm_total cqm_mbm_local dtherm ida arat pln pts hwp hwp_act_window hwp_epp hwp_pkg_req pku ospke spec_ctrl intel_stibp flush_l1d",
//     "tpm_version": "2.0",
//     "pcr_banks": [
//         "SHA1",
//         "SHA256"
//     ],
//     "no_of_sockets": "2",
//     "tpm_enabled": "true",
//     "txt_enabled": "true",
//     "tboot_installed": "true",
//     "is_docker_env": "false",
//     "hardware_features": {
//         "TXT": {
//             "enabled": true
//         },
//         "TPM": {
//             "enabled": true,
//             "meta": {
//                 "tpm_version": "2.0",
//                 "pcr_banks": "SHA1_SHA256"
//             }
//         }
//     },
//     "installed_components": [
//         "tagent"
//     ]
// }
//
//------------------------------------------------------------------------------------------------
// BootGuard machines will contain the 'CBNT' section in 'hardware_features'...
//------------------------------------------------------------------------------------------------
//
//    "hardware_features": {
//         "CBNT": {
// 	            "enabled": true,
// 	            "meta": {
// 	                "profile": "BTGP4",
// 	                "msr": "mk ris kfm"
// 	            }
// 	        },
//         "TXT": {
//             "enabled": true
//         },
//         "TPM": {
//             "enabled": true,
//             "meta": {
//                 "tpm_version": "2.0",
//                 "pcr_banks": "SHA1_SHA256"
//             }
//         }
//    }, //
//
//------------------------------------------------------------------------------------------------
// Secure Boot systems will contain the 'secure_boot_enabled' meta section under 'UEFI' section in 'hardware_features'...
//------------------------------------------------------------------------------------------------
//
//    "hardware_features": {
//         "TXT": {
// 	            "enabled": false
// 	        },
// 	        "TPM": {
// 	            "enabled": true,
// 	            "meta": {
// 	                "tpm_version": "2.0",
// 	                "pcr_banks": "SHA1_SHA256"
// 	            }
// 	        },
// 	        "UEFI": {
//		    "enabled": "true",
//		    "meta": {
//			"secure_boot_enabled": false
//		    }
//		},
//     },
//

type HardwareFeature struct {
	Supported bool `json:"supported,string"`
	Enabled   bool `json:"enabled,string"`
}

type CBNT struct {
	HardwareFeature
	Meta struct {
		Profile string `json:"profile"`
		MSR     string `json:"msr"`
	} `json:"meta"`
}

type TPM struct {
	HardwareFeature
	Meta struct {
		TPMVersion string `json:"tpm_version"`
		PCRBanks   string `json:"pcr_banks"`
	} `json:"meta"`
}

type UEFI struct {
	HardwareFeature
	Meta struct {
		SecureBootEnabled bool `json:"secure_boot_enabled"`
	} `json:"meta"`
}

type PlatformInfo struct {
	ErrorCode           int    `json:"errorCode"`
	OSName              string `json:"os_name"`
	OSVersion           string `json:"os_version"`
	BiosVersion         string `json:"bios_version"`
	VMMName             string `json:"vmm_name"`
	VMMVersion          string `json:"vmm_version"`
	ProcessorInfo       string `json:"processor_info"`
	HostName            string `json:"host_name"`
	BiosName            string `json:"bios_name"`
	HardwareUUID        string `json:"hardware_uuid"`
	ProcessorFlags      string `json:"process_flags"`
	TPMVersion          string `json:"tpm_version"`
	NumberOfSockets     int    `json:"no_of_sockets,string"`
	TPMEnabled          bool   `json:"tpm_enabled,string"`
	TXTEnabled          bool   `json:"txt_enabled,string"`
	TbootInstalled      bool   `json:"tboot_installed,string"`
	IsDockerEnvironment bool   `json:"is_docker_env,string"`
	HardwareFeatures    struct {
		TXT  HardwareFeature `json:"TXT"`
		TPM  TPM             `json:"TPM"`
		CBNT CBNT            `json:"CBNT"`
		UEFI UEFI            `json:"UEFI"`
		PFR  HardwareFeature `json:"PFR"`
		BMC  HardwareFeature `json:"BMC"`
	} `json:"hardware_features"`
	InstalledComponents []string `json:"installed_components"`
}

func GetPlatformInfo() (*PlatformInfo, error) {
	var err, rerr error
	platformInfo := PlatformInfo{}

	platformInfo.OSName, rerr = OSName()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting OS Name")
	}

	platformInfo.OSVersion, rerr = OSVersion()
	if rerr != nil {
		err = errors.Wrap(rerr, "Error getting OS Version")
	}

	platformInfo.BiosVersion, rerr = BiosVersion()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting BIOS Version")
	}

	platformInfo.VMMName, rerr = VMMName()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting VMM Name")
	}

	platformInfo.VMMVersion, rerr = VMMVersion()
	if rerr != nil {
		return nil, errors.Wrap(err, "Error getting VMM Version")
	}

	platformInfo.ProcessorInfo, rerr = ProcessorID()
	if rerr != nil {
		return nil, errors.Wrap(err, "Error getting Processor ID")
	}

	platformInfo.HostName, rerr = HostName()
	if rerr != nil {
		return nil, errors.Wrap(err, "Error getting Host Name")
	}

	platformInfo.BiosName, rerr = BiosName()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting BIOS Name")
	}

	platformInfo.HardwareUUID, rerr = HardwareUUID()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting Hardware UUID")
	}

	processorFlags, rerr := ProcessorFlags()
	if rerr != nil {
		processorFlags = []string{}
		err = errors.Wrap(err, "Error getting Processor Flags")
	}
	platformInfo.ProcessorFlags = strings.Join(processorFlags, " ")

	platformInfo.TPMVersion, rerr = TPMVersion()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting TPM Version")
	}

	platformInfo.NumberOfSockets, rerr = NoOfSockets()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting number of processor sockets")
	}

	platformInfo.TPMEnabled, rerr = TPMEnabled()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting TPM Status")
	}

	platformInfo.TXTEnabled, rerr = TXTEnabled()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting TXT Status")
	}

	platformInfo.TbootInstalled, rerr = TbootInstalled()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting TBOOT status")
	}

	platformInfo.IsDockerEnvironment = utils.IsContainerEnv()
	platformInfo.HardwareFeatures.TXT.Enabled = platformInfo.TXTEnabled
	platformInfo.HardwareFeatures.TXT.Supported = platformInfo.TXTEnabled
	platformInfo.HardwareFeatures.TPM.Enabled = platformInfo.TPMEnabled
	platformInfo.HardwareFeatures.TPM.Supported = platformInfo.TPMEnabled
	platformInfo.HardwareFeatures.TPM.Meta.TPMVersion = platformInfo.TPMVersion
	platformInfo.InstalledComponents = []string{"tagent"}

	platformInfo.HardwareFeatures.CBNT, rerr = GetCBNTHardwareFeature()
	if rerr != nil {
		err = errors.Wrap(err, "Error getting CBNT information")
	}

	platformInfo.HardwareFeatures.UEFI, err = GetUEFIHardwareFeature()
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting UEFI information")
	}

	if WLAIsInstalled() {
		platformInfo.InstalledComponents = append(platformInfo.InstalledComponents, WLAGENT)
	}

	if err != nil {
		platformInfo.ErrorCode = 1
	}

	return &platformInfo, err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Run 'which wlagent'.  If the command returns '0' (no error) then workload-agent is installed.
func WLAIsInstalled() bool {
	cmd := exec.Command("which", WLAGENT)

	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
