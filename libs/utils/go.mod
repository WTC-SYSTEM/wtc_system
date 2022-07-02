module github.com/WTC-SYSTEM/wtc_system/libs/utils

go 1.18

replace github.com/WTC-SYSTEM/wtc_system/libs/logging => ../logging

require github.com/WTC-SYSTEM/wtc_system/libs/logging v0.0.0-00010101000000-000000000000

require (
	github.com/sirupsen/logrus v1.8.1 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)
