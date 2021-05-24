// +build linux

/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */

package platforminfo

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

var (
	ErrNoVMMError                   = errors.New("VMM is not installed")
	ErrMissingVirshVersionFileError = errors.New("Empty virsh version file, assuming no VMM installed")
	ErrNoTPMError                   = errors.New("No TPM chip present on host")
	ErrTPMStatusError               = errors.New("Error in getting the TPM status")
)

const (
	tpm12File            = "/sys/class/tpm/tpm0/device/enabled"
	tpm20File0           = "/sys/class/tpm/tpm0/device/description"
	tpm20File1           = "/sys/class/tpm/tpm0/device/firmware_node/description"
	classMiscDevCapsFile = "/sys/class/misc/tpm0/device/caps"
	classTpmDevCapsFile  = "/sys/class/tpm/tpm0/device/caps"
	tpmDevFile           = "/dev/tpm0"
	tpmVersion12         = "1.2"
	tpmVersion20         = "2.0"
	txtMLELaunchMsg      = "TXT measured launch:"
	CBNT_MSR_OFFSET      = 0x13A
	TXT_MSR_OFFSET       = 0x3A
	MSR_DEVICE           = "/dev/cpu/0/msr"
	CBNT_PROCESSOR_FLAGS = "mk ris kfm"
	CBNT_PROFILE_0       = "BTGP0"
	CBNT_PROFILE_3       = "BTGP3"
	CBNT_PROFILE_4       = "BTGP4"
	CBNT_PROFILE_5       = "BTGP5"
	CBNT_PROFILE_0_FLAGS = 0x0
	CBNT_PROFILE_3_FLAGS = 0x6d
	CBNT_PROFILE_4_FLAGS = 0x51
	CBNT_PROFILE_5_FLAGS = 0x7d
	TXT_STATUS_ENABLED   = 0x3
	SECURE_BOOT_ENABLED  = "Secure Boot: enabled"
)

var (
	biosNameCmd      = []string{"dmidecode", "-s", "bios-vendor"}
	biosVersionCmd   = []string{"dmidecode", "-s", "bios-version"}
	hardwareUUIDCmd  = []string{"dmidecode", "-s", "system-uuid"}
	osInfoCmd        = []string{"lsb_release", "-a"}
	processorInfoCmd = []string{"dmidecode", "--type", "processor"}
	dockerVersionCmd = []string{"dockerd", "-v"}
	virshVersionCmd  = []string{"virsh", "-v"}
	hostNameCmd      = []string{"hostname"}
	noSocketsCmd     = []string{"lscpu"}
	txtEnabledCmd    = []string{"txt-stat"}
	suefiBootCtlCmd  = []string{"bootctl", "status"}
)

// BiosName retrieves the host BIOS name.
// An example of the host bios name is Intel Corporation
func BiosName() (string, error) {
	//Output of command:
	//Intel Corporation
	result, err := readAndParseFromCommandLine(biosNameCmd)
	if err != nil {
		return "", err
	}
	biosName := ""
	// Split the output separated by new line into list
	for i := range result {
		biosName = strings.TrimSpace(result[i])
		break
	}
	return biosName, err
}

// BiosVersion retrieves the host BIOS version.
// An example of the host bios version is S5500.86B.01.00.0060.090920111354
func BiosVersion() (string, error) {
	//Output of command:
	//S5500.86B.01.00.0060.090920111354
	result, err := readAndParseFromCommandLine(biosVersionCmd)
	if err != nil {
		return "-1", err
	}

	biosVersion := ""
	for i := range result {
		biosVersion = strings.TrimSpace(result[i])
		break
	}

	return biosVersion, err
}

// HardwareUUID retrieves the host hardware UUID.
// An example of the host hardware UUID is 4219B2F5-C25F-6AF2-573C-35B0DF557236
func HardwareUUID() (string, error) {
	//Output of command:
	//4219B2F5-C25F-6AF2-573C-35B0DF557236
	result, err := readAndParseFromCommandLine(hardwareUUIDCmd)
	if err != nil {
		return "-1", err
	}

	hardwareUUID := ""
	for i := range result {
		hardwareUUID = strings.TrimSpace(result[i])
		break
	}
	return hardwareUUID, err
}

// OSName retrieves the host OS name.
// An example of the host OS name is Ubuntu
func OSName() (string, error) {

	osName, _, err := osNameAndVersion()
	return osName, err
}

