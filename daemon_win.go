// +build windows

package docker

const dockerExe = "C:\\bin\\docker.exe"
const dockerdExe = ""
const dockerHome = "C:\\ProgramData\\docker\\"
const dockerTrustStore = dockerHome+"trust/private"

func (p Plugin) startDaemon() {
	// this is a no-op on windows
}
