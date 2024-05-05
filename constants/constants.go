package constants

const NAME = "PassEz"
const VER = "(WIP)"

const DEFAULT_ITERATION = 100000
const DEFAULT_SALT_LEN = 32
const DEFAULT_KEY_LEN = 32

const ENC_FILE_EXT = ".json"
const PARAMS_FILE_EXT = ".params"

type Service struct {
	ServiceName string       `json:"serviceName"`
	Credentials []Credential `json:"credentials"`
}

type Credential struct {
	Identifier string `json:"identifier"`
	Password   []byte `json:"password"`
}