// OSVersion retrieves the host OS version.
// An example of the host OS version is 11.10
func OSVersion() (string, error) {
	_, osVersion, err := osNameAndVersion()
	return osVersion, err
}

// osNameandVersion retrieves the host OS name and version
func osNameAndVersion() (string, string, error) {
	/*
	  Sample response of 'lsb_release -a'

	  No LSB modules are available.
	  Distributor ID: Ubuntu
	  Description:    Ubuntu 11.10
	  Release:        11.10
	  Codename:       oneiric
	*/
	result, err := readAndParseFromCommandLine(osInfoCmd)
	if err != nil {
		return "", "-1", err
	}

	osName, osVersion := "", ""
	for i := range result {
		if strings.HasPrefix(result[i], "Distributor ID:") {
			osName = strings.TrimSpace(strings.TrimPrefix(result[i], "Distributor ID:"))
		} else if strings.HasPrefix(result[i], "Release:") {
			osVersion = strings.TrimSpace(strings.TrimPrefix(result[i], "Release:"))
		}
	}
	return osName, osVersion, err
}

// NoOfSockets retrieves the number of CPU sockets on the host platform.
func NoOfSockets() (int, error) {
	result, err := readAndParseFromCommandLine(noSocketsCmd)
	if err != nil {
		return -1, err
	}

	noOfSockets := 0
	for i := range result {
		if strings.HasPrefix(result[i], "Socket(s):") {
			noOfSockets, _ = strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(result[i], "Socket(s):")))
			break
		}
	}
	return noOfSockets, err
}

// ProcessorFlags retrieves the processor flags.
// Some examples of host processor flags are FPU, VME
func ProcessorFlags() ([]string, error) {
	_, flags, err := processorInfo()
	return flags, err
}

// ProcessorID retrieves the processor ID.
// An example of the processor ID is ABCDFGHI012345JK
func ProcessorID() (string, error) {
	id, _, err := processorInfo()
	return id, err
}

func processorInfo() (string, []string, error) {
	/*
		  Sample output of 'dmidecode --type processor'

		  Processor Information
				  Socket Designation: CPU1
				  Type: Central Processor
				  Family: Xeon
				  Manufacturer: Intel(R) Corporation
				  ID: C2 06 02 00 FF FB EB BF
				  Signature: Type 0, Family 6, Model 44, Stepping 2
				  Flags:
						  FPU (Floating-point unit on-chip)
						  VME (Virtual mode extension)
						  DE (Debugging extension)
						  PSE (Page size extension)
						  TSC (Time stamp counter)
						  ...
				  Version: Intel(R) Xeon(R) Gold 6140 CPU @ 2.30GHz
				  Voltage: 1.6 V
				  External Clock: 100 MHz
				  Max Speed: 4000 MHz
				  Current Speed: 2300 MHz
				  Status: Populated, Enabled
				  Upgrade: Socket LGA3647-1
				  L1 Cache Handle: 0x000E
				  L2 Cache Handle: 0x000F
				  L3 Cache Handle: 0x0010
				  Serial Number: Not Specified
				  Asset Tag: UNKNOWN
				  Part Number: Not Specified
				  Core Count: 18
				  Core Enabled: 18
				  Thread Count: 36
				  Characteristics:
						  64-bit capable
						  Multi-Core
						  Hardware Thread
						  Execute Protection
						  Enhanced Virtualization
						  Power/Performance Control
	*/
	result, err := readAndParseFromCommandLine(processorInfoCmd)
	if err != nil {
		return "-1", make([]string, 0), err
	}

	id := ""
	flags := make([]string, 0)
	flagRegExp, _ := regexp.Compile(`^\t*[A-Z0-9\-]+ \([A-Za-z0-9\- ]+\)$`)
	for i := range result {
		line := result[i]
		if strings.HasPrefix(line, "ID:") {
			id = strings.Join(strings.Fields(line)[1:], " ")
		} else if flagRegExp.MatchString(line) {
			flags = append(flags, strings.Fields(line)[0])
		} else if strings.HasPrefix(line, "Version:") {
			break
		}
	}
	return id, flags, err
}

// VMMName retrives the name of the hypervisor running on the host.
// An example of the hypervisor name is Docker.
func VMMName() (string, error) {
	vmmName, _, err := vmmInfo()
	return vmmName, err
}

// VMMVersion retrieves the version of the hypervisor running on the host.
// An example of the hypervisor version is 1.13.1
func VMMVersion() (string, error) {
	_, vmmVersion, err := vmmInfo()
	return vmmVersion, err
}

