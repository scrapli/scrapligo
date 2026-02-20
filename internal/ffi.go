package internal

type driverOptions struct {
	loggerCallback uintptr
	loggerLevel    uintptr
	loggerLevelLen uintptr

	port             *uint16
	transportKind    uintptr
	transportKindLen uintptr

	cli struct {
		definitionStr    uintptr
		definitionStrLen uintptr
	}

	netconf struct {
		errorTag             uintptr
		errorTagLen          uintptr
		preferredVersion     uintptr
		preferredVersionLen  uintptr
		messagePollInterval  *uint64
		capabilitiesCallback uintptr
	}

	session struct {
		readSize                *uint64
		readMinDelayNs          *uint64
		readMaxDelayNs          *uint64
		returnChar              uintptr
		returnCharLen           uintptr
		operationTimeoutNs      *uint64
		operationMaxSearchDepth *uint64
		recordDestination       uintptr
		recordDestinationLen    uintptr
		recorderCallback        uintptr
	}

	auth struct {
		username                uintptr
		usernameLen             uintptr
		password                uintptr
		passwordLen             uintptr
		privateKeyPath          uintptr
		privateKeyPathLen       uintptr
		privateKeyPassphrase    uintptr
		privateKeyPassphraseLen uintptr
		lookups                 struct {
			keys     uintptr
			keysLens uintptr
			vals     uintptr
			valsLens uintptr
			count    uint16
		}
		forceInSessionAuth             *bool
		bypassInSessionAuth            *bool
		usernamePattern                uintptr
		usernamePatternLen             uintptr
		passwordPattern                uintptr
		passwordPatternLen             uintptr
		privateKeyPassphrasePattern    uintptr
		privateKeyPassphrasePatternLen uintptr
	}

	transport struct {
		bin struct {
			bin                 uintptr
			binLen              uintptr
			extraOpenArgs       uintptr
			extraOpenArgsLen    uintptr
			overrideOpenArgs    uintptr
			overrideOpenArgsLen uintptr
			sshConfigPath       uintptr
			sshConfigPathLen    uintptr
			knownHostsPath      uintptr
			knownHostsPathLen   uintptr
			enableStrictKey     *bool
			termHeight          *uint16
			termWidth           *uint16
		}

		ssh2 struct {
			knownHostsPath                   uintptr
			knownHostsPathLen                uintptr
			libssh2Trace                     *bool
			proxyJumpHost                    uintptr
			proxyJumpHostLen                 uintptr
			proxyJumpPort                    *uint16
			proxyJumpUsername                uintptr
			proxyJumpUsernameLen             uintptr
			proxyJumpPassword                uintptr
			proxyJumpPasswordLen             uintptr
			proxyJumpPrivateKeyPath          uintptr
			proxyJumpPrivateKeyPathLen       uintptr
			proxyJumpPrivateKeyPassphrase    uintptr
			proxyJumpPrivateKeyPassphraseLen uintptr
			proxyJumpLibssh2Trace            *bool
		}

		test struct {
			f    uintptr
			fLen uintptr
		}
	}
}
