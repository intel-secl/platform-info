/*
 * Copyright (C) 2019 Intel Corporation
 * SPDX-License-Identifier: BSD-3-Clause
 */
package main

import (
	"encoding/json"
	"fmt"
	"intel/isecl/lib/common/v4/validation"
	p "intel/isecl/lib/platform-info/v4/platforminfo"
	"os"
)

func main() {

	// if no arguments were provided, just dump platform-info json to stdout
	if len(os.Args[1:]) < 1 {
		platformInfo, err := p.GetPlatformInfo()
		if err != nil {
			fmt.Printf("main: Error while fetching platform info: %s", err)
			os.Exit(1)
		}

		// serialize to json
		b, err := json.MarshalIndent(platformInfo, "", "\t")
		if err != nil {
			fmt.Printf("main: Error while serializing platform info: %s", err)
			os.Exit(1)
		}

		fmt.Printf(string(b))
		os.Exit(0)
	}

	inputValArr := os.Args[0:]
	if valErr := validation.ValidateStrings(inputValArr); valErr != nil {
		fmt.Println("Invalid string format")
		os.Exit(1)
	}

	var methodName = os.Args[1]
	switch methodName {
	case "BiosName":
		fmt.Println("Bios Name method called")
		biosName, err := p.BiosName()
		if err != nil {
			fmt.Println("Error during retrieval of Bios Name: ", err)
			os.Exit(1)
		}
		fmt.Println("Bios Name:\n", biosName)

	case "BiosVersion":
		fmt.Println("Bios Version method called")
		biosVersion, err := p.BiosVersion()
		if err != nil {
			fmt.Println("Error during retrieval of Bios Version: ", err)
			os.Exit(1)
		}
		fmt.Println("Bios Version:\n", biosVersion)

	case "HardwareUUID":
		fmt.Println("Hardware UUID method called")
		hardwareUUID, err := p.HardwareUUID()
		if err != nil {
			fmt.Println("Error during retrieval of Hardware UUID: ", err)
			os.Exit(1)
		}
		fmt.Println("Hardware UUID:\n", hardwareUUID)

	case "OSName":
		fmt.Println("OS Name method called")
		osName, err := p.OSName()
		if err != nil {
			fmt.Println("Error during retrieval of OS Name: ", err)
			os.Exit(1)
		}
		fmt.Println("OS Name:\n", osName)

	case "OSVersion":
		fmt.Println("OS Version method called")
		osVersion, err := p.OSVersion()
		if err != nil {
			fmt.Println("Error during retrieval of OS Version: ", err)
			os.Exit(1)
		}
		fmt.Println("OS Version:\n", osVersion)

	case "ProcessorFlags":
		fmt.Println("Processor Flags method called")
		processorFlags, err := p.ProcessorFlags()
		if err != nil {
			fmt.Println("Error during retrieval of Processor Flags: ", err)
			os.Exit(1)
		}
		fmt.Println("Processor Flags:\n", processorFlags)

	case "ProcessorID":
		fmt.Println("Processor ID method called")
		processorID, err := p.ProcessorID()
		if err != nil {
			fmt.Println("Error during retrieval of Processor ID: ", err)
			os.Exit(1)
		}
		fmt.Println("Processor Info:\n", processorID)

	case "VMMName":
		fmt.Println("VMM Name method called")
		vmmName, err := p.VMMName()
		if err != nil {
			fmt.Println("Error during retrieval of VMM Name: ", err)
			os.Exit(1)
		}
		fmt.Println("VMM Name:\n", vmmName)

	case "VMMVersion":
		fmt.Println("VMM Version method called")
		vmmVersion, err := p.VMMVersion()
		if err != nil {
			fmt.Println("Error during retrieval of VMM Version: ", err)
			os.Exit(1)
		}
		fmt.Println("VMM Version:\n", vmmVersion)

	case "TPMVersion":
		fmt.Println("TPM Version method called")
		tpmVersion, err := p.TPMVersion()
		if err != nil {
			fmt.Println("Error during retrieval of TPM Version: ", err)
			os.Exit(1)
		}
		fmt.Println("TPM Version:\n", tpmVersion)

	case "HostName":
		fmt.Println("Host Name method called")
		hostName, err := p.HostName()
		if err != nil {
			fmt.Println("Error during retrieval of Host Name: ", err)
			os.Exit(1)
		}
		fmt.Println("Host Name:\n", hostName)

	case "NoOfSockets":
		fmt.Println("Number of Sockets method called")
		noOfSockets, err := p.NoOfSockets()
		if err != nil {
			fmt.Println("Error during retrieval of Number of Sockets: ", err)
			os.Exit(1)
		}
		fmt.Println("Number of sockets:\n", noOfSockets)

	case "TPMEnabled":
		fmt.Println("TPM Enabled method called")
		tpmEnabled, err := p.TPMEnabled()
		if err != nil {
			fmt.Println("Error during retrieval of status of TPM: ", err)
			os.Exit(1)
		}
		fmt.Println("TPM Enabled:\n", tpmEnabled)

	case "TXTEnabled":
		fmt.Println("TXT Enabled method called")
		txtEnabled, err := p.TXTEnabled()
		if err != nil {
			fmt.Println("Error during retrieval of status of TXT: ", err)
			os.Exit(1)
		}
		fmt.Println("TXT Enabled:\n", txtEnabled)

	default:
		fmt.Println("Invalid method name mentioned.\nExpected values: BiosName, BiosVersion, HardwareUUID, " +
			"OSName, OSVersion, ProcessorFlags, ProcessorInfo, VMMName, VMMVersion, TPMVersion, HostName, " +
			"NoOfSockets, TPMEnabled, TXTEnabled")
	}
}