// vmmNameAndVersion retrieves the name and version of the hypervisor running on the host.
func vmmInfo() (string, string, error) {
	// Check if docker is installed
	/*
	  Sample response of 'dockerd -v'

	  Docker version 19.03.5, build 633a0ea
	*/
	result, err := readAndParseFromCommandLine(dockerVersionCmd)

	// if no error occurs, assume docker is installed, return name and version...
	if err == nil && len(result) != 0 {
		response := result[0]
		parts := strings.Split(response, " ")
		return parts[0], strings.TrimSuffix(parts[2], ","), nil
	}

	// Check if virsh is installed...
	result, err = readAndParseFromCommandLine(virshVersionCmd)
	if err == nil && len(result) != 0 {
		response := result[0]
		parts := strings.Split(response, " ")
		return "Virsh", parts[0], nil
	}

	// if neither is installed, just return empty strings to satisfy platforminfo json
	return "", "", nil
}

// HostName retrieves the network hostname.
// An example of the hostname is Redhat
func HostName() (string, error) {
	result, err := readAndParseFromCommandLine(hostNameCmd)
	if err != nil {
		return "", err
	}

	hostName := strings.TrimSpace(result[0])
	return hostName, err
}

// TPMVersion retrieves the version of the installed Trusted Platform Module (TPM).
// An example of the host TPM version is 1.2
func TPMVersion() (string, error) {
	_, err := os.Stat(tpmDevFile)
	if os.IsNotExist(err) {
		return "0", nil
	}
	_, err = os.Stat(classMiscDevCapsFile)
	if !os.IsNotExist(err) {
		return tpmVersion12, nil
	}
	_, err = os.Stat(classTpmDevCapsFile)
	if !os.IsNotExist(err) {
		return tpmVersion12, nil
	}
	return tpmVersion20, nil
}

// TPMEnabled indicates whether the Trusted Platform Module is enabled or not.
func TPMEnabled() (bool, error) {
	version, err := TPMVersion()
	if version == "0" {
		return false, nil
	}

	tpmEnabled := false
	if version == tpmVersion12 {
		_, err := readAndParseFromCommandLine([]string{"cat", tpm12File})
		tpmEnabled = err == nil
	} else if version == tpmVersion20 {
		_, err = os.Stat(tpm20File0)
		tpmEnabled = !os.IsNotExist(err)
		if !tpmEnabled {
			_, err = os.Stat(tpm20File1)
			if !os.IsNotExist(err) {
				tpmEnabled = true
			} else {
				err = ErrTPMStatusError
			}
		}
	}
	return tpmEnabled, err
}

// Run txt-stat to see if tboot is installed
func TbootInstalled() (bool, error) {
	result, err := readAndParseFromCommandLine(txtEnabledCmd)
	if err != nil {
		return false, nil
	}

	txtStatus := false
	for i := range result {
		trimmed := strings.TrimSpace(result[i])
		if strings.HasPrefix(trimmed, txtMLELaunchMsg) {
			txtStatus = strings.Contains(trimmed, "TRUE")
			break
		}
	}
	return txtStatus, err
}

// Utility function that reads an unsigned long long from /dev/cpu/0/msr at offset
// 'offset'
func ReadMSR(offset int64) (uint64, error) {

	msr, err := os.Open(MSR_DEVICE)
	if err != nil {
		return 0, errors.Wrapf(err, "platforminfo:ReadMSR(): Error opening msr")
	}

	defer msr.Close()

	_, err = msr.Seek(offset, 0)
	if err != nil {
		return 0, errors.Wrapf(err, "platforminfo:ReadMSR(): Could not seek to location %x", offset)
	}

	results := make([]byte, 8)
	length, err := msr.Read(results)
	if length < 8 {
		return 0, errors.New("platforminfo:ReadMSR(): Reading the msr returned the incorrect length")
	} else if err != nil {
		return 0, errors.Wrapf(err, "platforminfo:ReadMSR(): There was an error reading msr at offset %x", offset)
	}

	return binary.LittleEndian.Uint64(results), nil
}

// Utility function that returns the integer value between 'hibit' and 'lowbit'
func GetBits(value uint64, hibit uint, lowbit uint) (uint64, error) {
	bits := hibit - lowbit + 1
	if bits > 64 {
		return 0, errors.New("platforminfo:GetBits(): Invalid hi/low bit parameters")
	}

	value >>= lowbit
	value &= (uint64(1) << bits) - 1
	return value, nil
}

// Read the MSR at 0x3a to determine if txt is enabled (i.e. bits 1:0 return '3')
// similar to running "rdmsr -f 1:0 0x3a"
func TXTEnabled() (bool, error) {

	txtBits, err := ReadMSR(TXT_MSR_OFFSET)
	if err != nil {
		return false, err
	}

	txtStatus, err := GetBits(txtBits, 1, 0)
	if err != nil {
		return false, err
	}

	if txtStatus == TXT_STATUS_ENABLED {
		return true, nil
	}

	return false, nil
}

// Provide the CBNT bits, read bits 7:0 to determine the boot guard
// profile.  This function assumes that cbnt is enabled when being called.  The value
// 0x51 indicates BTGP4, 0x7d is BTGP5.
//
// Similar to "rdmsr -f 7:0 0x13A"
func GetCBNTProfile(cbntBits uint64) (string, error) {

	cbntProfileFlags, err := GetBits(cbntBits, 7, 0)
	if err != nil {
		return "", err
	}

	switch cbntProfileFlags {
	case CBNT_PROFILE_0_FLAGS:
		return CBNT_PROFILE_0, nil
	case CBNT_PROFILE_3_FLAGS:
		return CBNT_PROFILE_3, nil
	case CBNT_PROFILE_4_FLAGS:
		return CBNT_PROFILE_4, nil
	case CBNT_PROFILE_5_FLAGS:
		return CBNT_PROFILE_5, nil
	}

	return "", fmt.Errorf("platforminfo:GetCBNTProfile(): Could not determine CBNT profile with flags %x", cbntProfileFlags)
}

// If CBNT is supported and enabled, create a valid "CBNT" structure.  Otherwise
// return a empty structure (i.e. CBNT is not enabled).
func GetCBNTHardwareFeature() (CBNT, error) {

	var cbnt CBNT

	cbntBits, err := ReadMSR(CBNT_MSR_OFFSET)
	if err != nil {
		return cbnt, err
	}

	// Check if CBNT is supported
	// similar to 'rdmsr -f 32:32 0x13A
	cbntSupported, err := GetBits(cbntBits, 32, 32)
	if err != nil {
		return cbnt, err
	}

	// CBNT is supported, create a CBNT structure and populate it.
	if cbntSupported == 1 {
		cbnt.HardwareFeature.Supported = true

		// Verify if cbnt is enabled
		// MSR[7:4] denotes the value of (Verify/Measure/FACB)
		bitsVMF, err := GetBits(cbntBits, 7, 4)
		if err != nil {
			return cbnt, err
		}

		// MSR[0] denotes the value of (BTG enabled and passed startup ACM)
		bitBTGEnabled, err := GetBits(cbntBits, 0, 0)
		if err != nil {
			return cbnt, err
		}

		if bitsVMF != 0 || bitBTGEnabled != 0 {
			cbnt.Meta.Profile, err = GetCBNTProfile(cbntBits)
			if err != nil {
				return cbnt, err
			}
			cbnt.Enabled = cbnt.Meta.Profile != CBNT_PROFILE_0
			cbnt.Meta.MSR = CBNT_PROCESSOR_FLAGS
		}
	}

	return cbnt, nil
}

// If the current boot is in UEFI mode and Secure Boot is enabled, create a valid "UEFI" structure with "SecureBootEnabled" meta section.
// Otherwise, return empty structure (i.e. Current boot is not UEFI and secure boot is not enabled).
func GetUEFIHardwareFeature() (UEFI, error) {

	var uefi UEFI
	//Checking the presence of efi directory to make sure the current boot is in UEFI mode.
	//In case of legacy mode this directory wont be available.
	if _, err := os.Stat("/sys/firmware/efi"); os.IsNotExist(err) {
		return uefi, nil
	} else {
		uefi.HardwareFeature.Enabled = true
		uefi.HardwareFeature.Supported = true
	}

	// call 'bootctl status'. Ignore errors because 'bootctl' can return '1'
	// even when valid output is present onstdout.  Just see if stdout contains
	// "Secure Boot: enabled" to determine if SUEFI is enabled.
	cmd := exec.Command(suefiBootCtlCmd[0], suefiBootCtlCmd[1])
	stdout, _ := cmd.Output()
	if strings.Contains(string(stdout), SECURE_BOOT_ENABLED) {
		uefi.Meta.SecureBootEnabled = true
	}

	return uefi, nil
}
